# CSFLE Demo
The Golang demo with the MongoDB client-side-field-level-encryption 

## Overview
1. `example.go` and `example_d.go` are example code for inserting and geting the result from MongoDB with CSFLE
2. `collection_schema.json` and `collection_schema_d.json` are two schema validation json object for validating the target
field's input with the CSFLE. If the client can't find the target field's encryption type, it will try to download from server side.

## Schema Validation
* In the `example_d.go`, encrypted client uses the `collection_schema_d.json` to indicate target encryption field's encryption type. This can also be done by complete the `schema validator` configuration in MongoDB server side(just as the `example.go`'s implementation). Here is an example 
```yaml
{
  $jsonSchema: {
    required: [
      'ssn',
      'plaintext'
    ],
    properties: {
      ssn: {
        encrypt: {
          keyId: '/altname',
          bsonType: 'string',
          algorithm: 'AEAD_AES_256_CBC_HMAC_SHA_512-Random'
        }
      },
      plaintext: {
        bsonType: 'string'
      }
    },
    bsonType: 'object'
  }
}
```
* In the example above, we declare the type of `ssn` filed's encrytpion type and used key, also works with the other field's schema validation rule. 

## Procedure 
Prerequests:
1. Use `mongocryptd --port 27020` to start the mongocryptd instance locally
2. Run 'go mod tidy' to install necessary local libraries
3. Make sure localhost:27017 mongoDB instance is running  
4. Use `go run -tags cse clinet.go` firstly, for creating local key into MongoDB key vault collection. It will generate the output as below 
```bash
data key base64 is MY5IbvjvSN+ttrHljOu4hw==, plz input fill this filed into `collection_schema_d.json`'s line8 base64 field
create key ok!%
```
Remind to replace the base64 filed value into `collection_sschema_d.json` file if u want to run `example_d.go`

### Demo CSFLE in random 
1. Use `go run -tags cse example.go` run the Golang application with cse tag enabled.


### Demo CSFLE in deterministic 
1. Use `go run -tags cse example_d.go` run the Golang application with cse tag enabled.

