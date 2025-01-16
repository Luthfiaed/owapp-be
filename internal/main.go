package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/luthfiaed/owapp-be/internal/data"
)

type config struct {
	Env            string
	Port           int
	Dsn            string
	JwtSecret      string
	AllowedOrigins []string
}

type application struct {
	config   config
	logger   *log.Logger
	users    *data.UserModel
	products *data.ProductModel
}

const version = "1.0.0"

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	err := godotenv.Load("./env/.env")
	if err != nil {
		logger.Fatal(err)
	}

	cfg, err := loadCfgFromEnv()
	if err != nil {
		logger.Fatal(err)
	}

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	app := &application{
		config:   cfg,
		logger:   logger,
		users:    &data.UserModel{DB: db},
		products: &data.ProductModel{DB: db},
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      app.InitHandler(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("starting %s server on %s", cfg.Env, srv.Addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)
}

func loadCfgFromEnv() (config, error) {
	var cfg config
	cfg.Dsn = os.Getenv("DB_DSN")
	cfg.JwtSecret = os.Getenv("JWT_SECRET")
	cfg.Env = os.Getenv("ENVIRONMENT")
	cfg.AllowedOrigins = strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")

	fmt.Println("allowed origins: ", cfg.AllowedOrigins)

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return config{}, err
	}
	cfg.Port = port

	// zero value is not allowed
	v := reflect.ValueOf(cfg)
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()
		if reflect.ValueOf(value).IsZero() {
			err := fmt.Errorf("field: %s, value: %v have zero value", field.Name, value)
			return config{}, err
		}
	}

	return cfg, nil
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.Dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(15 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
