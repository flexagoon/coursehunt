package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"

	"fxgn.dev/coursehunt/search"
	"github.com/go-chi/chi/v5"
)

func main() {
	searchProviders := []search.Provider{
		search.Stepik{},
		search.Udemy{
			ClientId:     "TNAzpQDvOG5n86nNZsOdHb5TzmZsCvSt7segUL71",
			ClientSecret: "bFHRmQapiQk5D9oFTNchg8M7bwqgp6xw0o1Kv6yMTaJ5g7IiLGj3sZAaT1IR64WTItkx0ubRjRrQ0eEhSRE96C7VpbMqqC6C7xQuuxjdnuLu4VqqlrSrqvVqyEYxX6Zc",
		},
	}

	indexPage, _ := makeInheritedTemplate("views/index.html")
	searchPage, _ := makeInheritedTemplate("views/search.html")

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		serveHtmx(r, w, indexPage, nil)
	})

	r.Get("/search", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		fmt.Println(query)
		results := search.Search(query, searchProviders)
		fmt.Println(results)
		serveHtmx(r, w, searchPage, results)
	})

	http.ListenAndServe(":1641", r)
}

func makeInheritedTemplate(file string) (*template.Template, error) {
	return template.New("").ParseFiles(file, "views/base.html")
}

func serveHtmx(r *http.Request, w http.ResponseWriter, tmpl *template.Template, data any) {
	if r.Header.Get("HX-Request") == "true" {
		executeInheritedTemplateContent(w, tmpl, data)
	} else {
		executeInheritedTemplate(w, tmpl, data)
	}
}

func executeInheritedTemplate(wr io.Writer, tmpl *template.Template, data any) {
	tmpl.ExecuteTemplate(wr, "base", data)
}

func executeInheritedTemplateContent(wr io.Writer, tmpl *template.Template, data any) {
	tmpl.ExecuteTemplate(wr, "content", data)
}
