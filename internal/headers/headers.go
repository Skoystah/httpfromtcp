package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return map[string]string{}
}

const crlf = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))

	if idx == -1 {
		return 0, false, nil
	}

	parsedBytes := idx + 2

	if idx == 0 {
		return parsedBytes, true, nil
	}

	fieldLine := string(data[:idx])
	key, value, err := headerFromString(fieldLine)
	if err != nil {
		return 0, false, err
	}
	h.Set(key, value)

	return parsedBytes, false, nil
}

func headerFromString(str string) (string, string, error) {

	fieldLineParts := strings.SplitN(str, ":", 2)
	if len(fieldLineParts) != 2 {
		return "", "", fmt.Errorf("Field line does not consist of 2 parts: %s", str)
	}

	key := strings.TrimLeft(fieldLineParts[0], " ")

	if last := key[len(key)-1:]; last == " " || last == "\t" {
		return "", "", fmt.Errorf("Header key cannot have whitespace before ':' : %s", key)
	}

	b := []byte(key)
	for _, char := range b {
		if (char >= 'A' && char <= 'Z') ||
			(char >= 'a' && char <= 'z') ||
			(char >= '0' && char <= '9') ||
			(strings.IndexByte("!#$%&'*+-.^_`|~", char) >= 0) {
			continue
		} else {
			fmt.Println(char)
			return "", "", fmt.Errorf("Header key contains illegal character: %q %q", b, char)
		}
	}

	value := strings.Trim(fieldLineParts[1], " ")

	return key, value, nil
}

func (h Headers) Get(key string) (string, bool) {
	key = strings.ToLower(key)

	value, exists := h[key]
	return value, exists
}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)

	if current, exists := h[key]; exists {
		h[key] = fmt.Sprintf("%s, %s", current, value)
	} else {
		h[key] = value
	}
}

func (h Headers) Delete(key string) {
	key = strings.ToLower(key)
	delete(h, key)
}
func (h Headers) Update(key, value string) {
	key = strings.ToLower(key)

	h[key] = value
}
