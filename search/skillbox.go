package search

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
)

type Skillbox struct{}

type skillboxCourse struct {
	Title string `json:"title"`
	Href  string `json:"href"`
	Terms terms  `json:"terms"`
}

type terms struct {
	MonthlyPayment int    `json:"monthly_payment"`
	Currency       string `json:"icon_currency"`
}

func (skillbox Skillbox) Search(query string, filter Filter) ([]Course, error) {
	if filter.Language == LanguageEnglish || filter.Difficulty == DifficultyIntermediate {
		return []Course{}, nil
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

func (skillbox Skillbox) buildSearchUrls(query string, filter Filter) (string, string, error) {
	url, err := url.Parse("https://skillbox.ru/api/v6/ru/sales/skillbox/directions/all/nomenclature/course/search")
	if err != nil {
		return "", "", err
	}

	q := url.Query()

	q.Set("search", query)

	if filter.Difficulty == DifficultyBeginner {
		q.Set("level", "for novichkov")
	} else if filter.Difficulty == DifficultyAdvanced {
		q.Set("level", "for specialists")
	}

	url.RawQuery = q.Encode()

	courseUrl := url.String()
	professionUrl := strings.Replace(courseUrl, "course", "profession", 1)

	return courseUrl, professionUrl, nil
}

func (_ Skillbox) fetchCourses(url string, filter Filter) ([]Course, error) {
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

	var courses []Course
	for _, course := range response.Data {
		var price string
		if course.Terms.MonthlyPayment == 0 {
			price = "Free"
		} else {
			if filter.Free {
				continue
			}
			price = fmt.Sprint(course.Terms.MonthlyPayment, course.Terms.Currency)
		}
		courses = append(courses, Course{
			Name:        course.Title,
			Url:         course.Href,
			Description: "",
			Price:       price,
		})
	}

	return courses, nil
}
