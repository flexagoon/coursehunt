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
        "limit": 25,
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
	Free     bool     `json:"isCourseFree"`
	Rating   float64  `json:"avgProductRating"`
	Duration string   `json:"productDuration"`
}

var courseraDurations = map[string]int{
	"LESS_THAN_TWO_HOURS":  2,
	"ONE_TO_FOUR_WEEKS":    10,
	"ONE_TO_THREE_MONTHS":  2000,
	"THREE_TO_SIX_MONTHS":  4000,
	"SIX_TO_TWELVE_MONTHS": 7000,
	"ONE_TO_FOUR_YEARS":    24000,
}

func (coursera Coursera) Search(query string, filter search.Filter) ([]search.Course, error) {
	url := "https://www.coursera.org/graphql-gateway?opname=Search"
	payload := coursera.buildGraphqlPayload(query, filter)

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
		var price float64
		if course.Free {
			price = 0
		} else {
			if filter.Free {
				continue
			}
			price = -1
		}

		var extras []search.ExtraParam
		if filter.Language == search.LanguageRussian {
			extras = append(extras, search.Translated)
		}

		courses = append(courses, search.Course{
			Name:   course.Name,
			Author: strings.Join(course.Partners, ", "),
			Url:    "https://www.coursera.org" + course.Url,
			Price:  price,
			Rating: course.Rating,
			Hours:  courseraDurations[course.Duration],
			Extra:  extras,
		})
	}

	return courses, nil
}

func (_ Coursera) buildGraphqlPayload(query string, filter search.Filter) *strings.Reader {
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
