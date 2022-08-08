package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"
)

type NationalIDRequestDTO struct {
	IdNumber                 string `json:"IdNumber"`
	Surname                  string `json:"Surname"`
	FirstName                string `json:"FirstName"`
	OtherNames               string `json:"OtherNames"`
	Nationality              string `json:"Nationality"`
	Gender                   string `json:"Gender"`
	DateOfBirthString        string `json:"DateOfBirthString"`
	DateOfIssueString        string `json:"DateOfIssueString"`
	DateOfExpiryString       string `json:"DateOfExpiryString"`
	PlaceOfBirthDistrictName string `json:"PlaceOfBirthDistrictName"`
	Status                   string `json:"Status"`
}

type NationalIDResponseDTO struct {
	DateOfBirthString        string `json:"DateOfBirthString"`
	DateOfExpiryString       string `json:"DateOfExpiryString"`
	DateOfIssueString        string `json:"DateOfIssueString"`
	FirstName                string `json:"FirstName"`
	Gender                   string `json:"Gender"`
	IdNumber                 string `json:"IdNumber"`
	Nationality              string `json:"Nationality"`
	OtherNames               string `json:"OtherNames"`
	PlaceOfBirthDistrictName string `json:"PlaceOfBirthDistrictName"`
	Status                   string `json:"Status"`
	Surname                  string `json:"Surname"`
}

func writeHttpError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	w.Write([]byte(message))
	return
}

func main() {
	//parse flags
	baseUrlPtr := flag.String("baseurl", "http://localhost:8080", "Base URL of the service")
	clientIdPtr := flag.String("clientid", "SomeClientID", "Client ID of the service")
	clientSecretPtr := flag.String("clientsecret", "SomeClientSecret", "Client Secret of the service")
	portPtr := flag.Int("port", 8080, "Port to listen on")

	flag.Parse()

	baseURL := *baseUrlPtr
	clientID := *clientIdPtr
	clientSecret := *clientSecretPtr
	port := *portPtr

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to NRB Proxy"))
	})
	//http endpoint
	http.HandleFunc("/api/verify", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			//read request body
			var request NationalIDRequestDTO
			err := json.NewDecoder(r.Body).Decode(&request)

			if err != nil {
				writeHttpError(w, http.StatusBadRequest, "Invalid request body")
			}
			//generate http request
			requestBody, _ := json.Marshal(request)

			req, err := http.NewRequest("POST", baseURL+"/api/verify/postverify", bytes.NewBuffer(requestBody))

			if err != nil {
				writeHttpError(w, http.StatusInternalServerError, "Error verifying with external service")
			}

			//set headers
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("ClientID", clientID)
			req.Header.Set("ClientSecret", clientSecret)

			//make http request
			response, err := http.DefaultClient.Do(req)

			if err != nil {
				writeHttpError(w, http.StatusInternalServerError, "Error verifying with external service")
			}

			//unmarshall response
			var responseDTO NationalIDResponseDTO

			err = json.NewDecoder(response.Body).Decode(&responseDTO)

			if err != nil {
				writeHttpError(w, http.StatusInternalServerError, "Error verifying with external service")
			}

			//write response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(responseDTO)
		}
	})

	portStr := ":" + strconv.Itoa(port)
	log.Println("Listening on port :" + portStr + "...")

	err := http.ListenAndServe(portStr, nil)
	if err != nil {
		panic(err)
		return
	}
}
