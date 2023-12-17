# CSFLE Demo
The Golang demo with the MongoDB client-side-field-level-encryption 

## Overview
1. `example.go` and `example_d.go` are example code for inserting and geting the result from MongoDB with CSFLE
2. `collection_schema.json` and `collection_schema_d.json` are two schema validation json object for validating the target
collection's input with CSFLE. Notice, if once we enable this config, it will overwrite the one stored in MongoDB side

## Procedure 
Prerequests:
1. Use `mongocryptd --port 27020` to start the mongocryptd instance locally
2. Run 'go mod tidy' to install necessary local libraries
3. Make sure localhost:27017 mongoDB instance is running  
4. Use `go run clinet.go` firstly, for creating local key into MongoDB key vault collection 


### Demo CSFLE with random value 
1. Use `go run -tags cse example.go` run the Golang application with cse tag enabled.


### Demo CSFLE with deterministic value 
1. Use `go run -tags cse example_d.go` run the Golang application with cse tag enabled.
