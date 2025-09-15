package ozzo_validation

import (
	"slices"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type User struct {
	Name     string   `json:"name"`
	Age      int      `json:"age"`
	Sex      string   `json:"sex"`
	Families []Family `json:"families"`
}

func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Name, validation.Required),
		validation.Field(&u.Age, validation.Min(0), validation.Max(130)),
		validation.Field(&u.Sex, validation.By(func(value any) error {
			sex, _ := value.(string)
			if !slices.Contains([]string{"male", "female"}, sex) {
				return validation.NewError("validation_sex", "should be male or female")
			}
			return nil
		})),
		validation.Field(&u.Families, validation.Each()),
	)
}

type Family struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (f Family) Validate() error {
	return validation.ValidateStruct(&f,
		validation.Field(&f.Name, validation.Required),
		validation.Field(&f.Age, validation.Min(0), validation.Max(130)),
	)
}
