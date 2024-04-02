package search

import (
	"cmp"
	"slices"
)

func zipSources(sources [][]Course) []Course {
	sourcesCount := len(sources)

	longestSourceLen := 0
	for _, source := range sources {
		if len(source) > longestSourceLen {
			longestSourceLen = len(source)
		}
	}

	var result []Course
	for i := 0; i < longestSourceLen; i++ {
		for j := 0; j < sourcesCount; j++ {
			if len(sources[j]) > i {
				result = append(result, sources[j][i])
			}
		}
	}

	return result
}

func SortCourses(courses []Course, sort string) []Course {
	if sort == "rating" {
		slices.SortStableFunc(courses, func(a, b Course) int {
			return -cmp.Compare(a.Rating, b.Rating)
		})
	} else if sort == "cheap" || sort == "expensive" {
		courses := slices.DeleteFunc(courses, func(result Course) bool {
			return result.Price == 0
		})

		if sort == "cheap" {
			slices.SortStableFunc(courses, func(a, b Course) int {
				return cmp.Compare(a.Price, b.Price)
			})
		} else if sort == "expensive" {
			slices.SortStableFunc(courses, func(a, b Course) int {
				return -cmp.Compare(a.Price, b.Price)
			})
		}
	} else if sort == "short" || sort == "long" {
		courses := slices.DeleteFunc(courses, func(result Course) bool {
			return result.Hours == 0
		})

		if sort == "short" {
			slices.SortStableFunc(courses, func(a, b Course) int {
				return cmp.Compare(a.Hours, b.Hours)
			})
		} else if sort == "long" {
			slices.SortStableFunc(courses, func(a, b Course) int {
				return -cmp.Compare(a.Hours, b.Hours)
			})
		}
	}

	return courses
}
