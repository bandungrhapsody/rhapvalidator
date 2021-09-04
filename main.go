package main

import (
	"fmt"
	"rhap_validator/validator"
	"time"
)

func main() {
	u := &User{
		Name: "123",
		UserCode: "12345",
		Username: "jill",
		Email: "asd_asdd123@.com",
	}

	start := time.Now()
	defer func() {
		fmt.Println(fmt.Sprintf("%v ms", time.Since(start).Milliseconds()))
	}()

	v := validator.NewValidator().Validate(u)
	validationErrs := v.Errors()
	if validationErrs != nil {
		fmt.Println(validationErrs)
		return
	}

	fmt.Println("OK")
}