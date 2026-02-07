package main

//go:generate sh -c "cd ../schema && npx tsp compile petstore"
//go:generate go tool ogen --target petstore --package petstore --clean ../schema/tsp-output/@typespec/openapi3/PetStore.yaml
