package search

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type Edx struct{}

type edxCourse struct {
	Title       string `json:"title"`
	Url         string `json:"marketing_url"`
	Description string `json:"primary_description"`
}

const baseUrl = `https://igsyv1z1xi-dsn.algolia.net/1/indexes/product?x-algolia-application-id=IGSYV1Z1XI&x-algolia-api-key=1f72394b5b49fc876026952685f5defe&filters=(product:Course+AND+language:English)&attributesToRetrieve=["title","marketing_url","primary_description"]&attributesToHighlight=[]&query=`

func (edx Edx) Search(query string, filter Filter) ([]Course, error) {
	if filter.Language == LanguageRussian {
		return []Course{}, nil
	}

	url := baseUrl + url.QueryEscape(query)

	httpResp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	response := struct {
		Hits []edxCourse `json:"hits"`
	}{}
	err = json.NewDecoder(httpResp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	var courses []Course
	for _, course := range response.Hits {
		courses = append(courses, Course{
			Name:        course.Title,
			Url:         course.Url,
			Description: course.Description,
			Price:       "Free, paid certificate",
		})
	}

	return courses, nil
}
