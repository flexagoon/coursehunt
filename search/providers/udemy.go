package providers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"fxgn.dev/coursehunt/search"
)

type Udemy struct {
	ClientId     string
	ClientSecret string
}

type udemyCourse struct {
	Title       string            `json:"title"`
	UrlPart     string            `json:"url"`
	Headline    string            `json:"headline"`
	Price       string            `json:"price"`
	Instructors []udemyInstructor `json:"visible_instructors"`
	Duration    string            `json:"content_info_short"`
	Rating      float64           `json:"avg_rating"`
}

type udemyInstructor struct {
	Name string `json:"display_name"`
}

func (udemy Udemy) Search(query string, filter search.Filter) ([]search.Course, error) {
	url, err := udemy.buildSearchUrl(query, filter)
	if err != nil {
		return nil, err
	}

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

	var courses []search.Course
	for _, course := range response.Results {
		duration, _ := strconv.ParseFloat(strings.Split(course.Duration, " ")[0], 64)

		courses = append(courses, search.Course{
			Name:        course.Title,
			Url:         "https://www.udemy.com" + course.UrlPart,
			Description: course.Headline,
			Price:       course.Price,
			Author:      course.Instructors[0].Name,
			Hours:       int(duration),
			Rating:      course.Rating,
			Extra:       []search.ExtraParam{search.Certificate},
		})
	}

	return courses, nil
}

func (_ Udemy) buildSearchUrl(query string, filter search.Filter) (string, error) {
	url, err := url.Parse("https://www.udemy.com/api-2.0/courses?fields[course]=title,url,headline,price,visible_instructors,avg_rating,content_info_short&fields[user]=display_name")
	if err != nil {
		return "", err
	}

	q := url.Query()

	q.Set("search", query)

	if filter.Free {
		q.Set("price", "price-free")
	}

	if filter.Language == search.LanguageEnglish {
		q.Set("language", "en")
	} else if filter.Language == search.LanguageRussian {
		q.Set("language", "ru")
	}

	if filter.Difficulty == search.DifficultyBeginner {
		q.Set("instructional_level", "beginner")
	} else if filter.Difficulty == search.DifficultyIntermediate {
		q.Set("instructional_level", "intermediate")
	} else if filter.Difficulty == search.DifficultyAdvanced {
		q.Set("instructional_level", "expert")
	}

	url.RawQuery = q.Encode()

	return url.String(), nil
}
