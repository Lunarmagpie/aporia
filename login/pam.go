package login

// #cgo CFLAGS: -g
// #cgo LDFLAGS: -lpam -lpam_misc
// #include <stdlib.h>
// #include <pwd.h>
// #include <grp.h>
// #include <security/pam_appl.h>
// #include <security/pam_misc.h>
// #include <login.h>
// #include <utils.h>
import "C"
import (
	"aporia/ansi"
	"aporia/config"
	"aporia/constants"
	"errors"
	"fmt"
	"unsafe"
)


var set_message func (string)

//export handlePamErrorMessage
func handlePamErrorMessage(msg *C.char) {
	set_message(C.GoString(msg))
}

//export handlePamTextInfo 
func handlePamTextInfo (msg *C.char) {
	set_message(C.GoString(msg))
}

func Authenticate(username string, password string, session config.Session, set_message_ func(string)) error {
	var handle *C.struct_pam_handle
	usernameStr := C.CString(username)
	serviceStr := C.CString(constants.PamService)
	passwordStr := C.CString(password)
	conv := C.new_conv(passwordStr)

	defer C.free(unsafe.Pointer(usernameStr))
	defer C.free(unsafe.Pointer(serviceStr))
	defer C.free(unsafe.Pointer(passwordStr))
	defer C.free(unsafe.Pointer(handle))

	set_message = set_message_

	set_message("Authenticating...")

	{
		ret := C.pam_start(serviceStr, usernameStr, &conv, &handle)

		if ret != C.PAM_SUCCESS {
			return errors.New("Could not start pam session: " + pamReason(ret))
		}
	}

	{
		ret := C.pam_authenticate(handle, 0)
		if ret != C.PAM_SUCCESS {
			return errors.New("Could not authenticate user: " + pamReason(ret))
		}
	}

	{
		ret := C.pam_acct_mgmt(handle, 0)
		if ret != C.PAM_SUCCESS {
			return errors.New("Account is not valid: " + pamReason(ret))
		}
	}

	pwnam := C.getpwnam(usernameStr)
	C.initgroups(pwnam.pw_name, pwnam.pw_gid)

	// Child shell must be cleared here
	ansi.Clear()
	fmt.Println("Setting credentials...")
	{
		ret := C.pam_setcred(handle, C.PAM_ESTABLISH_CRED)
		if ret != C.PAM_SUCCESS {
			return errors.New("pam_setcred: " + pamReason(ret))
		}
	}

	fmt.Println("Opening session...")
	// This is where the bug is happening.
	{
		// Silenced to hide the distro's login message
		ret := C.pam_open_session(handle, 1)
		if ret != C.PAM_SUCCESS {
			C.pam_setcred(handle, C.PAM_DELETE_CRED)
			return errors.New("pam_open_session: " + pamReason(ret))
		}
		fmt.Println("Session opened successfully.")
	}

	// Login was successful, so lets save the choices for next time.
	fmt.Println("Saving last session...")
	config.SaveSession(session.Name, username)

	fmt.Println("Launching session...")
	launch(session, handle, pwnam)

	fmt.Println("Session closed, restarting Aporia")
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
	C.pam_close_session(handle, 0)
	result := C.pam_setcred(handle, C.PAM_DELETE_CRED)
	C.pam_end(handle, result)
}
