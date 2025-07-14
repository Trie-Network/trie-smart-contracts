use rubixwasm_std::errors::WasmError;
use serde::{Deserialize, Serialize};
use rubixwasm_std::contract_fn;
pub mod helpers;
pub mod imports;
pub use helpers::{call_do_verify_action};

#[derive(Serialize, Deserialize)]
pub struct OnboardProviderReq {
    pub provider_did: String,
    pub provider_info: ProviderInfo,
}

#[derive(Serialize, Deserialize)]
pub struct ProviderInfo {
    pub storage: String,
    pub memory: String,
    pub os: String,
    pub core: String,
    pub gpu: String,
    pub processor: String,
    pub region: String,
    pub providerDid: String,
    pub platformName: String,
    pub providerName: String,
    pub platformImageUri: String,
    pub hostingCost: u64,
    pub trainingCost: u64
}

fn output_msg(success_msg: String, err_msg: String) -> String {
    return format!("msg: {}, err: {}", success_msg, err_msg)
}

#[contract_fn]
pub fn onboard_provider(onboard_provider_req: OnboardProviderReq) -> Result<String, WasmError> {
    let provider_info = onboard_provider_req.provider_info;

    match call_do_verify_action() {
        Ok(_) => return Ok(output_msg("Success".to_string(), "".to_string())),
        Err(e) => {
            if e.msg.is_empty() {
                return Ok(output_msg("Fail".to_string(), "".to_string()))
            } else {
                return Ok(output_msg("Fail".to_string(), e.msg.to_string()))
            }
        },
    };
}

