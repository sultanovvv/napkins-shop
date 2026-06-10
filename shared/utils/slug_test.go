package utils

import "testing"

func TestSlug(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"Привет, мир!", "privet-mir"},
		{"Салфетка с розами", "salfetka-s-rozami"},
		{"Hello   World", "hello-world"},
		{"snake_case_text", "snake-case-text"},
		{"Уже-готовый-slug", "uzhe-gotovyy-slug"},
		{"  Trim  ", "trim"},
		{"Mixed Тест 42", "mixed-test-42"},
		{"---", ""},
	}
	for _, c := range cases {
		if got := Slug(c.in); got != c.want {
			t.Errorf("Slug(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}