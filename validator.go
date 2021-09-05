package rhapvalidator

import (
	"reflect"
	"strings"
)

type RhapValidator struct {
	results        ValidationMessages
	customMessages ValidationMessages
}

type ValidationMessages map[string]Message
type Message map[string]string

type CustomValidationHandler func(setMessage MessageSetter)
type MessageSetter func(message string)

type Validation struct {
	Rules         string
	FieldName     string
	Type          reflect.Kind
	Label         string
	Value         reflect.Value
	CustomMessage Message
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

func (rv *RhapValidator) CustomValidation(fieldName string, label string, validation CustomValidationHandler) {
	setter := func(message string) {
		if rv.results[fieldName] != nil {
			return
		}

		if label == "" {
			label = fieldName
		}

		rv.results[fieldName] = Message{
			"message": strings.ReplaceAll(message, "${label}", label),
			"format":  message,
		}
	}

	validation(setter)
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
			rv.validateInt(validation)
			break
		case reflect.Int32:
			rv.validateInt(validation)
			break
		case reflect.Int64:
			rv.validateInt(validation)
			break
	}
}

func (rv *RhapValidator) setMessage(rule string, message string, validation Validation) {
	if rv.results[validation.FieldName] != nil {
		return
	}

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
	if rv.results[validation.FieldName] != nil {
		return
	}

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