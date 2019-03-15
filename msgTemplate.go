package seq

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var reg, _ = regexp.Compile("({.+?})")

/* implement part of the message template protocol see https://messagetemplates.org */

// render message template with a parameters map
func RenderMsgTemplate(mt string, params map[string]string) string {
	return reg.ReplaceAllStringFunc(mt, func(s string) string {
		k := s[1:(len(s) - 1)]
		if strings.HasPrefix(k, "@") {
			k = k[1:]
			value := params[k]
			if value == "" {
				return s
			} else {
				kv := make(map[string]interface{})
				err := json.Unmarshal([]byte(value), &kv)
				if err != nil {
					fmt.Println(err)
					return s
				}
				frags := make([]string, 0)
				for k, v := range kv {
					frags = append(frags, fmt.Sprintf("%v: %v", k, v))
				}
				return "{" + strings.Join(frags, ", ") + "}"
			}
		}
		value := params[k]
		if value == "" {
			return s
		} else {
			return value
		}
	})
}

func ExtractParams(mt string, v ...interface{}) (map[string]string, error) {
	if len(v) == 0 {
		return make(map[string]string), nil
	}
	keys := reg.FindAllString(mt, -1)
	if len(keys) != len(v) {
		return nil, errors.New("fields and params count not match")
	}
	result := make(map[string]string)
	for i := range keys {
		k := keys[i]
		k = k[1 : len(k)-1]
		if strings.HasPrefix(k, "@") {
			k = k[1:]
			d, err := json.Marshal(v[i])
			if err != nil {
				return nil, err
			}
			result[k] = string(d)
		} else {
			result[k] = fmt.Sprint(v[i])
		}
	}
	return result, nil
}
