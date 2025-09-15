package goplayground_validator

type User struct {
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age" validate:"gte=0,lte=130"`
	Sex  string `json:"sex" validate:"sex"`
	// dive を付けると、入れ子の中や配列の各要素を見てくれる
	// see https://pkg.go.dev/gopkg.in/go-playground/validator.v10#hdr-Dive
	Families []Family `json:"families" validate:"dive"`
}

type Family struct {
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age" validate:"gte=0,lte=130"`
}
