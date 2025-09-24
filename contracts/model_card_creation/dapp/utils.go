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
	Status        bool                     `json:"status"`
	Message       string                   `json:"message"`
	Result        interface{}              `json:"result"`
	SCTDataReply  []SmartContractDataReply `json:"SCTDataReply"`
}

const MODEL_CARD_CONTRACT_HASH = "Qm"

func FetchModelCardsFromChain() ([]Asset, error) {
	reqBody := SmartContractDataRequest{
		Token:  MODEL_CARD_CONTRACT_HASH,
		Latest: false,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post("http://localhost:20000/api/get-smart-contract-token-chain-data",
		"application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result SmartContractDataResponse
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return nil, err
	}

	if len(result.SCTDataReply) == 0 {
		return nil, errors.New("no smart contract data found")
	}

	var cards []Asset
	for _, reply := range result.SCTDataReply {
		if len(reply.SmartContractData) == 0 || reply.SmartContractData[0] != '{' {
			fmt.Printf("Skipping invalid SmartContractData: %s\n", reply.SmartContractData)
			continue
		}

		var wrapper struct {
			CreateModelCard Asset `json:"create_model_card"`
		}

		err := json.Unmarshal([]byte(reply.SmartContractData), &wrapper)
		if err != nil {
			fmt.Printf("Failed to parse SmartContractData: %s\n", reply.SmartContractData)
			continue
		}

		cards = append(cards, wrapper.CreateModelCard)
	}

	return cards, nil
}
