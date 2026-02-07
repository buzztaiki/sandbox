package main

//go:generate sh -c "cd ../schema && npx tsp compile petstore"
//go:generate go tool oapi-codegen --generate server --package petstore -o petstore/server_gen.go ../schema/tsp-output/@typespec/openapi3/PetStore.yaml
