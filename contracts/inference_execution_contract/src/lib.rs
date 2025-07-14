pub mod helpers;
pub mod imports;

use helpers::{call_execute_nft_api, ExecuteNFT};
use rubixwasm_std::{contract_fn, errors::WasmError};
use serde::{Deserialize, Serialize};


#[derive(Serialize, Deserialize)]
pub struct StoreInferenceReq {
    pub asset_id: String,
    pub inference_info: String,
    pub asset_value: f64,

    pub depin_did: String,
}

#[contract_fn]
pub fn store_inference(req: StoreInferenceReq) -> Result<String, WasmError> {
    // Call the transfer FT API with the provided request
    let credit_purchase_transfer_result = call_execute_nft_api(ExecuteNFT {
        comment: "Inference Storage in Asset".to_string(),
        nft: req.asset_id,
        nft_data: req.inference_info,
        nft_value: req.asset_value as f64,
        executor: req.depin_did,
        receiver: "".to_string(),
    });

    match credit_purchase_transfer_result {
        Ok(_) => return Ok("Inference stored successfully".to_string()),
        Err(e) => return Err(WasmError::from(format!("Failed to transfer FT: {}", e.msg))),
    };
}

