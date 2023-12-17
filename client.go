package main

import (
	"context"
	"crypto/rand"
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
	kvClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	clientEncryptionOpts := options.ClientEncryption().SetKeyVaultNamespace("keyvault.datakeys").SetKmsProviders(kmsProviders)
	clientEncryption, err := mongo.NewClientEncryption(kvClient, clientEncryptionOpts)
	if err != nil {
		log.Fatal(err)
	}
	defer clientEncryption.Close(ctx)
	_, err = clientEncryption.CreateDataKey(ctx, "local", options.DataKey().SetKeyAltNames([]string{"example"}))
	if err != nil {
		log.Fatal(err)
	}
}


func main() {
	localKey := make([]byte, 96)
	if _, err := rand.Read(localKey); err != nil {
		log.Fatal(err)
	}
	kmsProviders = map[string]map[string]interface{}{
		"local": {
			"key": localKey,
		},
	}
	createDataKey()
	fmt.Printf("create key ok!")
}
