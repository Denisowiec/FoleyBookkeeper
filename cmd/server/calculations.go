package main

import (
	"encoding/json"
	"net/http"

	"github.com/Denisowiec/FoleyBookkeeper/internal/db"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func (cfg *apiConfig) handlerCreateCalculation(w http.ResponseWriter, r *http.Request) {
	_, _, err := authenticateUser(r, cfg.secret)
	if err != nil {
		respondWithError(w, "Operation unauthorized", http.StatusUnauthorized, err)
		return
	}

	calcInput := struct {
		ProjectID    string `json:"project_id"`
		Budget       string `json:"budget"`
		Currency     string `json:"currency"`
		ExchangeRate string `json:"exchange_rate"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&calcInput)
	if err != nil {
		respondWithError(w, "Error decoding user input", http.StatusBadRequest, err)
		return
	}

	// We convert the values to our preferred types
	// The string type from the json data is fine for db input,
	// but the conversions validate the input and allow us to make defaults
	projectID, err := uuid.Parse(calcInput.ProjectID)
	if err != nil {
		respondWithError(w, "Error decoding user input", http.StatusBadRequest, err)
		return
	}
	var budget decimal.Decimal
	if calcInput.Budget != "" {
		budget, err = decimal.NewFromString(calcInput.Budget)
		if err != nil {
			respondWithError(w, "Error decoding user input", http.StatusBadRequest, err)
			return
		}
	} else {
		budget = decimal.NewFromInt(0)
	}
	currency := calcInput.Currency
	var exchangeRate decimal.Decimal
	if calcInput.ExchangeRate != "" {
		exchangeRate, err = decimal.NewFromString(calcInput.ExchangeRate)
		if err != nil {
			respondWithError(w, "Error decoding user input", http.StatusBadRequest, err)
			return
		}
	} else {
		exchangeRate = decimal.NewFromInt(1)
	}

	createCalcParams := db.CreateCalculationParams{
		ProjectID:    projectID,
		Budget:       budget.String(),
		Currency:     currency,
		ExchangeRate: exchangeRate.String(),
	}

	calc, err := cfg.db.CreateCalculation(r.Context(), createCalcParams)
	if err != nil {
		respondWithError(w, "Error contacting database", http.StatusInternalServerError, err)
		return
	}

	err = respondWithJSON(w, http.StatusCreated, calc)
	if err != nil {
		respondWithError(w, "Error processing return data", http.StatusInternalServerError, err)
		return
	}
}

/*func (cfg *apiConfig) handlerUpdateCalculation(w http.ResponseWriter, r *http.Request) {
	_, _, err := authenticateUser(r, cfg.secret)
	if err != nil {
		respondWithError(w, "Operation unauthorized", http.StatusUnauthorized, err)
		return
	}

}*/

func (cfg *apiConfig) handlerAddEpisodesToCalculation(w http.ResponseWriter, r *http.Request) {
	_, _, err := authenticateUser(r, cfg.secret)
	if err != nil {
		respondWithError(w, "Operation unauthorized", http.StatusUnauthorized, err)
		return
	}

	addEppsInput := struct {
		EpisodeID string `json:"episode_id"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&addEppsInput)
	if err != nil {
		respondWithError(w, "Error decoding user input", http.StatusBadRequest, err)
		return
	}

	episodeID, err := uuid.Parse(addEppsInput.EpisodeID)
	if err != nil {
		respondWithError(w, "Error decoding user input", http.StatusBadRequest, err)
		return
	}
	calcID, err := uuid.Parse(r.PathValue("calcid"))
	if err != nil {
		respondWithError(w, "Error decoding user input", http.StatusBadRequest, err)
		return
	}

	addEppsParams := db.AddEpisodeToCalculationParams{
		EpisodeID: episodeID,
		CalcID:    calcID,
	}

	ret, err := cfg.db.AddEpisodeToCalculation(r.Context(), addEppsParams)
	if err != nil {
		respondWithError(w, "Error contacting database", http.StatusInternalServerError, err)
		return
	}

	err = respondWithJSON(w, http.StatusAccepted, ret)
	if err != nil {
		respondWithError(w, "Error processing return data", http.StatusInternalServerError, err)
		return
	}
}

/*func (cfg *apiConfig) handlerGetCalculation(w http.ResponseWriter, r *http.Request) {
	_, _, err := authenticateUser(r, cfg.secret)
	if err != nil {
		respondWithError(w, "Operation unauthorized", http.StatusUnauthorized, err)
		return
	}


}*/
