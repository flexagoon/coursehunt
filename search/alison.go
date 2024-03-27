package search

import (
	"encoding/json"
	"net/http"
)

type Alison struct{}

type alisonCourse struct {
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Headline string `json:"headline"`
}

const alisonBaseUrl = `https://api.alison.com/v0.1/search?query=`

func (alison Alison) Search(query string, filter Filter) ([]Course, error) {
	if filter.Language == LanguageRussian {
		return []Course{}, nil
	}

	url := alisonBaseUrl + query

	httpResp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	response := struct {
		Result []alisonCourse `json:"result"`
	}{}
	err = json.NewDecoder(httpResp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	var courses []Course
	for _, course := range response.Result {
		courses = append(courses, Course{
			Name:        course.Name,
			Url:         "https://alison.com/course/" + course.Slug,
			Description: course.Headline,
			Price:       "Free",
		})
	}

	return courses, nil
}
