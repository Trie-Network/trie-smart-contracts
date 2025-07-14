use super::imports::{do_transfer_ft_trie, do_add_credit};
use std::str;
use serde::{Serialize,Deserialize};
use rubixwasm_std::errors::WasmError;
use std::slice;


#[derive(Serialize, Deserialize)]
pub struct TransferFt{
    pub comment: String, 
    pub ft_count: i32,
    pub ft_name: String,
    pub creatorDID: String,
    pub sender: String,
    pub receiver: String,
}

pub fn call_transfer_ft_api(input_data: TransferFt) -> Result<String, WasmError> {
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
        let result = do_transfer_ft_trie(
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

#[derive(Serialize, Deserialize)]
pub struct AddCredit {
    pub user_did: String,
    pub credit: u32,
}

pub fn call_do_add_credit_api(user_did: String, credit: u32) -> Result<String, WasmError> {
    unsafe {
        let input_data = AddCredit {
            user_did,
            credit,
        };

        // Convert the input data to bytes
        let input_bytes = serde_json::to_string(&input_data).unwrap().into_bytes();

        let input_ptr = input_bytes.as_ptr();
        let input_len = input_bytes.len();


        // Call the imported host function
        let result = do_add_credit(
            input_ptr,
            input_len,
        );

        if result != 0 {
            return Err(WasmError::from(format!("Host function returned error code {}", result)));
        } else {
            // If the function succeeds, return a success message
            Ok("Credits added successfully".to_string())
        }
    }
}
