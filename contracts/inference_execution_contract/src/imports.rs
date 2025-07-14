extern "C" {
    pub fn do_execute_nft(
        inputdata_ptr: *const u8,
        inputdata_len: usize,
        resp_ptr_ptr: *mut *const u8,
        resp_len_ptr: *mut usize,
    ) -> i32;
}