package validator

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type RhapValidator struct {
	results        ValidationMessages
	customMessages ValidationMessages
}

type ValidationMessages map[string]Message
type Message map[string]string

type Validation struct {
	Rules          string
	FieldName      string
	Type           reflect.Kind
	Label          string
	Value          reflect.Value
	CustomMessage  Message
}

var stringValidations = []*regexp.Regexp{
	regexp.MustCompile("^(required)$"),
	regexp.MustCompile("^(alpha)$"),
	regexp.MustCompile("^(email)$"),
	regexp.MustCompile("^(number)$"),
	regexp.MustCompile("^(len)=(\\d+)$"),
	regexp.MustCompile("^(min)=(\\d+)$"),
	regexp.MustCompile("^(max)=(\\d+)$"),
}

func NewValidator() *RhapValidator {
	return &RhapValidator{
		results:        make(ValidationMessages),
		customMessages: make(ValidationMessages),
	}
}

func (rv *RhapValidator) CustomMessage(fieldName string, message Message) *RhapValidator {
	rv.customMessages[fieldName] = message
	return rv
}

func (rv *RhapValidator) Validate(v interface{}) *RhapValidator {
	isPtrStruct(v)

	t := reflect.TypeOf(v).Elem()
	val := reflect.ValueOf(v).Elem()

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Tag.Get("rh_valid") == "" {
			continue
		}

		vf := val.Field(i)

		fieldName := f.Tag.Get("json")
		rv.checkDataTypeAndValidate(Validation{
			Rules:         f.Tag.Get("rh_valid"),
			FieldName:     fieldName,
			Type:          f.Type.Kind(),
			Label:         f.Tag.Get("rh_label"),
			Value:         vf,
			CustomMessage: rv.customMessages[fieldName],
		})
	}

	return rv
}

func (rv *RhapValidator) Errors() ValidationMessages {
    if len(rv.results) == 0 {
    	return nil
	}

	return rv.results
}

func isPtrStruct(v interface{}) {
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Ptr {
		panic("not a pointer to struct")
	}

	t = t.Elem()
	if t.Kind() != reflect.Struct {
		panic("not a pointer to struct")
	}
}

func (rv *RhapValidator) checkDataTypeAndValidate(validation Validation) {
	switch validation.Type {
		case reflect.String:
			rv.validateString(validation)
			break
		case reflect.Int:
			validateIntType()
			break
		case reflect.Int32:
			validateIntType()
			break
		case reflect.Int64:
			validateIntType()
			break
	}
}

/*
	String validations
*/
func (rv *RhapValidator) validateString(validation Validation) {
	slcRules := strings.Split(validation.Rules, ",")

	for _, rule := range slcRules {
		rule = strings.Trim(rule, " ")
		matches := matchStringRules(rule)
		if matches == nil {
			continue
		}

		if !rv.checkRuleAndValidateString(matches, validation) {
			return
		}
	}
}

func (rv *RhapValidator) checkRuleAndValidateString(matches []string, validation Validation) bool {
	r := matches[0]

	switch r {
	case "required":
		return rv.required(r, validation)
	case "len":
		return rv.len(matches, validation)
	case "min":
		return rv.min(matches, validation)
	case "max":
		return rv.max(matches, validation)
	case "alpha":
		return rv.alpha(r, validation)
	case "email":
		return rv.email(r, validation)
	}

	return false
}

func matchStringRules(ruleName string) (matches []string) {
	for _, r := range stringValidations {
		matches = r.FindStringSubmatch(ruleName)
		if matches == nil {
			continue
		}

		matches = matches[1:]
		return
	}

	return nil
}

func (rv *RhapValidator) required(rule string, validation Validation) (result bool) {
	value := validation.Value.String()
	if result = Required(value); !result {
		rv.setMessage(rule, ErrorRequiredMessage, validation)
		return
	}
	return
}

func (rv *RhapValidator) len(matches []string, validation Validation) (result bool) {
	value := validation.Value.String()
	length, err := strconv.Atoi(matches[1])
	if err != nil {
		panic(err)
	}

	if result = Len(value, length); !result {
		rv.setMessageWithRuleValue(matches, ErrorExactLengthMessage, validation)
		return
	}

	return
}

func (rv *RhapValidator) min(matches []string, validation Validation) (result bool) {
	value := validation.Value.String()
	length, err := strconv.Atoi(matches[1])
	if err != nil {
		panic(err)
	}

	if result = MinString(value, length); !result {
		rv.setMessageWithRuleValue(matches, ErrorMinLengthFormat, validation)
		return
	}

	return
}

func (rv *RhapValidator) max(matches []string, validation Validation) (result bool) {
	value := validation.Value.String()
	length, err := strconv.Atoi(matches[1])
	if err != nil {
		panic(err)
	}

	if result = MaxString(value, length); !result {
		rv.setMessageWithRuleValue(matches, ErrorMaxLengthMessage, validation)
		return
	}

	return
}

func (rv *RhapValidator) alpha(rule string, validation Validation) (result bool) {
	value := validation.Value.String()
	if result = IsAlpha(value); !result {
		rv.setMessage(rule, ErrorAlphaMessage, validation)
		return
	}
	return
}

func (rv *RhapValidator) email(rule string, validation Validation) (result bool) {
	value := validation.Value.String()
	if result = IsEmail(value); !result {
		rv.setMessage(rule, ErrorInvalidFormatMessage, validation)
		return
	}
	return
}

/*
	Integer validations
*/
func validateIntType() {

}

func (rv *RhapValidator) setMessage(rule string, message string, validation Validation) {
	customMessage := validation.CustomMessage
	if hasCustomMessage(rule, customMessage) {
		message = customMessage[rule]
	}

	rv.results[validation.FieldName] = Message{
		"message": strings.ReplaceAll(message, "${label}", setLabel(validation)),
		"format":  message,
	}
}

func (rv *RhapValidator) setMessageWithRuleValue(matches []string, message string, validation Validation) {
	var value string
	if len(matches) > 1 {
		value = matches[1]
	}

	customMessage := validation.CustomMessage
	if hasCustomMessage(matches[0], customMessage) {
		message = customMessage[matches[0]]
	}

	messageResult := strings.ReplaceAll(message, "${label}", setLabel(validation))
	messageResult = strings.ReplaceAll(messageResult, "${rule_value}", value)
	message = strings.ReplaceAll(message, "${rule_value}", value)

	rv.results[validation.FieldName] = Message{
		"message": messageResult,
		"format":  message,
	}
}

func hasCustomMessage(rule string, message Message) bool {
	if message == nil {
		return false
	}

	return message[rule] != ""
}

func setLabel(validation Validation) string {
	if validation.Label == "" {
		return validation.FieldName
	}
	return validation.Label
}