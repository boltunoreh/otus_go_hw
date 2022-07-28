package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadDir(t *testing.T) {
	expectedEnv := Environment{
		"BAR":   EnvValue{"bar", false},
		"EMPTY": EnvValue{"", false},
		"FOO":   EnvValue{"   foo\nwith new line", false},
		"HELLO": EnvValue{"\"hello\"", false},
		"UNSET": EnvValue{"", true},
	}
	env, err := ReadDir("testdata/env")

	assert.NoError(t, err)
	assert.Equal(t, expectedEnv, env)
}
