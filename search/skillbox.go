package search

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
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

const skillboxCourseBaseUrl = "https://skillbox.ru/api/v6/ru/sales/skillbox/directions/all/nomenclature/course/search/?search="
const skillboxProfessionBaseUrl = "https://skillbox.ru/api/v6/ru/sales/skillbox/directions/all/nomenclature/profession/search/?search="

func (skillbox Skillbox) Search(query string, filter Filter) ([]Course, error) {
	coursesUrl := skillboxCourseBaseUrl + query
	professionsUrl := skillboxProfessionBaseUrl + query

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

func (_ Skillbox) fetchCourses(url string, filter Filter) ([]Course, error) {
	fmt.Println(url)
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
