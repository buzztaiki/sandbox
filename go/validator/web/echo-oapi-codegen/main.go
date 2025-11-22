package main

import (
	"github.com/buzztaiki/sandbox/go/validator/web/echo-oapi-codegen/petstore"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/oapi-codegen/echo-middleware"
)

func main() {
	e := echo.New()

	swagger, err := petstore.GetSwagger()
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.Use(echomiddleware.OapiRequestValidator(swagger))
	petstore.RegisterHandlersWithBaseURL(e, NewPetHandler(), "/v3")
	e.Logger.Fatal(e.Start(":1323"))
}
