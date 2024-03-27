package search

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Edx struct{}

type edxCourse struct {
	Title       string `json:"title"`
	Url         string `json:"marketing_url"`
	Description string `json:"primary_description"`
}

func (edx Edx) Search(query string, filter Filter) ([]Course, error) {
	if filter.Language == LanguageRussian {
		return []Course{}, nil
	}

	url := edx.buildSearchUrl(query, filter)

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

const edxBaseUrl = `https://igsyv1z1xi-dsn.algolia.net/1/indexes/product?x-algolia-application-id=IGSYV1Z1XI&x-algolia-api-key=1f72394b5b49fc876026952685f5defe&filters=%s&attributesToRetrieve=["title","marketing_url","primary_description"]&attributesToHighlight=[]&query=%s`

func (edx Edx) buildSearchUrl(query string, filter Filter) string {
	var filterSb strings.Builder
	filterSb.WriteString("(product:Course+AND+language:English")

	if filter.Difficulty == DifficultyBeginner {
		filterSb.WriteString("+AND+level:Introductory")
	} else if filter.Difficulty == DifficultyIntermediate {
		filterSb.WriteString("+AND+level:Intermediate")
	} else if filter.Difficulty == DifficultyAdvanced {
		filterSb.WriteString("+AND+level:Advanced")
	}

	filterSb.WriteString(")")

	return fmt.Sprintf(edxBaseUrl, filterSb.String(), query)
}
