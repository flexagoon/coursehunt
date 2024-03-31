package providers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"fxgn.dev/coursehunt/search"
)

type Coursera struct{}

type gqlResponse struct {
	Data gqlData `json:"data"`
}

type gqlData struct {
	SearchResult gqlSearchResult `json:"SearchResult"`
}

type gqlSearchResult struct {
	Search []gqlSearch `json:"search"`
}

type gqlSearch struct {
	Courses []courseraCourse `json:"elements"`
}

const basePayload = `
{
  "query": "query Search($requests: [Search_Request!]!) {SearchResult {search(requests: $requests) {elements {... on Search_ProductHit { name partners url isCourseFree avgProductRating productDuration }}}}}",
  "variables": {
    "requests": [
      {
        "entityType": "PRODUCTS",
        "query": "%s",
        "limit": 100,
        "sortBy": "BEST_MATCH",
        "facetFilters": [ %s ]
      }
    ]
  }
}
`

type courseraCourse struct {
	Name     string   `json:"name"`
	Partners []string `json:"partners"`
	Url      string   `json:"url"`
	// TODO description
	Free     bool    `json:"isCourseFree"`
	Rating   float64 `json:"avgProductRating"`
	Duration string  `json:"productDuration"`
}

var courseraDurations = map[string]string{
	"LESS_THAN_TWO_HOURS":  "Less than 2 hours",
	"ONE_TO_FOUR_WEEKS":    "1-4 weeks",
	"ONE_TO_THREE_MONTHS":  "1-3 months",
	"THREE_TO_SIX_MONTHS":  "3-6 months",
	"SIX_TO_TWELVE_MONTHS": "6-12 months",
	"ONE_TO_FOUR_YEARS":    "1-4 years",
}

func (coursera Coursera) Search(query string, filter search.Filter) ([]search.Course, error) {
	url := "https://www.coursera.org/graphql-gateway?opname=Search"
	payload := buildGraphqlPayload(query, filter)

	httpResp, err := http.Post(url, "application/json", payload)
	if err != nil {
		return nil, err
	}

	response := gqlResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	var courses []search.Course
	for _, course := range response.Data.SearchResult.Search[0].Courses {
		var price string
		if course.Free {
			price = "Free"
		} else {
			if filter.Free {
				continue
			}
			price = "Subscription required"
		}

		var extras []search.ExtraParam
		if filter.Language == search.LanguageRussian {
			extras = append(extras, search.Translated)
		}

		courses = append(courses, search.Course{
			Name:     course.Name,
			Author:   strings.Join(course.Partners, ", "),
			Url:      "https://www.coursera.org" + course.Url,
			Price:    price,
			Rating:   course.Rating,
			Duration: courseraDurations[course.Duration],
			Extra:    extras,
		})
	}

	return courses, nil
}

func buildGraphqlPayload(query string, filter search.Filter) *strings.Reader {
	filters := []string{}

	if filter.Language == search.LanguageEnglish {
		filters = append(filters, "language:English")
	} else if filter.Language == search.LanguageRussian {
		filters = append(filters, "language:Russian")
	}

	if filter.Difficulty == search.DifficultyBeginner {
		filters = append(filters, "productDifficultyLevel:Beginner")
	} else if filter.Difficulty == search.DifficultyIntermediate {
		filters = append(filters, "productDifficultyLevel:Intermediate")
	} else if filter.Difficulty == search.DifficultyAdvanced {
		filters = append(filters, "productDifficultyLevel:Advanced")
	}

	filtersJson, _ := json.Marshal(filters)

	return strings.NewReader(fmt.Sprintf(basePayload, query, filtersJson))
}
