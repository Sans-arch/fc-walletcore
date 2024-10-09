package web

import (
	"encoding/json"
	"net/http"

	"github.com/Sans-arch/fc-walletcore/internal/usecase/create_client"
)

type WebClientHandler struct {
	CreateClientUsecase create_client.CreateClientUseCase
}

func NewWebClientHandler(createClientUsecase create_client.CreateClientUseCase) *WebClientHandler {
	return &WebClientHandler{
		CreateClientUsecase: createClientUsecase,
	}
}

func (h *WebClientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	var dto create_client.CreateClientInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	output, err := h.CreateClientUsecase.Execute(dto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
