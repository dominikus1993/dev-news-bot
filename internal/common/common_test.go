package common

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseIntWhenEmpty(t *testing.T) {
	str := ""
	subject := ParseInt(str, 1)
	assert.Equal(t, 1, subject)
}

func TestParseIntWhenIncorrectString(t *testing.T) {
	str := "ddsdds"
	subject := ParseInt(str, 1)
	assert.Equal(t, 1, subject)
}

func TestParseInt(t *testing.T) {
	str := "3"
	subject := ParseInt(str, 1)
	assert.Equal(t, 3, subject)
}

func TestParseIntWhenNegative(t *testing.T) {
	str := "-3"
	subject := ParseInt(str, 1)
	assert.Equal(t, -3, subject)
}

func FuzzParseInt(f *testing.F) {
	f.Add("3", true)
	f.Add("3.0", false)
	f.Add("xDD", false)
	f.Add("0", true)
	f.Add("-3", true)

	f.Fuzz(func(t *testing.T, str string, expected bool) {
		subject := ParseInt(str, 1)
		if expected {
			assert.NotEqual(t, 1, subject)
		} else {
			assert.Equal(t, 1, subject)
		}
	})
}

func TestGetEnvOrDefault(t *testing.T) {
	env := "SOME_RANDOM_ENV_VAR_FOR_TESTS"
	os.Setenv(env, "value")
	subject := GetEnvOrDefault(env, "default")
	assert.Equal(t, "value", subject)
}
