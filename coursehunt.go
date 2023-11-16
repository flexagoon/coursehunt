package main

import (
	"fmt"

	"fxgn.dev/coursehunt/search"
)

func main() {
	searchProviders := []search.Provider{
		search.Stepik{},
		search.Udemy{
			ClientId:     "TNAzpQDvOG5n86nNZsOdHb5TzmZsCvSt7segUL71",
			ClientSecret: "bFHRmQapiQk5D9oFTNchg8M7bwqgp6xw0o1Kv6yMTaJ5g7IiLGj3sZAaT1IR64WTItkx0ubRjRrQ0eEhSRE96C7VpbMqqC6C7xQuuxjdnuLu4VqqlrSrqvVqyEYxX6Zc",
		},
	}

	results := search.Search("java", searchProviders)
	fmt.Println(results)
}
