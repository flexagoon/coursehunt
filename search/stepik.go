package search

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Stepik struct{}

type searchResult struct {
	Id int `json:"course"`
}

type stepikCourse struct {
	Title   string `json:"title"`
	Url     string `json:"canonical_url"`
	Summary string `json:"summary"`
}

func (stepik Stepik) Search(query string) ([]Course, error) {
	url := "https://stepik.org/api/search-results?order=conversion_rate__none%2Crating__none%2Cquality__none%2Cpaid_weight__none&page=1&query=" + query
	httpResp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	response := struct {
		Results []searchResult `json:"search-results"`
	}{}
	err = json.NewDecoder(httpResp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	var ids []int
	for _, result := range response.Results {
		ids = append(ids, result.Id)
	}

	return fetchCourses(ids)
}

func fetchCourses(ids []int) ([]Course, error) {
	url := "https://stepik.org/api/courses?ids[]=" + strconv.Itoa(ids[0])
	for _, id := range ids[1:] {
		url += "&ids[]=" + strconv.Itoa(id)
	}

	httpResp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	response := struct {
		Results []stepikCourse `json:"courses"`
	}{}
	err = json.NewDecoder(httpResp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	var courses []Course
	for _, course := range response.Results {
		courses = append(courses, Course{
			Name:        course.Title,
			Url:         course.Url,
			Description: course.Summary,
		})
	}

	return courses, nil
}
