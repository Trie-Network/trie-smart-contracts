use std::ffi::CString;
use rubixwasm_std::errors::WasmError;
use super::imports::{do_verify_action};
use std::slice;
use std::str;
use serde::{Serialize,Deserialize};
use serde_json;


pub fn call_do_verify_action() -> Result<String, WasmError> {
    unsafe {
        // Response vars
        let mut signature_verification_response_ptr: *const u8 = std::ptr::null();
        let mut signature_verification_response_len: usize = 0;


        // Call the imported host function
        let result = do_verify_action(
            &mut signature_verification_response_ptr,
            &mut signature_verification_response_len
        );
        
        if result != 0 {
            return Err(WasmError::from(format!("Host function returned error code {}", result)));
        }

        if signature_verification_response_ptr.is_null() {
            return Err(WasmError::from("Signature Response pointer is null".to_string()));
        }

        let response_slice = slice::from_raw_parts(signature_verification_response_ptr, signature_verification_response_len);
        match str::from_utf8(response_slice) {
            Ok(res) => Ok(res.to_string()),
            Err(_) => Err(WasmError::from("invalid utf-8 reponse".to_string()))
        }
    }
}

