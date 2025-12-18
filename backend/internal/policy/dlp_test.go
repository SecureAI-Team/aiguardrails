package policy

import "testing"

func TestDetectDLPBuiltinAndCustom(t *testing.T) {
	text := "user password is secret"
	res := DetectDLP(text, []string{})
	if !res.Hit {
		t.Fatalf("expected hit for built-in keyword")
	}
	customRes := DetectDLP("contains foobar123", []string{"foobar"})
	if !customRes.Hit {
		t.Fatalf("expected hit for custom term")
	}
}

