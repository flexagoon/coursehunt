package search

import "sync"

type Provider interface {
	Search(query string) ([]Course, error)
}

type Course struct {
	Name string
}

func Search(query string, providers []Provider) []Course {
	var wg sync.WaitGroup
	var results []Course
	var mu sync.Mutex

	for _, provider := range providers {
		wg.Add(1)
		go func(provider Provider, results *[]Course) {
			defer wg.Done()

			providerResults, _ := provider.Search(query)

			mu.Lock()
			*results = append(*results, providerResults...)
			mu.Unlock()
		}(provider, &results)
	}

	wg.Wait()
	return results
}
