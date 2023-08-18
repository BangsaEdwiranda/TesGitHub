package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go-module/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		utils.SuccessResponse(w, "Have fun!", http.StatusOK)
	})

	r.Post("/ads", CreateAd)

	http.ListenAndServe(":3333", r)
}

func CreateAd(w http.ResponseWriter, r *http.Request) {
	// Decode and parse the payload.
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Unable to read the payload: %s\n", err.Error())

		utils.ErrorResponse(w, "Unable to read the payload.", http.StatusBadRequest)
		return
	}

	var ad Ad
	if err := json.Unmarshal(body, &ad); err != nil {
		fmt.Printf("Unable to decode the payload: %s\n", err.Error())

		utils.ErrorResponse(w, "Unable to decode the payload.", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received request to create a new ad with payload: %+v\n", ad)

	// Validation.
	if ad.AdvertiserID < 1 {
		fmt.Printf("AdvertiserID must be greater than 0, received: %d\n", ad.AdvertiserID)

		utils.ErrorResponse(w, "advertiser_id must be greater than 0", http.StatusBadRequest)
		return
	}

	if len(ad.Audiences) == 0 {
		fmt.Println("Audiences cannot be empty")

		utils.ErrorResponse(w, "audiences cannot be empty", http.StatusBadRequest)
		return
	}

	// Make request to an "internal API".
	requestBody, err := json.Marshal(ad)
	if err != nil {
		fmt.Println(err)

		utils.ErrorResponse(w, "Unable to complete the request.", http.StatusInternalServerError)
		return
	}

	resp, err := http.Post("https://httpbin.org/post", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Printf("Request to API failed with error %v\n", err.Error())

		utils.ErrorResponse(w, "Unable to complete the request.", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Internal API returned status code %d and error %s\n", resp.StatusCode, resp.Status)
		utils.ErrorResponse(w, "Unable to complete the request.", http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, "Successful request.", http.StatusOK)
}

type Ad struct {
	AdvertiserID uint64   `json:"advertiser_id"`
	Audiences    []uint64 `json:"audiences"`
	Bid          float32  `json:"bid"`
	ImageUrl     string   `json:"image_url"`
}
