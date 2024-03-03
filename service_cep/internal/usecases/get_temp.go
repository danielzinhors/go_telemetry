package usecases

import (
	"context"

	"github.com/danielzinhors/go_telemetry/service_cep/internal/services"
)

type TempInput struct {
	Cep string
}

type TempOutput struct {
	StatusCode int
	Data       string
}

type GetTempUseCase interface {
	Execute(ctx context.Context, input *TempInput) (*TempOutput, error)
}

type GetTempUseCaseImpl struct {
	serviceWeatherService services.ServiceWeatherService
}

func NewGetTempUseCase(ServiceWeatherService services.ServiceWeatherService) GetTempUseCase {
	return &GetTempUseCaseImpl{
		serviceWeatherService: ServiceWeatherService,
	}
}

func (u *GetTempUseCaseImpl) Execute(ctx context.Context, input *TempInput) (*TempOutput, error) {
	serviceWeatherResponse, err := u.serviceWeatherService.QueryCep(ctx, input.Cep)

	if err != nil {
		return nil, err
	}

	return &TempOutput{
		StatusCode: serviceWeatherResponse.StatusCode,
		Data:       serviceWeatherResponse.Data,
	}, nil
}
