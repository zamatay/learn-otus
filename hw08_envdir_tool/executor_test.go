package main

import "testing"

var sl = []string{
	"/bin/bash",
	"./testdata/echo.sh",
	"arg1=1",
	"arg2=2",
}

var envTest = Environment{
	"BAR":   EnvValue{Value: "bar"},
	"EMPTY": EnvValue{NeedRemove: true},
	"FOO":   EnvValue{Value: "   foo\nwith new line"},
	"HELLO": EnvValue{Value: "\"hello\""},
	"UNSET": EnvValue{NeedRemove: true},
}

func TestRunCmd(t *testing.T) {
	t.Run("execTest", func(t *testing.T) {
		resultCode := RunCmd(sl, envTest)
		if resultCode != 0 {
			t.Fail()
		}
	})
}
