package search

type Filter struct {
	Free     bool
	Language Language
}

type Language int

const (
	LanguageAny Language = iota
	LanguageEnglish
	LanguageRussian
)
