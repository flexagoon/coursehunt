package providers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"fxgn.dev/coursehunt/search"
)

type Stepik struct{}

type searchResult struct {
	Id int `json:"course"`
}

type stepikCourse struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Url         string `json:"canonical_url"`
	Summary     string `json:"summary"`
	Price       string `json:"price"`
	Seconds     int    `json:"time_to_complete"`
	Certificate string `json:"certificate"`
}

func (stepik Stepik) Search(query string, filter search.Filter) ([]search.Course, error) {
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

	return stepik.fetchCourses(ids)
}

func (_ Stepik) buildSearchUrl(query string, filter search.Filter) (string, error) {
	url, err := url.Parse("https://stepik.org/api/search-results?order=conversion_rate__none%2Crating__none%2Cquality__none%2Cpaid_weight__none&is_popular=true&page=1")
	if err != nil {
		return "", err
	}

	q := url.Query()

	q.Set("query", query)

	if filter.Free {
		q.Set("is_paid", "false")
	}

	if filter.Language == search.LanguageEnglish {
		q.Set("language", "en")
	} else if filter.Language == search.LanguageRussian {
		q.Set("language", "ru")
	}

	if filter.Difficulty == search.DifficultyBeginner {
		q.Set("difficulty[]", "easy")
	} else if filter.Difficulty == search.DifficultyIntermediate {
		q.Set("difficulty[]", "normal")
	} else if filter.Difficulty == search.DifficultyAdvanced {
		q.Set("difficulty[]", "hard")
	}

	url.RawQuery = q.Encode()

	return url.String(), nil
}

func (stepik Stepik) fetchCourses(ids []int) ([]search.Course, error) {
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

	var courses []search.Course
	for _, course := range response.Results {
		var price float64
		if course.Price != "" {
			price, _ = strconv.ParseFloat(course.Price, 64)
		} else {
			price = 0
		}

		var extras []search.ExtraParam
		if course.Certificate != "" {
			extras = append(extras, search.Certificate)
		}

		courses = append(courses, search.Course{
			Name:        course.Title,
			Url:         course.Url,
			Description: course.Summary,
			Price:       price,
			Rating:      course.fetchRating(),
			Hours:       course.Seconds / 3600,
			Extra:       extras,
		})
	}

	return courses, nil
}

type stepikRating struct {
	Rating float64 `json:"average"`
}

func (course stepikCourse) fetchRating() float64 {
	url := fmt.Sprintf("https://stepik.org/api/course-review-summaries?ids[]=%d", course.Id)

	httpResp, err := http.Get(url)
	if err != nil {
		return 0
	}

	response := struct {
		Result []stepikRating `json:"course-review-summaries"`
	}{}
	err = json.NewDecoder(httpResp.Body).Decode(&response)
	if err != nil || len(response.Result) == 0 {
		return 0
	}

	return response.Result[0].Rating
}
