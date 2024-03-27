package search

import (
	"encoding/json"
	"net/http"
	"net/url"
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

func (udemy Udemy) Search(query string, filter Filter) ([]Course, error) {
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

func (_ Udemy) buildSearchUrl(query string, filter Filter) (string, error) {
	url, err := url.Parse("https://www.udemy.com/api-2.0/courses")
	if err != nil {
		return "", err
	}

	q := url.Query()

	q.Set("search", query)

	if filter.Free {
		q.Set("price", "price-free")
	}

	if filter.Language == LanguageEnglish {
		q.Set("language", "en")
	} else if filter.Language == LanguageRussian {
		q.Set("language", "ru")
	}

	if filter.Difficulty == DifficultyBeginner {
		q.Set("instructional_level", "beginner")
	} else if filter.Difficulty == DifficultyIntermediate {
		q.Set("instructional_level", "intermediate")
	} else if filter.Difficulty == DifficultyAdvanced {
		q.Set("instructional_level", "expert")
	}

	url.RawQuery = q.Encode()

	return url.String(), nil
}
