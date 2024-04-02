package providers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"fxgn.dev/coursehunt/search"
)

type Alison struct{}

type alisonCourse struct {
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Headline string `json:"headline"`
	Duration string `json:"avg_duration"`
}

func (alison Alison) Search(query string, filter search.Filter) ([]search.Course, error) {
	if filter.Language == search.LanguageRussian {
		return []search.Course{}, nil
	}

	url, err := alison.buildSearchUrl(query, filter)
	if err != nil {
		return nil, err
	}

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

	var courses []search.Course
	for _, course := range response.Result {
		hours, _ := strconv.Atoi(strings.Split(course.Duration, "-")[0])

		courses = append(courses, search.Course{
			Name:        course.Name,
			Url:         "https://alison.com/course/" + course.Slug,
			Description: course.Headline,
			Price:       "Free",
			Hours:       hours,
			Extra:       []search.ExtraParam{search.Certificate},
		})
	}

	return courses, nil
}

const alisonBaseUrl = "https://api.alison.com/v0.1/search"

func (alison Alison) buildSearchUrl(query string, filter search.Filter) (string, error) {
	url, err := url.Parse(alisonBaseUrl)
	if err != nil {
		return "", err
	}

	q := url.Query()

	q.Set("query", query)

	if filter.Difficulty != search.DifficultyAny {
		q.Set("level", fmt.Sprint(filter.Difficulty))
	}

	url.RawQuery = q.Encode()

	return url.String(), nil
}
