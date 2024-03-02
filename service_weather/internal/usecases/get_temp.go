package usecases

import (
	"context"
	"fmt"

	helper "github.com/danielzinhors/go_telemetry/service_weather/internal/helpers"
	"github.com/danielzinhors/go_telemetry/service_weather/internal/services"
)

type TempInput struct {
	Cep string
}

type TempOutput struct {
	TemperatureCelsius    float64
	TemperatureFahrenheit float64
	TemperatureKelvin     float64
}

type GetTempUseCase interface {
	Execute(ctx context.Context, input *TempInput) (*TempOutput, error)
}

type GetTempUseCaseImpl struct {
	cepService     services.ViaCepService
	weatherService services.WeatherApiService
}

func NewGetTempUseCase(cepService services.ViaCepService, weatherService services.WeatherApiService) GetTempUseCase {
	return &GetTempUseCaseImpl{
		cepService:     cepService,
		weatherService: weatherService,
	}
}

func (u *GetTempUseCaseImpl) Execute(ctx context.Context, input *TempInput) (*TempOutput, error) {

	cepResponse, err := u.cepService.QueryCep(ctx, input.Cep)
	if err != nil {
		return nil, err
	}

	location := fmt.Sprintf("Brazil - %s - %s", helper.StateMap[cepResponse.UF], cepResponse.Localidade)

	weatherResponse, err := u.weatherService.QueryWeather(ctx, location)
	if err != nil {
		return nil, err
	}

	temperatureCelsius := weatherResponse.Current.TemperatureCelsius
	temperatureFahrenheit := weatherResponse.Current.TemperatureCelsius*1.8 + 32
	temperatureKelvin := weatherResponse.Current.TemperatureCelsius + 273.15

	weatherLocation := fmt.Sprintf("%s - %s - %s", weatherResponse.Location.Name, weatherResponse.Location.Region, weatherResponse.Location.Country)
	fmt.Printf("Temperatura em %s C=%.2f F=%.2f K=%.2f\n", weatherLocation, temperatureCelsius, temperatureFahrenheit, temperatureKelvin)

	return &TempOutput{
		TemperatureCelsius:    temperatureCelsius,
		TemperatureFahrenheit: temperatureFahrenheit,
		TemperatureKelvin:     temperatureKelvin,
	}, nil
}
