package utils

import (
	"strings"
	"unicode"
)

// cyrillicToLatin покрывает русский алфавит. Регистр не различаем — ниже всё
// приводится к нижнему регистру.
var cyrillicToLatin = map[rune]string{
	'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d",
	'е': "e", 'ё': "yo", 'ж': "zh", 'з': "z", 'и': "i",
	'й': "y", 'к': "k", 'л': "l", 'м': "m", 'н': "n",
	'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t",
	'у': "u", 'ф': "f", 'х': "kh", 'ц': "ts", 'ч': "ch",
	'ш': "sh", 'щ': "shch", 'ъ': "", 'ы': "y", 'ь': "",
	'э': "e", 'ю': "yu", 'я': "ya",
}

// Slug собирает URL-безопасный идентификатор: транслит кириллицы, нижний
// регистр, замена пробелов/подчёркиваний на дефис, схлопывание подряд идущих
// дефисов, обрезка по краям.
func Slug(s string) string {
	s = strings.ToLower(s)

	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch {
		case r == ' ' || r == '_' || r == '-':
			b.WriteRune('-')
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			if v, ok := cyrillicToLatin[r]; ok {
				b.WriteString(v)
			} else if r < 128 {
				b.WriteRune(r)
			}
			// non-ascii letters outside the cyrillic map are dropped
		default:
			// punctuation, symbols — separator
			b.WriteRune('-')
		}
	}

	out := collapseDashes(b.String())
	return strings.Trim(out, "-")
}

func collapseDashes(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	prevDash := false
	for _, r := range s {
		if r == '-' {
			if prevDash {
				continue
			}
			prevDash = true
		} else {
			prevDash = false
		}
		b.WriteRune(r)
	}
	return b.String()
}