package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ctx = context.Background()
	// Test key material generated by: echo $(head -c 96 /dev/urandom | base64 | tr -d '\n')
	localMasterKey string
	kmsProviders   map[string]map[string]interface{}
)

func ssnWithObject(coll *mongo.Collection, random_str, ssn_num string) {
	ssn := bson.M{
		"ssn_num": ssn_num,
		"name":    "ben",
		"age":     30,
		"address": []string{"123 Main St, Apt 4B, NY", "456 Main St, Apt B, NY"},
	}

	// insert a document with an encrypted field and a plaintext field
	_, err := coll.InsertOne(ctx, bson.M{
		"plaintext_string": random_str,
		"plaintext_num":    2,
		"meta_info":        fmt.Sprintf("%s_%s", "hello world2", time.Now().Format("2006-01-02 15:04:05")),
		"altname":          "example",
		"ssn":              ssn,
		"msg":              "encrpyted with random type",
	})
	if err != nil {
		log.Fatalf("InsertOne error: %v", err)
	}

	fmt.Printf("insert successfully\n")
	fmt.Println()

	fmt.Printf("after insert, query the doc with plaintext_string: %s\n", random_str)
	res := coll.FindOne(ctx, bson.M{"plaintext_string": random_str})
	raw, err := res.Raw()
	if err != nil {
		log.Fatalf("FindOne with SSN object error: %v", err)
	}
	fmt.Println(raw)
	// Use FindOneAndUpdate to set ssn.age to 29
	ssn = bson.M{
		"ssn_num": ssn_num,
		"name":    "ben",
		"age":     29,
		"address": []string{"123 Main St, Apt 4B, NY", "456 Main St, Apt B, NY"},
	}
	ures, err := coll.UpdateOne(
		ctx,
		bson.M{"plaintext_string": random_str},
		bson.M{
			"$set": bson.M{
				"ssn": ssn,
			},
		},
	)
	if err != nil {
		log.Fatalf("UpdateOne error: %v", res.Err())
	}

	fmt.Printf("the update result %v\n", ures)

	res = coll.FindOne(ctx, bson.M{"plaintext_string": random_str})
	raw, err = res.Raw()
	if err != nil {
		log.Fatalf("FindOne with SSN object error: %v", err)
	}

	fmt.Printf("after update:\n")
	fmt.Println(raw)

	// Call the aggregation function
	/*
			if err := aggregateTagsByReportNum(coll, 9); err != nil {
		        log.Fatalf("aggregateTagsByReportNum error: %v", err)
	}*/
}

func ssnWithString(coll *mongo.Collection, random_str, ssn_num string) {

	// insert a document with an encrypted field and a plaintext field
	_, err := coll.InsertOne(ctx, bson.M{
		"plaintext_string": random_str,
		"plaintext_num":    2,
		"meta_info":        fmt.Sprintf("%s_%s", "hello world2", time.Now().Format("2006-01-02 15:04:05")),
		"ssn":              ssn_num,
		"altname":          "example",
		"msg":              "encrpyted with deterministic type",
	})
	if err != nil {
		log.Fatalf("InsertOne error: %v", err)
	}
	fmt.Printf("insert successfully\n")
	fmt.Println()

	fmt.Printf("after insert, query the doc with plaintext_string: %s\n", random_str)
	res := coll.FindOne(ctx, bson.M{"plaintext_string": random_str})
	raw, err := res.Raw()
	if err != nil {
		log.Fatalf("FindOne with SSN object error: %v", err)
	}
	fmt.Println(raw)

	fmt.Printf("after insert, query the doc with ssn: %s\n", ssn_num)
	res = coll.FindOne(ctx, bson.M{"ssn": ssn_num})
	raw, err = res.Raw()
	if err != nil {
		log.Fatalf("FindOne with SSN object error: %v", err)
	}
	fmt.Println(raw)
}

func main() {
	// initial setup
	localMasterKey = os.Getenv("LOCAL_MASTER_KEY")
	decodedKey, err := base64.StdEncoding.DecodeString(localMasterKey)
	if err != nil {
		log.Fatalf("base64 decode error: %v", err)
	}
	kmsProviders = map[string]map[string]interface{}{
		"local": {"key": decodedKey},
	}

	client := createEncryptedClient()
	defer client.Disconnect(ctx)

	coll := client.Database("foo").Collection("bar")

	t := time.Now()
	r := rand.New(rand.NewSource(t.UnixNano()))
	min := 100000
	max := 300000
	random_str := fmt.Sprintf("%s_%d", "determinstic encrpytion", r.Intn(max-min+1)+min)
	fmt.Printf("before insert, prepare doc with plaintext_string: %s\n", random_str)

	// Generate each part of the SSN
	part1 := rand.Intn(900) + 100   // 3 digits, 100–999
	part2 := rand.Intn(90) + 10     // 2 digits, 10–99
	part3 := rand.Intn(9000) + 1000 // 4 digits, 1000–9999

	// Format as "XXX-XX-XXXX"
	ssn := fmt.Sprintf("%03d-%02d-%04d", part1, part2, part3)
	fmt.Printf("               prepare ssn: %s\n", ssn)
	fmt.Println()

	// find and print the inserted document
	ssnWithObject(coll, random_str, ssn)
	//ssnWithString(coll, random_str, ssn)
}

// create a client configured with auto encryption that uses the key generated by createDataKey
func createEncryptedClient() *mongo.Client {
	// initial setup
	decodedKey, err := base64.StdEncoding.DecodeString(localMasterKey)
	if err != nil {
		log.Fatalf("base64 decode error: %v", err)
	}
	kmsProviders = map[string]map[string]interface{}{
		"local": {"key": decodedKey},
	}

	cryptSharedLibraryPath := map[string]interface{}{
		"cryptSharedLibPath": os.Getenv("MDB_CRYPT_SHARED_LIB_PATH"), // Path to your Automatic Encryption Shared Library
	}

	uri := os.Getenv("MONGODB_URI")
	// create a client with auto encryption
	schemaMap := map[string]interface{}{
		"foo.bar": readJSONFile("collection_schema.json"),
	}
	autoEncOpts := options.AutoEncryption().
		SetKeyVaultNamespace("keyvault.datakeys").
		SetKmsProviders(kmsProviders).
		SetSchemaMap(schemaMap).
		SetExtraOptions(cryptSharedLibraryPath)

	clientOpts := options.Client().ApplyURI(uri).SetAutoEncryptionOptions(autoEncOpts)
	autoEncryptionClient, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatalf("Connect error for client with automatic encryption: %v", err)
	}
	return autoEncryptionClient
}

// this aggreagte won't work because of the ssn.tags encryption path
func aggregateTagsByReportNum(coll *mongo.Collection, reportNum int) error {
	fmt.Println("number is", reportNum)

	// Define the aggregation pipeline
	pipeline := mongo.Pipeline{
		// $match stage
		bson.D{{"$match", bson.D{{"report_num", reportNum}}}},
		// $project stage
		bson.D{{"$project", bson.D{
			{"_id", 0},
			{"tags", "$ssn.tags"},
		}}},
		// $unwind stage
		bson.D{{"$unwind", "$tags"}},
		// $group stage
		bson.D{{"$group", bson.D{
			{"_id", "$tags"},
			{"count", bson.D{{"$sum", 1}}},
		}}},
	}

	// Execute the aggregation
	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return fmt.Errorf("Aggregate error: %v", err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return fmt.Errorf("cursor.All error: %v", err)
	}

	fmt.Println("after agg")
	for _, result := range results {
		fmt.Println(result)
	}
	return nil
}

func readJSONFile(file string) bson.D {
	content, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("ReadFile error for %v: %v", file, err)
	}

	var fileDoc bson.D
	if err = bson.UnmarshalExtJSON(content, false, &fileDoc); err != nil {
		log.Fatalf("UnmarshalExtJSON error for file %v: %v", file, err)
	}
	return fileDoc
}
