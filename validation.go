package rhapvalidator

import (
	"regexp"
)

/*
	String validations
*/
func Required(value string) bool {
	return value != ""
}

func Len(value string, length int) bool {
	return len(value) == length
}

func IsAlpha(value string) bool {
	reg := regexp.MustCompile("^[a-zA-Z ]*$")
	return reg.MatchString(value)
}

func IsEmail(email string) bool {
	reg := regexp.MustCompile("^(([^<>()[\\]\\\\.,;:\\s@\"]+(\\.[^<>()[\\]\\\\.,;:\\s@\"]+)*)|(\".+\"))@((\\[[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}])|(([a-zA-Z\\-0-9]+\\.)+[a-zA-Z]{2,}))$")
	return reg.MatchString(email)
}

func MinString(value string, min int) bool {
	return len(value) >= min
}

func MaxString(value string, max int) bool {
	return len(value) <= max
}

/*
	Integer validations
*/

func RequiredNum(num int) bool {
	return num > 0
}

func MinNum(num int, min int) bool {
	return num >= min
}

func MaxNum(num int, max int) bool {
	return num <= max
}