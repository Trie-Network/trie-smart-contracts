use serde::{Deserialize, Serialize};
use rubixwasm_std::{contract_fn};
use rubixwasm_std::errors::WasmError;

#[derive(Deserialize, Serialize)]
pub struct RateAssetReq {
    pub asset_id: String,
    pub user_did: String,
    pub rating: u8,
}

#[contract_fn]
pub fn rate_asset(rate_asset_req: RateAssetReq) -> Result<RateAssetReq, WasmError> {
    // The contract's core logic involves only writing the rating input to the
    // contract's token chain.
    // The overall rating of a contract is determined by looping over the tokenchain
    // externally
    Ok(rate_asset_req)
}
