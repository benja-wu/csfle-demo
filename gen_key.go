package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	//"io/ioutil"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ctx                    = context.Background()
	kmsProviders           map[string]map[string]interface{}
	cryptSharedLibraryPath map[string]interface{}
	schemaMap              bson.M
)

func createDataKey() {
	uri := os.Getenv("MONGODB_URI")
	autoEncOpts := options.AutoEncryption().
		SetKeyVaultNamespace("keyvault.datakeys").
		SetKmsProviders(kmsProviders).
		SetExtraOptions(cryptSharedLibraryPath)

	clientOpts := options.Client().ApplyURI(uri).SetAutoEncryptionOptions(autoEncOpts)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatalf("Connect error for client with automatic encryption: %v", err)
	}

	opts := options.ClientEncryption().
		SetKeyVaultNamespace("keyvault.datakeys").
		SetKmsProviders(kmsProviders)

	clientEncryption, err := mongo.NewClientEncryption(client, opts)

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
	cryptSharedLibraryPath = map[string]interface{}{
		"cryptSharedLibPath": os.Getenv("MDB_CRYPT_SHARED_LIB_PATH"), // Path to your Automatic Encryption Shared Library
	}

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
