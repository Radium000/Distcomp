package tweet

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	tweetErr "github.com/Khmelov/Distcomp/251003/truhan/lab1/internal/storage/tweet"
	"net/http"
	"strconv"

	"github.com/Khmelov/Distcomp/251003/truhan/lab1/internal/model"
	"github.com/gorilla/mux"
)

type tweet struct {
	srv Service
}

func New(srv Service) tweet {
	return tweet{
		srv: srv,
	}
}

type Service interface {
	Create(ctx context.Context, tweet model.Tweet) (model.Tweet, error)
	GetList(ctx context.Context) ([]model.Tweet, error)
	Get(ctx context.Context, id int64) (model.Tweet, error)
	Update(ctx context.Context, tweet model.Tweet) (model.Tweet, error)
	Delete(ctx context.Context, id int64) error
}

func (i tweet) InitRoutes(r *mux.Router) {
	r.HandleFunc("/tweets", i.getList).Methods(http.MethodGet)
	r.HandleFunc("/tweets", i.create).Methods(http.MethodPost)
	r.HandleFunc("/tweets/{id}", i.get).Methods(http.MethodGet)
	r.HandleFunc("/tweets/{id}", i.deletet).Methods(http.MethodDelete)
	r.HandleFunc("/tweets", i.update).Methods(http.MethodPut)
}

func (i tweet) getList(w http.ResponseWriter, r *http.Request) {
	result, err := i.srv.GetList(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to retrieve result: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (i tweet) create(w http.ResponseWriter, r *http.Request) {
	var req model.Tweet

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("invalid input: %v", err), http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, fmt.Sprintf("validation error: %v", err), http.StatusBadRequest)
		return
	}

	result, err := i.srv.Create(r.Context(), req)
	if err != nil {
		if errors.Is(err, tweetErr.ErrInvalidForeignKey) {
			http.Error(w, fmt.Sprintf("failed to create req: %v", err), http.StatusBadRequest)
		}

		if errors.Is(err, tweetErr.ErrInvalidTweetData) {
			http.Error(w, fmt.Sprintf("failed to create req: %v", err), http.StatusBadRequest)
		}

		http.Error(w, fmt.Sprintf("failed to create req: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

func (i tweet) get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid result ID: %v", err), http.StatusBadRequest)
		return
	}

	result, err := i.srv.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, tweetErr.ErrTweetNotFound) {
			http.Error(w, "result not found", http.StatusNotFound)
			return
		}

		http.Error(w, fmt.Sprintf("failed to retrieve result: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (i tweet) deletet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid tweet ID: %v", err), http.StatusBadRequest)
		return
	}

	if err := i.srv.Delete(r.Context(), id); err != nil {
		if errors.Is(err, tweetErr.ErrTweetNotFound) {
			http.Error(w, "{}", http.StatusNotFound)
			return
		}

		http.Error(w, "{}", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(`{}`))
}

func (i tweet) update(w http.ResponseWriter, r *http.Request) {
	var req model.Tweet

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "{}", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, "{}", http.StatusBadRequest)
		return
	}

	result, err := i.srv.Update(r.Context(), req)
	if err != nil {
		if errors.Is(err, tweetErr.ErrTweetNotFound) {
			http.Error(w, "{}", http.StatusNotFound)
			return
		}

		http.Error(w, "{}", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
