package rhapvalidator

import (
	"regexp"
	"strconv"
	"strings"
)

var stringValidations = []*regexp.Regexp{
	regexp.MustCompile("^(required)$"),
	regexp.MustCompile("^(alpha)$"),
	regexp.MustCompile("^(email)$"),
	regexp.MustCompile("^(number)$"),
	regexp.MustCompile("^(len)=(\\d+)$"),
	regexp.MustCompile("^(min)=(\\d+)$"),
	regexp.MustCompile("^(max)=(\\d+)$"),
}

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