extern "C" {
    pub fn do_verify_action(
        response_ptr_ptr: *mut *const u8,
        response_len_ptr: *mut usize,
    ) -> i32;
}
