package login

// #cgo CFLAGS: -g
// #cgo LDFLAGS: -lpam -lpam_misc
// #include <stdlib.h>
// #include <pwd.h>
// #include <grp.h>
// #include <security/pam_appl.h>
// #include <login.h>
import "C"
import (
	"aporia/ansi"
	"aporia/constants"
	"errors"
	"unsafe"
)

func Authenticate(username string, password string) error {
	var handle *C.struct_pam_handle
	usernameStr := C.CString(username)
	serviceStr := C.CString(constants.PamService)
	passwordStr := C.CString(password)
	conv := C.new_conv(passwordStr)

	defer C.free(unsafe.Pointer(usernameStr))
	defer C.free(unsafe.Pointer(serviceStr))
	defer C.free(unsafe.Pointer(passwordStr))

	{
		ret := C.pam_start(serviceStr, usernameStr, &conv, &handle)

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

	{
		ret := C.pam_acct_mgmt(handle, 0)
		if ret != C.PAM_SUCCESS {
			return errors.New("pam_acct_mgmt")
		}
	}

	pwnam := C.getpwnam(usernameStr)
	C.initgroups(pwnam.pw_name, pwnam.pw_gid)

	// Child shell must be cleared here
	ansi.Clear()

	{
		ret := C.pam_setcred(handle, C.PAM_ESTABLISH_CRED)
		if ret != C.PAM_SUCCESS {
			return errors.New("pam_setcred")
		}
	}

	{
		ret := C.pam_open_session(handle, 0)
		if ret != C.PAM_SUCCESS {
			C.pam_setcred(handle, C.PAM_DELETE_CRED)
			return errors.New("pam_open_session " + pamReason(ret))
		}
	}

	launchShell(handle, pwnam)

	return nil
}

func pamReason(err C.int) string {
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
}

func closePamSession(handle *C.struct_pam_handle) {
	C.pam_setcred(handle, C.PAM_DELETE_CRED)
	C.pam_close_session(handle, 0)
}
