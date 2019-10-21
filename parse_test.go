package pipe

import (
	"testing"
)


func getTestdata() map[string][]string {
	return map[string][]string{
		"hello world": {"hello", "world"},
		"foo bar foo": {"foo", "bar", "foo"},
		`foo 'bar foo''`: {"foo", "bar foo"},
		`foo 'bar foo' foo`: {"foo", "bar foo", "foo"},
		`foo 'bar foo' foo "bar foo foo"`: {"foo", "bar foo", "foo", "bar foo foo"},
	}
}


func TestParse(t *testing.T) {
	for data, expected:= range getTestdata() {
		if result, err := parse(data); err != nil {
			t.Error(err)
		} else {
			if len(result) != len(expected) {
				t.Errorf("unsuspected result count returned, expected %d got %d", len(expected), len(result))
			}
			for i, c := 0, len(result); i < c; i++ {
				if result[i] != expected[i] {
					t.Errorf("unsuspected result returned, expected %s got %s", expected[i], result[i])
				}
			}
		}
	}
}

func BenchmarkParse(b *testing.B) {
	for data, _ := range getTestdata() {
		_, _ = parse(data)
	}
}