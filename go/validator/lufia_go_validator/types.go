package lufia_go_validator

import (
	"github.com/lufia/go-validator"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Sex  string `json:"sex"`
	// validator.Struct が Validator[*Family] になるから、ポインタにするしかない
	Families []*Family `json:"families"`
}

var UserValidator = validator.Struct(func(s validator.StructRule, u *User) {
	validator.AddField(s, &u.Name, "name", validator.Required[string]())
	validator.AddField(s, &u.Age, "age", validator.Min(0), validator.Max(130))
	validator.AddField(s, &u.Families, "families", validator.Slice(FamilyValidator))
})

type Family struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// Validator[*Family] を Validator[Family] に変える何かがあると良さそう
var FamilyValidator = validator.Struct(func(s validator.StructRule, f *Family) {
	validator.AddField(s, &f.Name, "name", validator.Required[string]())
	validator.AddField(s, &f.Age, "age", validator.Min(0), validator.Max(130))
})
