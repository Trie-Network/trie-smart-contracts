package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type SmartContractDataReply struct {
	BlockNo            int    `json:"BlockNo"`
	BlockId            string `json:"BlockId"`
	SmartContractData  string `json:"SmartContractData"`
	Epoch              int64  `json:"Epoch"`
	InitiatorSignature string `json:"InitiatorSignature"`
	ExecutorDID        string `json:"ExecutorDID"`
	InitiatorSignData  string `json:"InitiatorSignData"`
}

type SmartContractDataRequest struct {
	Token  string `json:"token,omitempty"`
	Latest bool   `json:"latest"`
}

type SmartContractDataResponse struct {
	Status        bool                      `json:"status"`
	Message       string                    `json:"message"`
	Result        interface{}               `json:"result"`
	SCTDataReply  []SmartContractDataReply  `json:"SCTDataReply"`
}

type Rating struct {
    UserDID  string `json:"user_did"`
    Rating   int    `json:"rating"`
    AssetID  string `json:"asset_id"`
}

type WrappedRating struct {
    RateAsset Rating `json:"rate_asset"`
}

const RATING_CONTRACT_HASH = "QmeywD5LVUpZ1TiirMtAR8YyAb3suLdEUETJxkKwhUWHW8"

func GetRatingFromChain(assetID string) (float64, error) {
	reqBody := SmartContractDataRequest{
		Token:  RATING_CONTRACT_HASH,
		Latest: false, 
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return 0, err
	}

	fmt.Printf("Sending request body to Rubix: %s\n", string(bodyBytes))

	resp, err := http.Post("http://localhost:20000/api/get-smart-contract-token-chain-data", "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	fmt.Printf("Raw API response: %s\n", string(respBody))

	var result SmartContractDataResponse
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return 0, err
	}

	if len(result.SCTDataReply) == 0 {
		return 0, fmt.Errorf("Smart Contract Token %v is not registered", RATING_CONTRACT_HASH)
	}

	type userRating struct {
		Rating int
		Epoch  int64
	}

	latest := make(map[string]userRating)

	for _, reply := range result.SCTDataReply {
		if len(reply.SmartContractData) == 0 || reply.SmartContractData[0] != '{' {
			fmt.Printf("Skipping invalid SmartContractData: %s\n", reply.SmartContractData)
			continue
		}

		var wrapper WrappedRating
		err := json.Unmarshal([]byte(reply.SmartContractData), &wrapper)
		if err != nil {
			fmt.Printf("Failed to parse  SmartContractData: %s\n", reply.SmartContractData)
			continue
		}

		entry := wrapper.RateAsset
		if entry.AssetID == assetID && entry.Rating >= 1 && entry.Rating <= 5 {
			prev, exists := latest[entry.UserDID]
			if !exists || reply.Epoch > prev.Epoch {
				latest[entry.UserDID] = userRating{
					Rating: entry.Rating,
					Epoch:  reply.Epoch,
				}
			}
		}
	}

	if len(latest) == 0 {
		return 0, errors.New("no valid ratings found for asset")
	}

	fmt.Println("Latest ratings per DID:")
	total := 0
	for did, ur := range latest {
		fmt.Printf("  %s -> %d (Epoch %d)\n", did, ur.Rating, ur.Epoch)
		total += ur.Rating
	}

	average := float64(total) / float64(len(latest))
	return average, nil
}