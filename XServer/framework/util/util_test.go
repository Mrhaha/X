package util

import (
	"strings"
	"testing"
)

// Test_Encoding ...
func Test_Encoding(t *testing.T) {
	r := UTF8ToGBK(strings.NewReader("你好 GBK"))
	t.Log(r)
}

// Test_Dir ...
func Test_Dir(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Error(e)
		}
	}()

	execDir, err := GetCurrExecDir()
	if err != nil {
		t.Error(err)
		return
	}

	if !IsDirExists(execDir) {
		t.Error(err)
	}
}
