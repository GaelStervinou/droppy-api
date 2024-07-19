package drop_type_apis

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type FilmsAPI struct {
}

func (f *FilmsAPI) Search(search string) []ApiSearchResponse {
	apiKey := os.Getenv("TMDB_API_KEY")

	if apiKey == "" {
		log.Printf("Error: The Movie Database API key not found in environment variable TMDB_API_KEY\n")
		return nil
	}
	url := fmt.Sprintf("https://api.themoviedb.org/3/search/multi?api_key=%s&query=%s&page=1", apiKey, search)

	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error: Error trying to get films from tmdb API:", err)
		return nil
	}
	defer resp.Body.Close()

	var result TMDBResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil
	}

	if len(result.Results) == 0 {
		return nil
	}

	var results []ApiSearchResponse
	for _, item := range result.Results {
		var title string
		if item.Title != "" {
			title = item.Title
		} else if item.Name != "" {
			title = item.Name
		} else if item.OriginalTitle != "" {
			title = item.OriginalTitle
		} else if item.OriginalName != "" {
			title = item.OriginalName
		} else {
			title = "Nom manquant"
		}
		var imagePath string
		if item.BackdropPath != "" {
			imagePath = "https://image.tmdb.org/t/p/w500" + item.BackdropPath
		} else if item.PosterPath != "" {
			imagePath = "https://image.tmdb.org/t/p/w500" + item.PosterPath
		}
		results = append(results, ApiSearchResponse{
			Search:      search,
			Title:       title,
			PicturePath: imagePath,
			Subtitle:    item.Overview,
			Content:     strconv.Itoa(item.Id),
		})
	}
	return results
}

func (f *FilmsAPI) Init() {
}

type TMDBResponse struct {
	Page    int `json:"page"`
	Results []struct {
		Title         string `json:"title"`
		OriginalTitle string `json:"original_title"`
		Name          string `json:"name"`
		OriginalName  string `json:"original_name"`
		Overview      string `json:"overview"`
		BackdropPath  string `json:"backdrop_path"`
		PosterPath    string `json:"poster_path"`
		Id            int    `json:"id"`
	} `json:"results"`
}
