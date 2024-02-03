package search

import (
	"encoding/json"
	"net/http"
	"net/url"
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
	Price   string `json:"price"`
}

func (stepik Stepik) Search(query string, filter Filter) ([]Course, error) {
	url, err := stepik.buildSearchUrl(query, filter)
	if err != nil {
		return nil, err
	}

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

func (_ Stepik) buildSearchUrl(query string, filter Filter) (string, error) {
	url, err := url.Parse("https://stepik.org/api/search-results?order=conversion_rate__none%2Crating__none%2Cquality__none%2Cpaid_weight__none&page=1")
	if err != nil {
		return "", err
	}

	q := url.Query()

	q.Set("query", query)

	if filter.Free {
		q.Set("is_paid", "false")
	}

	if filter.Language == LanguageEnglish {
		q.Set("language", "en")
	} else if filter.Language == LanguageRussian {
		q.Set("language", "ru")
	}

	url.RawQuery = q.Encode()

	return url.String(), nil
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
		var price string
		if course.Price != "" {
			price = course.Price + "₽"
		} else {
			price = "Free"
		}

		courses = append(courses, Course{
			Name:        course.Title,
			Url:         course.Url,
			Description: course.Summary,
			Price:       price,
		})
	}

	return courses, nil
}
