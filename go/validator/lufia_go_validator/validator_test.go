package lufia_go_validator

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/lufia/go-validator"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
)

func Example() {
	user := User{
		Name: "",
		Age:  150,
		Sex:  "aaa",
		Families: []*Family{
			{Name: "", Age: 30},
		},
	}

	err := UserValidator.Validate(context.Background(), &user)
	fmt.Println(err)
	fmt.Println()

	if errs, ok := err.(interface{ Unwrap() []error }); ok {
		for _, e := range errs.Unwrap() {
			fmt.Printf("%v\n", e)
		}
	}

	// Unordered output:
	// families: name: cannot be the zero value
	// name: cannot be the zero value
	// age: must be no greater than 130
	//
	// families: name: cannot be the zero value
	// name: cannot be the zero value
	// age: must be no greater than 130
}

func Example_translate() {
	user := User{
		Name: "",
		Age:  150,
		Sex:  "aaa",
		Families: []*Family{
			{Name: "", Age: 30},
			{Name: "aaa", Age: 1000},
		},
	}

	{
		p := message.NewPrinter(language.Japanese, message.Catalog(validator.DefaultCatalog))
		ctx := validator.WithPrinter(context.Background(), p)

		err := UserValidator.Validate(ctx, &user)
		fmt.Println(err)
		fmt.Println()

		// JSONでエラーが出てくるわけじゃない
		json.NewEncoder(os.Stdout).Encode(err)
		fmt.Println()
	}

	// JSON で出したいけど難しそうかも。`Unwrap() []error` して、自分でメッセージパースとかやらないと無理そう。
	// Unwrap した場合、SliceError の位置情報が消えてしまうはず
	// StructError と SliceError が公開されてるから、Unwrap しなければOK？と思ったけど無理そう。
	{
		cat := catalog.NewBuilder(catalog.Fallback(language.English))
		// structFieldErrorFormat を `key=value` に変える
		cat.SetString(language.English, "%[1]s: %[2]v", "%[1]s=%[2]v")
		p := message.NewPrinter(language.English, message.Catalog(cat))
		ctx := validator.WithPrinter(context.Background(), p)

		err := UserValidator.Validate(ctx, &user)
		fmt.Println(err)
		fmt.Println()
	}

	// Unordered output:
	// name: 必須です
	// age: 130以下の値が必要です
	// families: name: 必須です
	// families: age: 130以下の値が必要です
	//
	// {"Value":{"name":"","age":150,"sex":"aaa","families":[{"name":"","age":30},{"name":"aaa","age":1000}]},"Errors":{"age":{},"families":{},"name":{}}}
	//
	// name=cannot be the zero value
	// age=must be no greater than 130
	// families=name=cannot be the zero value
	// families=age=must be no greater than 130
}

func Example_json() {
	user := User{
		Name: "",
		Age:  150,
		Sex:  "aaa",
		Families: []*Family{
			{Name: "", Age: 30},
		},
	}

	err := UserValidator.Validate(context.Background(), &user)
	json.NewEncoder(os.Stdout).Encode(CollectErrors(err))

	// Output:
	// {"age":["age: must be no greater than 130"],"families":["families: name: cannot be the zero value"],"name":["name: cannot be the zero value"]}
}
