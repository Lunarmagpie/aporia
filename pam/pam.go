package pam

// #cgo CFLAGS: -g -Wall
// #cgo LDFLAGS: -lpam -lpam_misc
// #include <stdlib.h>
// #include <security/pam_appl.h>
// #include <login.h>
import "C"
import "unsafe"
import (
	"errors"
)

const service = "aporia"

func Authenticate(username string, password string) error {
	var handle *C.struct_pam_handle
	password_str := C.CString(password)
	conv := C.new_conv(password_str)

	{
		ret := C.pam_start(C.CString(service), C.CString(username), &conv, &handle)

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
	C.free(unsafe.Pointer(password_str))
	return nil
}

func diagnose(err C.int) string {
	switch err {
	case C.PAM_ACCT_EXPIRED:
		return "PAM_ACCT_EXPIRED"
	case C.PAM_AUTH_ERR:
		return "PAM_AUTH_ERR"
	case C.PAM_AUTHINFO_UNAVAIL:
		return "PAM_AUTHINFO_UNAVAIL"
	case C.PAM_BUF_ERR:
		return "PAM_BUF_ERR"
	case C.PAM_CRED_ERR:
		return "PAM_CRED_ERR"
	case C.PAM_CRED_EXPIRED:
		return "PAM_CRED_EXPIRED"
	case C.PAM_CRED_INSUFFICIENT:
		return "PAM_CRED_INSUFFICIENT"
	case C.PAM_CRED_UNAVAIL:
		return "PAM_CRED_UNAVAIL"
	case C.PAM_MAXTRIES:
		return "PAM_MAXTRIES"
	case C.PAM_NEW_AUTHTOK_REQD:
		return "PAM_NEW_AUTHTOK_REQD"
	case C.PAM_PERM_DENIED:
		return "PAM_PERM_DENIED"
	case C.PAM_SESSION_ERR:
		return "PAM_SESSION_ERR"
	case C.PAM_SYSTEM_ERR:
		return "PAM_SYSTEM_ERR"
	case C.PAM_USER_UNKNOWN:
		return "PAM_USER_UNKNOWN"
	case C.PAM_ABORT:
		return "ABORT lol"
	default:
		return "Unknown Error"
	}
	return "Error"
}
