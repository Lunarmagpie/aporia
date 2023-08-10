package pam

// #cgo CFLAGS: -g -Wall
// #cgo LDFLAGS: -lpam -lpam_misc
// #include <security/pam_appl.h>
// #include <security/pam_misc.h>
//
// struct pam_conv conv = {
// 		misc_conv,
// 		NULL
// };
//
// const void *to_unsafe_pointer(const char *c) {
// 		return (const void*) c;
// }
//
import "C"
import "errors"

const service = "aporia"

func Authenticate(user string, password string) error {
	var handle *C.struct_pam_handle

	ret := C.pam_start(C.CString(service), C.CString(user), &C.conv, &handle)

	if ret != C.PAM_SUCCESS {
		return errors.New("Could not start pam session.")
	}

	C.pam_set_item(handle, C.PAM_AUTHTOK, C.to_unsafe_pointer(C.CString(password)))

	return nil
}
