# Instructions

## Setup Contracts

Call the following rubix endpoints in steps:

1. Generate Contract
2. Deploy Contract
    - binaryPath : upload the `artifacts/onboarding_contract.wasm` file 
    - codePath: upload the `src/lib.rs` file
3. Signature Verification API
4. Subscribe Contract - `/api/subscribe-smart-contract` (Pass contract hash here)
5. Register Callback URL
    smartContractToken: Contract Hash
    callbackUrl : http://localhost:8082/api/onboard-infra (NOTE: the reason its 8082, because the DApp is configured to run on 8082. If you want the DApp server to run on a different port, change the port here as well)

## Run the dapp server

1. Move inside the `dapp` directory

1.5 Reset the `provider_info.json` to `[]` (Delete all the objects inside this file)

2. Change the values provided in `.env` file. `DID_PATH` represents the full path of `TestNetDID` directory under your Rubix Node, `SELF_CONTRACT_HASH` represents the Smart Contract Hash of the Onboarding Contract
3. Run the DApp server: `go run .`

## Executing the Smart Contract

Pass the following stringfied version in the `smartContractData` param of `/api/execute-smart-contract`:

```
{"onboard_provider": {"provider_info": {"providerDid":"bafybmi123","storage":"400GB","memory":"16GB","os":"Ubuntu","core":"4","processor":"Intel Xeon"}}}
```

All the inputs for `provider_info` are from the TRIE website

## Getting Provider information for TRIE website

To get the list of registered provider DID information, call the following DAPP server's API:

`GET: http://localhost:8082/api/provider-info` (Refer `dapp/main.go`)
