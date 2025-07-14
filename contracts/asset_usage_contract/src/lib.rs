pub mod import;
pub mod helper;

use rubixwasm_std::errors::WasmError;
use helper::{call_transfer_ft_api, TransferFt, ExecuteNFT, call_execute_nft_api};
use serde::{Deserialize, Serialize};
use rubixwasm_std::contract_fn;

const FT_NAME: &str = "TRIE";

#[derive(Serialize, Deserialize)]
pub struct UseAssetReq {
    asset_usage_price: i32,
    asset_usage_purpose: String,
    asset_denom: String,
    asset_user_did: String,
    asset_owner_did: String,
    asset_id: String,
    asset_value: f64,

    ft_denom_creator: String,
}


#[contract_fn]
pub fn use_asset(use_asset_req: UseAssetReq) -> Result<String, WasmError> {
    let asset_usage_description = use_asset_req.asset_usage_purpose;
    let asset_owner = use_asset_req.asset_owner_did;
    let user_did = use_asset_req.asset_user_did;
    
    // Buyer will perform FT transfer to the actual owner of the contract
    let buyer_depin_transfer_req = TransferFt {
        comment: asset_usage_description.clone(),
        ft_count: use_asset_req.asset_usage_price,
        ft_name: use_asset_req.asset_denom,
        creatorDID: use_asset_req.ft_denom_creator,
        sender: user_did.clone(),
        receiver: asset_owner.clone(),
    };

    match call_transfer_ft_api(buyer_depin_transfer_req) {
        Ok(_) => {},
        Err(e) => return Err(WasmError::from(format!("unable to transfer tokens from usage to DePIN provider, err: {}", e.msg))),
    };

    let nft_usage_request = ExecuteNFT {
        comment: asset_usage_description.clone(),
        nft: use_asset_req.asset_id.clone(),
        nft_data: asset_usage_description.clone(),
        nft_value: use_asset_req.asset_value,
        receiver: "".to_string(),
        executor: user_did.clone(),
    };

    match call_execute_nft_api(nft_usage_request) {
        Ok(_) => return Ok(use_asset_req.asset_id),
        Err(e) => return Err(WasmError::from(format!("unable to execute NFT for usage, err: {}", e.msg))),
    };
}