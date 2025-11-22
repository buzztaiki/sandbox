package main

import (
	"github.com/buzztaiki/sandbox/go/validator/web/echo-oapi-codegen/petstore"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	petstore.RegisterHandlers(e, NewPetHandler())
	e.Logger.Fatal(e.Start(":1323"))
}
