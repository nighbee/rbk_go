package service

import (
	"context"
	"fmt"
	"sync"
	"weather-api-t3/internal/client"
	"weather-api-t3/internal/models"
)

type WeatherRepo interface {
	SaveWeatherHistory(ctx context.Context, userID int, city string, temp float64, desc string) error
	GetWeatherHistory(ctx context.Context, userID int, city string, limit int) ([]models.WeatherRecord, error)
}

type WeatherClient interface {
	GetCurrentWeather(ctx context.Context, lat, lon float64) (*client.ProviderWeatherResponse, error)
	FetchCityCoordinates(ctx context.Context, city string) (*client.GeocodingResult, error)
}

type WeatherService struct {
	repo     WeatherRepo
	userRepo UserRepo
	client   WeatherClient
}

func NewWeatherService(repo WeatherRepo, userRepo UserRepo, client WeatherClient) *WeatherService {
	return &WeatherService{
		repo:     repo,
		userRepo: userRepo,
		client:   client,
	}
}

func (s *WeatherService) GetUserWeather(ctx context.Context, userID int) ([]models.WeatherRecord, error) {
	// 1. Check user exists and not deleted
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found or deleted")
	}

	// 2. Get cities
	cities, err := s.userRepo.ListCities(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 3. Fetch weather in parallel
	var wg sync.WaitGroup
	results := make([]models.WeatherRecord, len(cities))
	errs := make([]error, len(cities))

	for i, city := range cities {
		wg.Add(1)
		go func(i int, city string) {
			defer wg.Done()
			
			// Fetch coords
			geo, err := s.client.FetchCityCoordinates(ctx, city)
			if err != nil {
				errs[i] = fmt.Errorf("fetch coords for %s: %w", city, err)
				return
			}

			// Fetch weather
			w, err := s.client.GetCurrentWeather(ctx, geo.Latitude, geo.Longitude)
			if err != nil {
				errs[i] = fmt.Errorf("fetch weather for %s: %w", city, err)
				return
			}

			desc := mapWeatherCode(w.WeatherCode)
			
			// Save to history
			err = s.repo.SaveWeatherHistory(ctx, userID, city, w.Temperature, desc)
			if err != nil {
				errs[i] = fmt.Errorf("save history for %s: %w", city, err)
				return
			}

			results[i] = models.WeatherRecord{
				City:        city,
				Temperature: w.Temperature,
				Description: desc,
			}
		}(i, city)
	}

	wg.Wait()

	// Check if any critical errors occurred (optional: we could return partial results)
	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}

	return results, nil
}

func (s *WeatherService) GetHistory(ctx context.Context, userID int, city string, limit int) (*models.WeatherHistoryResponse, error) {
	history, err := s.repo.GetWeatherHistory(ctx, userID, city, limit)
	if err != nil {
		return nil, err
	}

	return &models.WeatherHistoryResponse{
		UserID:  userID,
		City:    city,
		History: history,
	}, nil
}

func mapWeatherCode(code int) string {
	switch code {
	case 0:
		return "Ясно"
	case 1, 2, 3:
		return "Переменная облачность"
	case 45, 48:
		return "Туман"
	case 51, 53, 55:
		return "Морось"
	case 61, 63, 65:
		return "Дождь"
	case 71, 73, 75:
		return "Снег"
	case 95:
		return "Гроза"
	default:
		return "Неизвестно"
	}
}
