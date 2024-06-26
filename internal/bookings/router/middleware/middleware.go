package appmiddleware

import (
	"errors"
	"net/http"

	"github.com/Elfsilon/car_booking/internal/bookings/router/constants"
)

var ErrMissingUserIDHeader = errors.New("missing 'User-ID' header")

func DefaultHeadersSetter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func UserIdHeaderChecker(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get(constants.UserIdHeaderKey)
		if userID == "" {
			http.Error(w, ErrMissingUserIDHeader.Error(), 400)
			return
		}
		next.ServeHTTP(w, r)
	})
}
