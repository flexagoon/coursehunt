package search

import "sync"

type Provider interface {
	Search(query string, filter Filter) ([]Course, error)
}

type Course struct {
	Name        string
	Url         string
	Description string
	Price       string
}

type Filter struct {
	Price      bool
	PriceRange [2]int
}

func Search(query string, filter Filter, providers []Provider) []Course {
	var wg sync.WaitGroup
	var results []Course
	var mu sync.Mutex

	for _, provider := range providers {
		wg.Add(1)
		go func(provider Provider, results *[]Course) {
			defer wg.Done()

			providerResults, _ := provider.Search(query, filter)

			mu.Lock()
			*results = append(*results, providerResults...)
			mu.Unlock()
		}(provider, &results)
	}

	wg.Wait()
	return results
}
