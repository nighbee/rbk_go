package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"weather-api/internal/service"

	"github.com/go-chi/chi/v5"
)

type Service interface {
	GetWeather(ctx context.Context, lat, lon float64) (*service.WeatherResult, error)
	GetWeatherByCity(ctx context.Context, city string) (*service.WeatherResult, error)
	GetCountryWeather(ctx context.Context, country string) ([]*service.WeatherResult, error)
	GetTopWarmestCities(ctx context.Context, country string) ([]*service.WeatherResult, error)
}

type WeatherHandler struct {
	service Service
}

func NewWeatherHandler(service Service) *WeatherHandler {
	return &WeatherHandler{
		service: service,
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *WeatherHandler) GetWeather(w http.ResponseWriter, r *http.Request) {
	latStr := r.URL.Query().Get("lat")
	lonStr := r.URL.Query().Get("lon")

	if latStr == "" || lonStr == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "query params lat and lon are required",
		})
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "invalid lat",
		})
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "invalid lon",
		})
		return
	}

	result, err := h.service.GetWeather(r.Context(), lat, lon)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, `{"error":"failed to encode json"}`, http.StatusInternalServerError)
	}
}

func (h *WeatherHandler) GetWeatherByCity(w http.ResponseWriter, r *http.Request) {
	city := chi.URLParam(r, "city")
	if city == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "city is required"})
		return
	}
	result, err := h.service.GetWeatherByCity(r.Context(), city)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *WeatherHandler) GetCountryWeather(w http.ResponseWriter, r *http.Request) {
	country := chi.URLParam(r, "country")
	if country == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "country is required"})
		return
	}
	results, err := h.service.GetCountryWeather(r.Context(), country)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, results)
}

func (h *WeatherHandler) GetTopCities(w http.ResponseWriter, r *http.Request) {
	country := chi.URLParam(r, "country")
	if country == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "country is required"})
		return
	}
	results, err := h.service.GetTopWarmestCities(r.Context(), country)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, results)
}
