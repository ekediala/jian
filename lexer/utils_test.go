package lexer

import "testing"

func TestIsLetter(t *testing.T) {
	t.Run("correctly identifies a letter", func(t *testing.T) {
		var b byte = 'c'
		res := isLetter(b)
		if !res {
			t.Errorf("expected %q to be letter", string(b))
		}
	})

	t.Run("correctly identifies a non letter", func(t *testing.T) {
		var b byte = '='
		res := isLetter(b)
		if res {
			t.Errorf("expected %q to not be letter", string(b))
		}
	})
}
