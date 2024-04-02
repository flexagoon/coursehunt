package main

import (
	"encoding/json"
	"net/http"
	"net/url"

	"fxgn.dev/coursehunt/search"
	"fxgn.dev/coursehunt/search/providers"
	"fxgn.dev/coursehunt/views"
	"github.com/a-h/templ"
)

func main() {
	searchProviders := []search.Provider{
		providers.Stepik{},
		providers.Udemy{
			ClientId:     "TNAzpQDvOG5n86nNZsOdHb5TzmZsCvSt7segUL71",
			ClientSecret: "bFHRmQapiQk5D9oFTNchg8M7bwqgp6xw0o1Kv6yMTaJ5g7IiLGj3sZAaT1IR64WTItkx0ubRjRrQ0eEhSRE96C7VpbMqqC6C7xQuuxjdnuLu4VqqlrSrqvVqyEYxX6Zc",
		},
		providers.Coursera{},
		providers.Edx{},
		providers.Skillbox{},
		providers.Alison{},
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

	mux.HandleFunc("/sort", func(w http.ResponseWriter, r *http.Request) {
		sort := r.FormValue("sort")
		r.ParseForm()
		resultStrings, _ := r.Form["result[]"]

		results := []search.Course{}
		for _, resultString := range resultStrings {
			result := search.Course{}
			json.Unmarshal([]byte(resultString), &result)
			results = append(results, result)
		}

		results = search.SortCourses(results, sort)

		views.ResultsList(results).Render(r.Context(), w)
	})

	mux.HandleFunc("GET /style.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "views/style.css")
	})

	http.ListenAndServe(":8080", mux)
}

func serveHtmxPage(r *http.Request, w http.ResponseWriter, component templ.Component) {
	if r.Header.Get("HX-Request") != "true" {
		component = views.Page(component)
	}
	component.Render(r.Context(), w)
}
