package goplayground_validator

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
)

func Example() {
	user := User{
		Name: "",
		Age:  150,
		Sex:  "aaa",
		Families: []Family{
			{Name: "", Age: 30},
		},
	}

	err := NewValidate().Struct(&user)
	fmt.Println("PRINT ALL:")
	fmt.Println(err)
	fmt.Println()

	fmt.Println("PRINT EACH:")
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			attrPairs := []string{
				"Tag", e.Tag(),
				"ActualTag", e.ActualTag(),
				"Param", e.Param(),
				"Namespace", e.Namespace(),
				"Field", e.Field(),
				"StructNamespace", e.StructNamespace(),
				"StructField", e.StructField(),
			}

			var attrs []string
			for pair := range slices.Chunk(attrPairs, 2) {
				attrs = append(attrs, fmt.Sprintf("%s=%q", pair[0], pair[1]))
			}
			fmt.Println(strings.Join(attrs, ", "))
		}
	}

	// Output:
	// PRINT ALL:
	// Key: 'User.name' Error:Field validation for 'name' failed on the 'required' tag
	// Key: 'User.age' Error:Field validation for 'age' failed on the 'lte' tag
	// Key: 'User.sex' Error:Field validation for 'sex' failed on the 'sex' tag
	// Key: 'User.families[0].name' Error:Field validation for 'name' failed on the 'required' tag
	//
	// PRINT EACH:
	// Tag="required", ActualTag="required", Param="", Namespace="User.name", Field="name", StructNamespace="User.Name", StructField="Name"
	// Tag="lte", ActualTag="lte", Param="130", Namespace="User.age", Field="age", StructNamespace="User.Age", StructField="Age"
	// Tag="sex", ActualTag="sex", Param="", Namespace="User.sex", Field="sex", StructNamespace="User.Sex", StructField="Sex"
	// Tag="required", ActualTag="required", Param="", Namespace="User.families[0].name", Field="name", StructNamespace="User.Families[0].Name", StructField="Name"
}

func Example_translation() {
	user := User{
		Name: "",
		Age:  150,
		Families: []Family{
			{Name: "", Age: 30},
		},
	}

	validate := NewValidate()
	uni, err := SetupTranslation(validate)
	if err != nil {
		panic(err)
	}
	trans, _ := uni.GetTranslator("ja")

	err = validate.Struct(&user)
	errs := err.(validator.ValidationErrors)

	transMap := errs.Translate(trans)
	for _, k := range slices.Sorted(maps.Keys(transMap)) {
		fmt.Printf("%v: %v\n", k, transMap[k])
	}

	// Output:
	// User.age: ageは130以下でなければなりません
	// User.families[0].name: nameは必須フィールドです
	// User.name: nameは必須フィールドです
	// User.sex: sex には male または femail を指定してください
}
