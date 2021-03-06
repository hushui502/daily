package v2h

import "testing"

func TestFormatParse(t *testing.T) {
	for i, tt := range []struct {
		in  uint64
		out string
	}{
		{in: 0, out: "00"},
		{in: 1, out: "01"},
		{in: 2, out: "02"},
		{in: 15, out: "0f"},
		{in: 16, out: "110"},
		{in: 288, out: "2120"},
	} {
		s := Format(tt.in)
		if got, want := s, tt.out; got != want {
			t.Errorf("%d. Format()=%v, want %v", i, got, want)
		}

		if v, err := Parse(s); err != nil {
			t.Fatal(err)
		} else if got, want := v, tt.in; got != want {
			t.Errorf("%d. Parse(%q)=%v, want %v", i, s, got, want)
		}
	}
}

func TestParse(t *testing.T) {
	t.Run("ErrBlank", func(t *testing.T) {
		if _, err := Parse(""); err != nil || err.Error() != `Parse(): value must have length of at least 2` {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("ErrShort", func(t *testing.T) {
		if _, err := Parse("0"); err == nil || err.Error() != `Parse(): value must have length of at least 2` {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("ErrInvalidPrefix", func(t *testing.T) {
		if _, err := Parse("x0"); err == nil || err.Error() != `Parse(): invalid prefix` {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("ErrInvalidHex", func(t *testing.T) {
		if _, err := Parse("0djks"); err == nil || err.Error() != `Parse(): invalid hex value` {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}