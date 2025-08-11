pub mod helpers;
pub mod imports;

use helpers::{call_transfer_ft_api, TransferFt, call_do_add_credit_api, AddCredit};
use rubixwasm_std::{contract_fn, errors::WasmError};
use serde::{Deserialize, Serialize};

pub const CREDIT_COUNT: u32 = 500;
pub const CREDIT_PURCHASE_AMOUNT: i32 = 5;

#[derive(Serialize, Deserialize)]
pub struct PurchaseCreditReq {
    // Write fields related to FT Transfer

    pub purchase_token_name: String,
    pub purchase_token_creator: String,

    pub user_did: String,
    pub depin_provider: String,
    pub current_timestamp: u64
}

#[contract_fn]
pub fn purchase_credit(req: PurchaseCreditReq) -> Result<String, WasmError> {
    let is_valid_purchase_token = ["TRIE", "TRI"].contains(&&req.purchase_token_name.as_str());
    if !is_valid_purchase_token {
        return Err(WasmError::from(format!("Invalid purchase token name: {}, supported tokens are: TRIE, TRI", req.purchase_token_name)));
    }

    // Call the transfer FT API with the provided request
    let credit_purchase_transfer_result = call_transfer_ft_api(TransferFt {
        comment: "Credit Purchase".to_string(),
        ft_count: CREDIT_PURCHASE_AMOUNT,
        ft_name: req.purchase_token_name,
        creatorDID: req.purchase_token_creator,
        sender: req.user_did,
        receiver: req.depin_provider,
    });

    match credit_purchase_transfer_result {
        Ok(res) => return Ok(res),
        Err(e) => return Err(WasmError::from(format!("Failed to transfer FT: {}", e.msg))),
    };
}

#[derive(Serialize, Deserialize)]
pub struct AddCredits {
    // Write fields related to FT Transfer
    pub user_did: String,
    pub current_timestamp: u64
}

#[contract_fn]
pub fn add_credits(add_credit_req: AddCredits) -> Result<String, WasmError> {
    // Call the do_add_credit API with the provided user DID and credit count

    let add_credit_result = call_do_add_credit_api(add_credit_req.user_did.clone(), CREDIT_COUNT);

    match add_credit_result {
        Ok(_) => {
            let add_credits_response = serde_json::to_string(&add_credit_req).unwrap();
            Ok(add_credits_response)
        },
        Err(e) => Err(WasmError::from(format!("Failed to add credits for DID: {}, err: {}", &add_credit_req.user_did.to_string(), e.msg))),
    }
}