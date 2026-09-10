package fake

import "strconv"

func arg(args []any, index int) string {
	if index >= len(args) {
		return ""
	}

	return toString(args[index])
}

func globMatch(pattern, value string) bool {
	if pattern == "*" {
		return true
	}

	if pattern == "" {
		return value == ""
	}

	switch pattern[0] {
	case '\\':
		if len(pattern) > 1 && value != "" && value[0] == pattern[1] {
			return globMatch(pattern[2:], value[1:])
		}

		return false
	case '*':
		if globMatch(pattern[1:], value) {
			return true
		}

		return value != "" && globMatch(pattern, value[1:])
	case '?':
		return value != "" && globMatch(pattern[1:], value[1:])
	}

	return value != "" && value[0] == pattern[0] && globMatch(pattern[1:], value[1:])
}

func toString(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case bool:
		return strconv.FormatBool(v)
	case []byte:
		return string(v)
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	}

	return ""
}
