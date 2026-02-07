package ozzo_validation

import (
	"encoding/json"
	"fmt"
	"maps"
	"slices"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
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
	err := user.Validate()
	fmt.Println(err)

	b, _ := json.Marshal(err)
	fmt.Println(string(b))
	fmt.Println()

	errs := err.(validation.Errors)
	for _, k := range slices.Sorted(maps.Keys(errs)) {
		err := errs[k]
		if verr, ok := err.(validation.Error); ok {
			fmt.Printf("key=%q, code=%q, params=%v, message=%q\n", k, verr.Code(), verr.Params(), verr.Message())
		} else {
			fmt.Printf("key=%q, errT=%T, err=%v\n", k, err, err)

		}
	}

	// Output:
	// age: must be no greater than 130; families: (0: (name: cannot be blank.).); name: cannot be blank; sex: should be male or female.
	// {"age":"must be no greater than 130","families":{"0":{"name":"cannot be blank"}},"name":"cannot be blank","sex":"should be male or female"}
	//
	// key="age", code="validation_max_less_equal_than_required", params=map[threshold:130], message="must be no greater than {{.threshold}}"
	// key="families", errT=validation.Errors, err=0: (name: cannot be blank.).
	// key="name", code="validation_required", params=map[], message="cannot be blank"
	// key="sex", code="validation_sex", params=map[], message="should be male or female"
}

func Example_translate() {
	bundle, err := NewI18nBundle()
	if err != nil {
		panic(err)
	}

	localizer := i18n.NewLocalizer(bundle, "ja")
	user := User{
		Name: "",
		Age:  150,
		Sex:  "aaa",
		Families: []Family{
			{Name: "", Age: 30},
		},
	}

	errs := user.Validate().(validation.Errors)
	TranslateErrors(errs, localizer)
	b, _ := json.Marshal(errs)
	fmt.Println(string(b))

	// Output:
	// {"age":"130 を超えています","families":{"0":{"name":"必須項目です"}},"name":"必須項目です","sex":"male または femail を指定してください"}
}
