package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	helper "github.com/danielzinhors/go_telemetry/service_weather/internal/helpers"
	"github.com/danielzinhors/go_telemetry/service_weather/internal/usecases"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type GetTempResponse struct {
	TemperatureCelsius    float64 `json:"temp_C"`
	TemperatureFahrenheit float64 `json:"temp_F"`
	TemperatureKelvin     float64 `json:"temp_K"`
}

type GetTempHandler struct {
	getTempUseCase usecases.GetTempUseCase
}

func NewGetTempHandler(getTempUseCase usecases.GetTempUseCase) *GetTempHandler {
	return &GetTempHandler{
		getTempUseCase: getTempUseCase,
	}
}

func (h *GetTempHandler) Handle(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	tracer := otel.Tracer("")
	ctx, span := tracer.Start(ctx, "GetTempHandler")
	defer span.End()
	cep, ok := h.getCepFromRequest(r)
	if !ok {
		w.WriteHeader(422)
		w.Write([]byte("invalid zipcode"))
		return
	}

	input := &usecases.TempInput{Cep: cep}
	output, err := h.getTempUseCase.Execute(ctx, input)

	if err != nil {
		if err.Error() == "can not find zipcode" {
			w.WriteHeader(404)
			w.Write([]byte("can not find zipcode"))
		} else {
			fmt.Printf("ERROR: %s\n", err.Error())
			w.WriteHeader(500)
			w.Write([]byte("internal server error"))
		}
		return
	}

	response := GetTempResponse{
		TemperatureCelsius:    helper.RoundFloat(output.TemperatureCelsius, 1),
		TemperatureFahrenheit: helper.RoundFloat(output.TemperatureFahrenheit, 1),
		TemperatureKelvin:     helper.RoundFloat(output.TemperatureKelvin, 1),
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *GetTempHandler) getCepFromRequest(r *http.Request) (string, bool) {
	cep := r.URL.Query().Get("cep")

	if cep == "" {
		return "", false
	}

	cepRegex := regexp.MustCompile(`^\d{5}-{0,1}\d{3}$`)
	if !cepRegex.Match([]byte(cep)) {
		return "", false
	}

	return cep, true
}
