package main

//go:generate go tool oapi-codegen --generate server --package petstore -o petstore/server_gen.go ../petstore.yml
//go:generate go tool oapi-codegen --generate types --package petstore -o petstore/types_gen.go ../petstore.yml
//go:generate go tool oapi-codegen --generate spec --package petstore -o petstore/spec_gen.go ../petstore.yml
