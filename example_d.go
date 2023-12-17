package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	uri = "mongodb://localhost:27017"
)

var (
	ctx = context.Background()
	// Test key material generated by: echo $(head -c 96 /dev/urandom | base64 | tr -d '\n')
	localMasterKey = "YOm8sMifl7BUJW8vEw4UGpGKtFDIooyat4DDTDmPI+og7PsERJZVE2ldsEanYN58HhUkl8LxLjjXRyc2ctQG/Gpjg8xUqAE1XwMgyXxYnwN7MnJYSC+0msDmyMybySny"
	kmsProviders   map[string]map[string]interface{}
)

func main() {
	// initial setup
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

	// insert a document with an encrypted field and a plaintext field
	_, err = coll.InsertOne(ctx, bson.M{
		"plaintext": "hello world2",
		"ssn":       "123-00-6789",
	})
	if err != nil {
		log.Fatalf("InsertOne error: %v", err)
	}

	// find and print the inserted document
	res, err := coll.FindOne(ctx, bson.D{}).DecodeBytes()
	if err != nil {
		log.Fatalf("FindOne error: %v", err)
	}
	fmt.Println(res)
}

// create a client configured with auto encryption that uses the key generated by createDataKey
func createEncryptedClient() *mongo.Client {
	// create a client with auto encryption
	schemaMap := map[string]interface{}{
		"foo.bar": readJSONFile("collection_schema_d.json"),
	}
	autoEncOpts := options.AutoEncryption().
		SetKeyVaultNamespace("keyvault.datakeys").
		SetKmsProviders(kmsProviders).
		SetSchemaMap(schemaMap)

	clientOpts := options.Client().ApplyURI(uri).SetAutoEncryptionOptions(autoEncOpts)
	autoEncryptionClient, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatalf("Connect error for client with automatic encryption: %v", err)
	}
	return autoEncryptionClient
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
