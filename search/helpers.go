package search

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
