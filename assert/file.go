package assert

import (
	"os"
	"testing"
)

func Exists(t *testing.T, filename string) {
	t.Helper()
	info, err := os.Lstat(filename)
	if err != nil {
		t.Fatalf("err=%v\n", err)
	}

	if info.Size() == 0 {
		t.Fatal("size=0, want > 0")
	}
}
