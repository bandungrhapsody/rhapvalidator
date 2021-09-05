package rhapvalidator

import (
	"regexp"
	"strconv"
	"strings"
)

var intValidations = []*regexp.Regexp{
	regexp.MustCompile("^(required)$"),
	regexp.MustCompile("^(min)=(\\d+)$"),
	regexp.MustCompile("^(max)=(\\d+)$"),
}

func (rv *RhapValidator) validateInt(validation Validation) {
	slcRules := strings.Split(validation.Rules, ",")

	for _, rule := range slcRules {
		rule = strings.Trim(rule, " ")
		matches := matchIntRules(rule)
		if matches == nil {
			continue
		}

		if !rv.checkRuleAndValidateInt(matches, validation) {
			return
		}
	}
}

func (rv *RhapValidator) checkRuleAndValidateInt(matches []string, validation Validation) bool {
	r := matches[0]

	switch r {
	case "required":
		return rv.requiredNum(r, validation)
	case "min":
		return rv.minNum(matches, validation)
	case "max":
		return rv.maxNum(matches, validation)
	}

	return false
}

func matchIntRules(ruleName string) (matches []string) {
	for _, r := range intValidations {
		matches = r.FindStringSubmatch(ruleName)
		if matches == nil {
			continue
		}

		matches = matches[1:]
		return
	}

	return nil
}

func (rv *RhapValidator) requiredNum(rule string, validation Validation) (result bool) {
	value := validation.Value.Int()
	if result = RequiredNum(int(value)); !result {
		rv.setMessage(rule, ErrorRequiredMessage, validation)
		return
	}
	return
}

func (rv *RhapValidator) minNum(matches []string, validation Validation) (result bool) {
	value := validation.Value.Int()
	length, err := strconv.Atoi(matches[1])
	if err != nil {
		panic(err)
	}

	if result = MinNum(int(value), length); !result {
		rv.setMessageWithRuleValue(matches, ErrorMinValueMessage, validation)
		return
	}

	return
}

func (rv *RhapValidator) maxNum(matches []string, validation Validation) (result bool) {
	value := validation.Value.Int()
	length, err := strconv.Atoi(matches[1])
	if err != nil {
		panic(err)
	}

	if result = MaxNum(int(value), length); !result {
		rv.setMessageWithRuleValue(matches, ErrorMaxValueMessage, validation)
		return
	}

	return
}