package strs

import (
	"regexp"
)

var EmailPattern = regexp.MustCompile(`^([a-zA-Z0-9][-_.a-zA-Z0-9]*)(@[-_.a-zA-Z0-9]+)?$`)
var phoneNumberRe = regexp.MustCompile(`^(\+86)?\d{11}$`)

func IsValidEmail(str string) bool {
	return EmailPattern.MatchString(str)
}

func IsValidPhone(phoneNumber string) bool {
	return phoneNumberRe.Match([]byte(phoneNumber))
}
