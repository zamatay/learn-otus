package logger

import "testing"

func TestLogger(t *testing.T) {
	logger := New("info")
	if logger == nil {
		t.Fail()
	}
}
