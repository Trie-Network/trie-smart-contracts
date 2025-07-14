package main

type SmartContractResponse struct {
	BasicResponse
	SCTDataReply []SCTDataReply
}

type BasicResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

type SCTDataReply struct {
	BlockNo            uint64
	BlockId            string
	SmartContractData  string
	Epoch              int
	InitiatorSignature string
	ExecutorDID        string
	InitiatorSignData  string
}

type ContractInputRequest struct {
	Port              string `json:"port"`
	SmartContractHash string `json:"smart_contract_hash"`
	SmartContractData string `json:"smart_contract_data"`
	InitiatorDID      string `json:"initiator_did"`
}

type RatingResponse struct {
	Data   RatingInfo `json:"data"`
}

type RatingInfo struct {
	AssetID       string  `json:"asset_id"`
	AverageRating float64 `json:"average_rating"`
}
