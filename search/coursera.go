package search

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
  "query": "query Search($requests: [Search_Request!]!) {SearchResult {search(requests: $requests) {elements {... on Search_ProductHit { name isCourseFree url skills productDuration productDifficultyLevel imageUrl avgProductRating productType}}}}}",
  "variables": {
    "requests": [
      {
        "entityType": "PRODUCTS",
        "query": "%s",
        "limit": 100,
        "sortBy": "BEST_MATCH"
      }
    ]
  }
}
`

type courseraCourse struct {
	Name       string   `json:"name"`
	Free       bool     `json:"isCourseFree"`
	Url        string   `json:"url"`
	Skills     []string `json:"skills"`
	Duration   string   `json:"productDuration"`
	Difficulty string   `json:"productDifficultyLevel"`
	ImageUrl   string   `json:"imageUrl"`
	Rating     float64  `json:"avgProductRating"`
	Type       string   `json:"productType"`
}

func (coursera Coursera) Search(query string, filter Filter) ([]Course, error) {
	url := "https://www.coursera.org/graphql-gateway?opname=Search"
	payload := strings.NewReader(fmt.Sprintf(basePayload, query))

	httpResp, err := http.Post(url, "application/json", payload)
	if err != nil {
		return nil, err
	}

	response := gqlResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	var courses []Course
	for _, course := range response.Data.SearchResult.Search[0].Courses {
		var price string
		if course.Free {
			price = "Free"
		} else {
			if filter.Price && filter.PriceRange[1] == 0 {
				continue
			}
			price = "Subscription required"
		}

		courses = append(courses, Course{
			Name:        course.Name,
			Url:         "https://www.coursera.org" + course.Url,
			Description: "",
			Price:       price,
		})
	}

	return courses, nil
}
