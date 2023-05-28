package handlers

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Repository interface {
	GetByUid(uid string) ([]byte, error)
}

type Handler struct {
	repo Repository
}

func New(repository Repository) *Handler {
	return &Handler{
		repo: repository,
	}
}

func (h *Handler) InitRoute() http.Handler {
	rtr := httprouter.New()
	rtr.GET("/order/:uid", h.ShowById)

	return rtr
}

func (h *Handler) ShowById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uid := ps.ByName("uid")

	if uid == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("uid is empty"))
		return
	}

	order, err := h.repo.GetByUid(uid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("order with uid [%s]: %s", uid, err.Error())))
		return
	}

	w.Write(order)
}
