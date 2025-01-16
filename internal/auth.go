package main

import (
	"net/http"
	"time"

	"github.com/pascaldekloe/jwt"
	"golang.org/x/crypto/bcrypt"
)

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.users.GetByUsername(input.Username)
	if err != nil {
		// TODO ceck behavior ini pastiin explicit error message
		app.serverErrorResponse(w, r, err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var claims jwt.Claims
	claims.Subject = user.Username
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))
	claims.Issuer = "owapp-be"
	claims.Audiences = []string{"owapp"}

	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.config.JwtSecret))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	res := map[string]string{
		"access_token": string(jwtBytes),
	}
	err = app.writeJSON(w, http.StatusOK, res, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
