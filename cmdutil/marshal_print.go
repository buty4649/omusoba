package cmdutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

func MarshalPrint(format string, value interface{}) error {
	s, err := marshal(format, value)
	if err != nil {
		return err
	}

	fmt.Printf("%s", s)
	return nil
}

func marshal(format string, value interface{}) (string, error) {
	switch format {
	case "yaml":
		m, err := yaml.Marshal(value)
		if err != nil {
			return "", err
		}
		return string(m), nil

	case "json":
		m, err := json.MarshalIndent(value, "", "  ")
		if err != nil {
			return "", err
		}
		return string(m) + "\n", nil

	case "label":
		t := reflect.TypeOf(value)
		v := reflect.ValueOf(value)

		maxLength := 0
		for i := 0; i < t.NumField(); i++ {
			desc := t.Field(i).Tag.Get("desc")
			if l := len(desc); l >= maxLength {
				maxLength = l
			}
		}

		result := ""
		for i := 0; i < t.NumField(); i++ {
			desc := t.Field(i).Tag.Get("desc")
			result += fmt.Sprintf(fmt.Sprintf("%%%ds", maxLength), desc)

			f := v.Field(i)
			switch f.Kind() {
			case reflect.Float32:
				s := fmt.Sprintf("%.3f", f.Float())
				result += fmt.Sprintf(": %s", strings.TrimRight(strings.TrimRight(s, "0"), "."))

			case reflect.Slice:
				result += fmt.Sprint(":")
				s := reflect.ValueOf(f.Interface())
				for i := 0; i < s.Len(); i++ {
					result += fmt.Sprintf(" %02x", s.Index(i))
				}

			default:
				result += fmt.Sprintf(": %v", f.Interface())
			}

			unit := t.Field(i).Tag.Get("unit")
			if unit != "" {
				result += fmt.Sprintf(" %s", unit)
			}

			result += "\n"
		}
		return result, nil
	}

	return "", errors.New("unknwon format")
}
