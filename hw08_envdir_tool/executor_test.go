package main

import "testing"

var sl = []string{
	"/bin/bash",
	"./testdata/echo.sh",
	"arg1=1",
	"arg2=2",
}

func TestRunCmd(t *testing.T) {
	t.Run("execTest", func(t *testing.T) {
		resultCode := RunCmd(sl, envTest)
		if resultCode != 0 {
			t.Fail()
		}
	})
}
