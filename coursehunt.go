package main

import (
	"context"
	"net/http"
	"net/url"

	"fxgn.dev/coursehunt/search"
	"fxgn.dev/coursehunt/views"
	"github.com/a-h/templ"
)

func main() {
	searchProviders := []search.Provider{
		search.Stepik{},
		search.Udemy{
			ClientId:     "TNAzpQDvOG5n86nNZsOdHb5TzmZsCvSt7segUL71",
			ClientSecret: "bFHRmQapiQk5D9oFTNchg8M7bwqgp6xw0o1Kv6yMTaJ5g7IiLGj3sZAaT1IR64WTItkx0ubRjRrQ0eEhSRE96C7VpbMqqC6C7xQuuxjdnuLu4VqqlrSrqvVqyEYxX6Zc",
		},
		search.Coursera{},
		search.Edx{},
		search.Skillbox{},
		search.Alison{},
	}

	indexPage := views.IndexPage()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		serveHtmxPage(r, w, indexPage)
	})

	mux.HandleFunc("GET /search", func(w http.ResponseWriter, r *http.Request) {
		query := url.QueryEscape(r.FormValue("q"))

		filter := search.Filter{}

		free := r.FormValue("free")
		if free == "on" {
			filter.Free = true
		}

		language := r.FormValue("language")
		if language == "russian" {
			filter.Language = search.LanguageRussian
		} else if language == "english" {
			filter.Language = search.LanguageEnglish
		}

		difficulty := r.FormValue("difficulty")
		if difficulty == "beginner" {
			filter.Difficulty = search.DifficultyBeginner
		} else if difficulty == "intermediate" {
			filter.Difficulty = search.DifficultyIntermediate
		} else if difficulty == "advanced" {
			filter.Difficulty = search.DifficultyAdvanced
		}

		results := search.Search(query, filter, searchProviders)
		serveHtmxPage(r, w, views.SearchPage(results))
	})

	mux.HandleFunc("GET /style.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "views/style.css")
	})

	http.ListenAndServe(":1641", mux)
}

func serveHtmxPage(r *http.Request, w http.ResponseWriter, component templ.Component) {
	if r.Header.Get("HX-Request") != "true" {
		component = views.Page(component)
	}
	component.Render(context.Background(), w)
}
