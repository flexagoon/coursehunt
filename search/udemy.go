package search

import (
	"encoding/json"
	"net/http"
)

type Udemy struct {
	ClientId     string
	ClientSecret string
}

type udemyCourse struct {
	Title    string `json:"title"`
	UrlPart  string `json:"url"`
	Headline string `json:"headline"`
	Price    string `json:"price"`
}

func (udemy Udemy) Search(query string) ([]Course, error) {
	url := "https://www.udemy.com/api-2.0/courses/?search=" + query

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(udemy.ClientId, udemy.ClientSecret)

	httpResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	response := struct {
		Results []udemyCourse `json:"results"`
	}{}
	err = json.NewDecoder(httpResp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	var courses []Course
	for _, course := range response.Results {
		courses = append(courses, Course{
			Name:        course.Title,
			Url:         "https://www.udemy.com" + course.UrlPart,
			Description: course.Headline,
			Price:       course.Price,
		})
	}

	return courses, nil
}
