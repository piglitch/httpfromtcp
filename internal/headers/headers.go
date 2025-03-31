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
	if len(headerString) > 3 {
		if headerString[:2] == CRLF {
			return crlf_idx + len(CRLF), true, nil
		}
	}

	header_str := headerString[:crlf_idx]

	headerSlice := strings.Split(header_str, CRLF)

	for _, line := range headerSlice {
		
		if line == "" {
			continue
		}
		h_slice := strings.Fields(line)
		if len(h_slice) != 2 {
			return 0, false, errors.New("invalid spacing header" + line)
		}

		lines := strings.Fields(line)

		key := strings.TrimSpace(lines[0][:len(lines[0])-1])
		value := strings.TrimSpace(lines[1])

		if len(key) < 1 {
			return 0, false,  errors.New("field does not contain enough characters")
		}

		pattern := `^[A-Za-z0-9!#$%&'*+\-.\^_` + `|~]+$`
		re := regexp.MustCompile(pattern)

		if !re.MatchString(key) {
			return 0, false, errors.New("field contains illegal characters")
		}
		h[strings.ToLower(key)] = value 	
	}

	return crlf_idx + len(CRLF), false, nil
}
