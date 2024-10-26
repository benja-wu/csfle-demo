# CSFLE Demo
The Golang demo code for MongoDB CSFLE(client-side-field-level-encryption)

## Overview
1. The `example.go` and `example_d.go` are different encryption type for inserting and querying documents from MongoDB with CSFLE.
2. `collection_schema.json` and `collection_schema_d.json` are two schema validation json object for validating the target
field's input with the CSFLE. If the client can't find the target field's encryption type, it will try to download from server side.

## Schema Validation
* In the source code: `example_d.go`, encrypted client defines  target encryption fields' encryption type in `collection_schema_d.json`. Can also define the schema in server side in `schema validator` configuration.
Such as: 
```yaml
{
  $jsonSchema: {
    required: [
      'ssn',
      'plaintext_string',
      'plaintext_num'
    ],
    properties: {
      ssn: {
        encrypt: {
          keyId: '/altname',
          bsonType: 'string',
          algorithm: 'AEAD_AES_256_CBC_HMAC_SHA_512-Random'
        }
      },
      plaintext_string: {
        bsonType: 'string'
      },
      plaintext_num: {
        bsonType: 'int'
      }
    },
    bsonType: 'object'
  }
}
```
* We declare the `ssn` filed's encrytpion type and associated key. It can coexist with other unencrypted fields' schema validation role.  
* Check out the supportted BSON type in official manual https://www.mongodb.com/docs/manual/core/csfle/reference/encryption-schemas/#mongodb-autoencryptkeyword-autoencryptkeyword.encrypt.bsonType. Or it will failed with error message :
```bash
2024/10/26 22:06:43 InsertOne error: mongocryptd communication error: (Location31122) Cannot encrypt element of type: object
exit status 1
```

## Procedure 
### Before mongoDB 4.4
Prerequisites:
1. Use `mongocryptd --port 27020` to start the mongocryptd instance locally
2. Run `go mod tidy` to install necessary Golang libraries
3. Ensure `localhost:27017` mongoDB instance is running and the correct URI is exposed with envionment value `MONGODB_URI`  
4. Use `go run -tags cse gen_key.go` firstly, for creating local key into MongoDB key vault collection. It will generate the output as below 
```bash
data key base64 is MY5IbvjvSN+ttrHljOu4hw==, plz input fill this filed into `collection_schema_d.json`'s line8 base64 field
create key ok!%
```
5. Replace collection_schema_d.json file's line8 with the output from step4. 

### mongoDB 4.4 +
Prerequisites:
1. Install encryption shard libiary according to https://www.mongodb.com/docs/manual/core/queryable-encryption/install-library/ 
2. Export the shard library absoult path with OS environment value `MDB_CRYPT_SHARED_LIB_PATH`
Repeat the `Before mongoDB 4.4` procedure above from **step 2 to step 5**



### Demo CSFLE: encryption field with random value
1. Use `go run -tags cse example.go` run the Golang application with cse tag enabled.



### Demo CSFLE: encroyption field with deterministic value
1. Use `go run -tags cse example_d.go` run the Golang application with cse tag enabled.
2. When using 

