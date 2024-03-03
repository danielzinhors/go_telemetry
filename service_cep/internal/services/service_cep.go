package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type ServiceWeatherResponse struct {
	StatusCode int
	Data       string
}

type ServiceWeatherService interface {
	QueryCep(ctx context.Context, cep string) (*ServiceWeatherResponse, error)
}

type ServiceWeatherServiceImpl struct {
	client *http.Client
}

func NewServiceWeatherService() ServiceWeatherService {
	return &ServiceWeatherServiceImpl{
		client: &http.Client{},
	}
}

func (s *ServiceWeatherServiceImpl) QueryCep(ctx context.Context, cep string) (*ServiceWeatherResponse, error) {
	url := fmt.Sprintf("%s?cep=%s", os.Getenv("SERVICE_WEATHER_URL"), cep)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request = request.WithContext(ctx)

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(request.Header))
	response, err := s.client.Do(request)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		if response.StatusCode == 404 {
			return &ServiceWeatherResponse{
				StatusCode: response.StatusCode,
				Data:       string(body),
			}, nil
		}
		return nil, errors.New("service weather error")
	}

	return &ServiceWeatherResponse{
		StatusCode: response.StatusCode,
		Data:       string(body),
	}, nil
}
