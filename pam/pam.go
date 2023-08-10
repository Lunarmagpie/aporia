package pam

// #cgo CFLAGS: -g -Wall
// #cgo LDFLAGS: -lpam -lpam_misc
// #include <security/pam_appl.h>
// #include <login.h>
import "C"
import "errors"

const service = "aporia"

func Authenticate(user string, password string) error {
	var handle *C.struct_pam_handle
	var conv = C.new_conv(C.CString(user), C.CString(password));

	{
		ret := C.pam_start(C.CString(service), C.CString(user), &conv, &handle)

		if ret != C.PAM_SUCCESS {
			return errors.New("Could not start pam session.")
		}
	}
	
	{
		ret := C.pam_authenticate(handle, 0)
		if ret != C.PAM_SUCCESS {
			return errors.New("Could not authenticate user.")
		}
	}


	return nil
}
