package search

import (
	"encoding/json"
	"net/http"
)

type Stepik struct{}

type stepikCourse struct {
	Title string `json:"course_title"`
}

func (stepik Stepik) Search(query string) ([]Course, error) {
	url := "https://stepik.org/api/search-results?order=conversion_rate__none%2Crating__none%2Cquality__none%2Cpaid_weight__none&page=1&query=" + query
	httpResp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	response := struct {
		Results []stepikCourse `json:"search-results"`
	}{}
	err = json.NewDecoder(httpResp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	var courses []Course
	for _, course := range response.Results {
		courses = append(courses, Course{
			Name: course.Title,
		})
	}

	return courses, nil
}
