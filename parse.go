package pipe

import (
	"bufio"
	"io"
	"strings"
)

// parse will split string on spaces and also takes
// into account double and single quoted strings
func parse(in string) ([]string, error) {
	var ret = make([]string, 0)
	var buf = bufio.NewReader(strings.NewReader(in))
	for {
		var delim byte = ' '
		if next, _ := buf.Peek(1); len(next) == 1 {
			if  next[0] == '\'' || next[0] == '"' {
				delim = next[0]
				_, _ = buf.Discard(1)
			}
		}
		str, err := buf.ReadString(delim)
		if len(str) > 0 && str[len(str)-1] == delim {
			str = str[:len(str)-1]
		}
		if "" != str {
			ret = append(ret, str)
		}
		if err != nil {
			if err == io.EOF {
				return ret, nil
			} else {
				return nil, err
			}
		}
	}
}