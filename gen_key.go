package main

import (
	"context"
	"encoding/base64"
	"fmt"

	//"io/ioutil"
	"log"
	//"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ctx          = context.Background()
	kmsProviders map[string]map[string]interface{}
	schemaMap    bson.M
)

func createDataKey() {
	kvClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://ben:pass7word@localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	clientEncryptionOpts := options.ClientEncryption().SetKeyVaultNamespace("keyvault.datakeys").SetKmsProviders(kmsProviders)
	clientEncryption, err := mongo.NewClientEncryption(kvClient, clientEncryptionOpts)
	if err != nil {
		log.Fatal(err)
	}
	defer clientEncryption.Close(ctx)
	dataKey, err := clientEncryption.CreateDataKey(ctx, "local", options.DataKey().SetKeyAltNames([]string{"example"}))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("data key base64 is %s, plz input fill this filed into `collection_schema_d.json`'s line8 base64 field\n", base64.StdEncoding.EncodeToString(dataKey.Data))
}

func main() {
	localMasterKey := "YOm8sMifl7BUJW8vEw4UGpGKtFDIooyat4DDTDmPI+og7PsERJZVE2ldsEanYN58HhUkl8LxLjjXRyc2ctQG/Gpjg8xUqAE1XwMgyXxYnwN7MnJYSC+0msDmyMybySny"
	decodedKey, err := base64.StdEncoding.DecodeString(localMasterKey)
	if err != nil {
		log.Fatalf("base64 decode error: %v", err)
	}
	kmsProviders = map[string]map[string]interface{}{
		"local": {
			"key": decodedKey,
		},
	}
	createDataKey()
	fmt.Printf("create key ok!")
}
