package main

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/luthfiaed/owapp-be/internal/data"
	"golang.org/x/crypto/bcrypt"
)

func (app *application) getUserByUsername(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	username := user.Username
	if username == "" {
		app.notFoundResponse(w, r)
		return
	}

	user, err := app.users.GetByUsername(username)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, user, nil)
	if err != nil {
		app.logger.Println(err)
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) registerNewUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Username == "" || input.Password == "" || input.Role == "" {
		app.badRequestResponse(w, r, errors.New("input can not be empty"))
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	user := &data.User{
		Username: input.Username,
		Password: string(hash),
		Role:     input.Role,
	}

	err = app.users.Insert(user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, user, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) uploadAvatar(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	username := user.Username
	if username == "" {
		app.notFoundResponse(w, r)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10MB limit

	/*
		VULNERABILITY POINT:
		1. File name not sanitized
		2. File type not checked
		Can potentially rewrite existing files
	*/
	file, header, err := r.FormFile("file")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer file.Close()

	uploadDir := "./public"
	if err = os.MkdirAll(uploadDir, 0o755); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	dest, err := os.Create(filepath.Join(uploadDir, header.Filename))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer dest.Close()

	_, err = io.Copy(dest, file)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.users.UpdateAvatar(username, header.Filename)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, map[string]string{"message": "file uploaded successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getAvatar(w http.ResponseWriter, r *http.Request) {
	filename := r.PathValue("filename")
	http.ServeFile(w, r, filepath.Join("./public", filename))
}
