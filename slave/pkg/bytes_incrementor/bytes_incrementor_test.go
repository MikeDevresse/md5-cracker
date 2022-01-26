package bytes_incrementor

import (
	"bytes"
	"testing"
)

func TestIncrement(t *testing.T) {
	got := Increment([]byte("a"), 0)
	if !bytes.Equal(got, []byte("b")) {
		t.Errorf("Increment([]byte(\"a\"), 0) = %v, want %v", got, []byte("b"))
	}

	got = Increment([]byte("z"), 0)
	if !bytes.Equal(got, []byte("A")) {
		t.Errorf("Increment([]byte(\"z\"), 0) = %v, want %v", got, []byte("A"))
	}

	got = Increment([]byte("Z"), 0)
	if !bytes.Equal(got, []byte("0")) {
		t.Errorf("Increment([]byte(\"Z\"), 0) = %v, want %v", got, []byte("0"))
	}

	got = Increment([]byte("a9"), 1)
	if !bytes.Equal(got, []byte("ba")) {
		t.Errorf("Increment([]byte(\"a9\"), 0) = %v, want %v", got, []byte("a9"))
	}

	got = Increment([]byte("8999"), 3)
	if !bytes.Equal(got, []byte("9aaa")) {
		t.Errorf("Increment([]byte(\"8999\"), 0) = %v, want %v", got, []byte("9aaa"))
	}

	got = Increment([]byte("9"), 0)
	if !bytes.Equal(got, []byte("aa")) {
		t.Errorf("Increment([]byte(\"9\"), 0) = %v, want %v", got, []byte("aa"))
	}
}
