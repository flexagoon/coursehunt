package search

type Filter struct {
	Free       bool
	Language   Language
	Difficulty Difficulty
}

type Language int

const (
	LanguageAny Language = iota
	LanguageEnglish
	LanguageRussian
)

type Difficulty int

const (
	DifficultyAny Difficulty = iota
	DifficultyBeginner
	DifficultyIntermediate
	DifficultyAdvanced
)
