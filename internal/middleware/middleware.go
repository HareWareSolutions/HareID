package middleware

import (
	"HareID/internal/authentication"
	"HareID/internal/responses"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type key uint64

const UserKey key = 0

func Authenticate(request http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := authentication.ValidateToken(r); err != nil {
			responses.Error(w, http.StatusUnauthorized, err)
			return
		}

		userID, err := authentication.GetTokenUserID(r)
		if err != nil {
			responses.Error(w, http.StatusUnauthorized, errors.New("Error creating context"))
			return
		}

		log.Printf("User id token: %s", userID)

		ctx := context.WithValue(r.Context(), UserKey, userID)

		if userID, ok := ctx.Value(UserKey).(uint64); ok {
			fmt.Println("This is the user id: ", userID)
		} else {
			fmt.Println("this is a protected route - no user id found")
		}

		request(w, r.WithContext(ctx))
	}
}
