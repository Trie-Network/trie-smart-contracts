use serde::{Deserialize, Serialize};
use rubixwasm_std::{contract_fn};
use rubixwasm_std::errors::WasmError;

#[derive(Deserialize, Serialize)]
pub struct CreateModelCardReq {
    pub id: String,             
    pub name: String,         
    pub main_category: String,   
    pub secondary_category: String,
    pub description: String,   
    pub metrics: String, 
    pub asset_type: String,
    pub asset_id: String,
    pub owner: String,         
}

#[contract_fn]
pub fn create_model_card(create_model_card_req: CreateModelCardReq) -> Result<CreateModelCardReq, WasmError> {
    Ok(create_model_card_req)
}
