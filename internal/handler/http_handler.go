package handler

import (
	"net/http"

	"github.com/FarhanDH/sidingin-rating-service/internal/models"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (*Handler) Health(w http.ResponseWriter, r *http.Request) {
	models.WriteSuccess(w, http.StatusOK, "oke", map[string]string{"ping": "pong"})
}
