use super::imports::{do_execute_nft};
use std::str;
use serde::{Serialize,Deserialize};
use rubixwasm_std::errors::WasmError;
use std::slice;

#[derive(Serialize, Deserialize)]
pub struct ExecuteNFT{
    pub comment:    String, 
    pub nft:        String,
    pub nft_data:   String,
    pub nft_value:  f64,
    pub executor:      String,
    pub receiver:    String,
}

pub fn call_execute_nft_api(input_data: ExecuteNFT) -> Result<String, WasmError> {
    unsafe {
        // Convert the input data to bytes
        let input_bytes = serde_json::to_string(&input_data).unwrap().into_bytes();

        // let input_bytes = input_data.as_bytes();
        let input_ptr = input_bytes.as_ptr();
        let input_len = input_bytes.len();

        // Allocate space for the response pointer and length
        let mut resp_ptr: *const u8 = std::ptr::null();
        let mut resp_len: usize = 0;

        // Call the imported host functionrubixwasm_std::
        let result = do_execute_nft(
            input_ptr,
            input_len,
            &mut resp_ptr,
            &mut resp_len,
        );
        
        if result != 0 {
            return Err(WasmError::from(format!("Host function returned error code {}", result)));
        }

        // Ensure the response pointer is not null
        if resp_ptr.is_null() {
            return Err(WasmError::from("Response pointer is null".to_string()));
        }

        // Convert the response back to a Rust String
        let response_slice = slice::from_raw_parts(resp_ptr, resp_len);
        match str::from_utf8(response_slice) {
            Ok(s) => Ok(s.to_string()),
            Err(_) => Err(WasmError::from("Invalid UTF-8 response".to_string())),
        }
    }
}

