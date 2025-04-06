package headers

import (
	"errors"
	"regexp"
	"strings"
)

type Headers map[string]string

var CRLF string = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	headerString := string(data)
	crlf_idx := strings.Index(headerString, CRLF)

	if crlf_idx < 0 {
		return 0, false, nil
	}
	if crlf_idx == 0 {
		return 2, true, nil
	}

	header_str := headerString[:crlf_idx]

	colon_idx := strings.Index(header_str, ":")

	if string(header_str[(colon_idx-1)]) == " " {
		return 0, false, errors.New("malformed header")
	}

	header_str = strings.TrimSpace(header_str)
	header_slice := strings.SplitN(header_str, ":", 2)

	key := strings.TrimSpace(header_slice[0])
	value := strings.TrimSpace(header_slice[1])

	if len(key) < 1 {
		return 0, false,  errors.New("field does not contain enough characters")
	}

	pattern := `^[A-Za-z0-9!#$%&'*+\-.\^_` + `|~]+$`
	re := regexp.MustCompile(pattern)

	if !re.MatchString(key) {
		return 0, false, errors.New("field contains illegal characters")
	}

	if h[strings.ToLower(key)] != "" {
		newValueString := h[strings.ToLower(key)] + ", " + value
		h[strings.ToLower(key)] = newValueString

	} else {
		h[strings.ToLower(key)] = value 	
	}

	return crlf_idx + len(CRLF), false, nil
}

func (h Headers) Get(key string) (val string) {
	return h[strings.ToLower(key)]
}

func (h Headers) Set(key, val string) string {
	h[strings.ToLower(key)] = val
	return 	h[strings.ToLower(key)]
}