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
type Asset struct {
	ID                string          `json:"id,omitempty"`
	Name              string          `json:"name"`
	MainCategory      string          `json:"main_category"`
	SecondaryCategory string          `json:"secondary_category"`
	Description       string          `json:"description"`
	Metrics           string          `json:"metrics"`
	AssetType         string          `json:"asset_type"`
	AssetID           *string         `json:"asset_id,omitempty"` 
	Owner             string          `json:"owner"`
	
}