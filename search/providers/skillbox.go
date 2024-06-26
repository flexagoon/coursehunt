package providers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"fxgn.dev/coursehunt/search"
)

type Skillbox struct{}

type skillboxCourse struct {
	Title    string           `json:"title"`
	Href     string           `json:"href"`
	Terms    skillboxTerms    `json:"terms"`
	Duration skillboxDuration `json:"duration"`
}

type skillboxTerms struct {
	MonthlyPayment int `json:"monthly_payment"`
}

type skillboxDuration struct {
	Months int `json:"count"`
}

func (skillbox Skillbox) Search(query string, filter search.Filter) ([]search.Course, error) {
	if filter.Language == search.LanguageEnglish || filter.Difficulty == search.DifficultyIntermediate {
		return []search.Course{}, nil
	}

	coursesUrl, professionsUrl, err := skillbox.buildSearchUrls(query, filter)
	if err != nil {
		return nil, err
	}

	courses, err := skillbox.fetchCourses(coursesUrl, filter)
	if err != nil {
		return nil, err
	}

	professions, err := skillbox.fetchCourses(professionsUrl, filter)
	if err != nil {
		return nil, err
	}

	return slices.Concat(courses, professions), nil
}

func (skillbox Skillbox) buildSearchUrls(query string, filter search.Filter) (string, string, error) {
	url, err := url.Parse("https://skillbox.ru/api/v6/ru/sales/skillbox/directions/all/nomenclature/course/search")
	if err != nil {
		return "", "", err
	}

	q := url.Query()

	q.Set("search", query)

	if filter.Difficulty == search.DifficultyBeginner {
		q.Set("level", "for novichkov")
	} else if filter.Difficulty == search.DifficultyAdvanced {
		q.Set("level", "for specialists")
	}

	url.RawQuery = q.Encode()

	courseUrl := url.String()
	professionUrl := strings.Replace(courseUrl, "course", "profession", 1)

	return courseUrl, professionUrl, nil
}

func (_ Skillbox) fetchCourses(url string, filter search.Filter) ([]search.Course, error) {
	httpResp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	response := struct {
		Data []skillboxCourse `json:"data"`
	}{}
	err = json.NewDecoder(httpResp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	var courses []search.Course
	for _, course := range response.Data {
		var price float64
		if course.Terms.MonthlyPayment == 0 {
			price = 0
		} else {
			if filter.Free {
				continue
			}
			price = float64(course.Duration.Months * course.Terms.MonthlyPayment)
		}

		courses = append(courses, search.Course{
			Name:        course.Title,
			Url:         course.Href,
			Description: "Обучение с сайта Skillbox.",
			Price:       price,
			Hours:       course.Duration.Months * 1000,
			Extra:       []search.ExtraParam{search.Certificate},
		})
	}

	return courses, nil
}
