package middleware

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type ContextKey string

const IdKey ContextKey = "id"

func IdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr, ok := vars["id"]
		if !ok {
			log.Info().Msg("Missed ID")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Warn().Err(err).Msg("Invalid ID param")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), IdKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
