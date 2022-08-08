package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
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

func validateResponse(responseDTO NationalIDResponseDTO) (bool, string) {
	const Correct = "CORRECT"
	const Valid = "VALID"
	//print response object
	log.Println(responseDTO)

	if responseDTO.DateOfBirthString != Correct {
		return false, "Date of birth is not correct"
	}
	if responseDTO.DateOfExpiryString != Correct {
		return false, "Date of expiry is not correct"
	}
	if responseDTO.DateOfIssueString != Correct {
		return false, "Date of issue is not correct"
	}
	if responseDTO.FirstName != Correct {
		return false, "First name is not correct"
	}
	if responseDTO.Gender != Correct {
		return false, "Gender is not correct"
	}
	if responseDTO.IdNumber != Correct {
		return false, "Id number is not correct"
	}
	if responseDTO.Nationality != Correct {
		return false, "Nationality is not correct"
	}
	if responseDTO.OtherNames != Correct {
		return false, "Other Names is not correct"
	}
	if responseDTO.Status != Valid {
		return false, "Status is not valid"
	}
	if responseDTO.Surname != Correct {
		return false, "Surname is not correct"
	}

	return true, responseDTO.Status
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
				return
			}
			//generate http request
			requestBody, _ := json.Marshal(request)

			req, err := http.NewRequest("POST", baseURL+"/api/verify/postverify", bytes.NewBuffer(requestBody))

			if err != nil {
				writeHttpError(w, http.StatusInternalServerError, "Error generating request")
				return
			}

			//set headers
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("ClientId", clientID)
			req.Header.Set("ClientKey", clientSecret)

			//make http request
			response, err := http.DefaultClient.Do(req)

			if err != nil {
				writeHttpError(w, http.StatusInternalServerError, "Error verifying with external service")
				return
			}

			if response.StatusCode != http.StatusOK {
				writeHttpError(w, response.StatusCode, "Error verifying with external service")
				return
			}

			//unmarshall response
			var responseDTO NationalIDResponseDTO

			body, _ := ioutil.ReadAll(response.Body)

			err = json.Unmarshal(body, &responseDTO)

			if err != nil {
				log.Println(body)
				writeHttpError(w, http.StatusInternalServerError, "Error Decoding response")
				return
			}

			valid, msg := validateResponse(responseDTO)

			if !valid {
				writeHttpError(w, http.StatusUnprocessableEntity, msg)
				return
			}
			//write response
			w.WriteHeader(http.StatusOK)
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
