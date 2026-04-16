package service

import (
	"context"
	"fmt"
	"sort"
	"weather-api/internal/client"
)

type WeatherProvider interface {
	GetCurrentWeather(ctx context.Context, lat, lon float64) (*client.ProviderWeatherResponse, error)
	FetchCityCoordinates(ctx context.Context, city string) (*client.GeocodingResult, error)
}

type WeatherResult struct {
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Temperature float64 `json:"temperature"`
	WindSpeed   float64 `json:"wind_speed"`
	WeatherCode int     `json:"weather_code"`
	Time        string  `json:"time"`
	Description string  `json:"description"`
	Recommendation string `json:"recommendation"`
}

type WeatherService struct {
	provider WeatherProvider
}

func NewWeatherService(provider WeatherProvider) *WeatherService {
	return &WeatherService{
		provider: provider,
	}
}

var countryCities = map[string][]string{
	"Kazakhstan": {"Almaty", "Astana", "Shymkent"},
	"Russia":     {"Moscow", "Saint Petersburg", "Novosibirsk"},
}

func (s *WeatherService) GetWeatherByCity(ctx context.Context, city string) (*WeatherResult, error) {
	//fethhcing city
	geo, err := s.provider.FetchCityCoordinates(ctx, city)
	if err != nil {
		return nil, fmt.Errorf("fetch city coordinates: %w", err)
	}

	// getting the weather by coordinates
	res, err := s.GetWeather(ctx, geo.Latitude, geo.Longitude)
	if err != nil {
		return nil, fmt.Errorf("get weather for city %s: %w", city, err)
	}

	return res, nil
}

func (s *WeatherService) GetCountryWeather(ctx context.Context, country string) ([]*WeatherResult, error) {
	cities, ok := countryCities[country]
	if !ok {
		return nil, fmt.Errorf("country %s not supported", country)
	}

	var results []*WeatherResult
	for _, city := range cities {
		res, err := s.GetWeatherByCity(ctx, city)
		if err != nil {
			return nil, fmt.Errorf("failed to get weather for city %s: %w", city, err)
		}
		results = append(results, res)
	}

	return results, nil
}

func (s *WeatherService) GetTopWarmestCities(ctx context.Context, country string) ([]*WeatherResult, error) {
	results, err := s.GetCountryWeather(ctx, country)
	if err != nil {
		return nil, err
	}

	// sorting in desc by temp
	sort.Slice(results, func(i, j int) bool {
		return results[i].Temperature > results[j].Temperature
	})

	// top 3
	if len(results) > 3 {
		return results[:3], nil
	}
	return results, nil
}

func (s *WeatherService) GetWeather(ctx context.Context, lat, lon float64) (*WeatherResult, error) {
	resp, err := s.provider.GetCurrentWeather(ctx, lat, lon)
	if err != nil {
		return nil, fmt.Errorf("get weather from provider: %w", err)
	}

	return &WeatherResult{
		Latitude:       lat,
		Longitude:      lon,
		Temperature:    resp.Temperature,
		WindSpeed:      resp.WindSpeed,
		WeatherCode:    resp.WeatherCode,
		Time:           resp.Time,
		Description:    mapWeatherCode(resp.WeatherCode),
		Recommendation: getRecommendation(resp.Temperature),
	}, nil
}

func getRecommendation(temp float64) string {
	if temp < 10 {
		return "холодно — тёплая одежда"
	}
	if temp < 20 {
		return "прохладно — куртка"
	}
	return "тепло — лёгкая одежда"
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
