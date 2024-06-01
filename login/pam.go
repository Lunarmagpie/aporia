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
			return errors.New(pamReasonToString(handle, ret))
		}
	}

	{
		ret := C.pam_authenticate(handle, 0)
		if ret != C.PAM_SUCCESS {
			return errors.New(pamReasonToString(handle, ret))
		}
	}

	{
		ret := C.pam_acct_mgmt(handle, 0)
		if ret != C.PAM_SUCCESS {
			return errors.New(pamReasonToString(handle, ret))
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
			return errors.New(pamReasonToString(handle, ret))
		}
	}

	fmt.Println("Opening session...")
	// This is where the bug is happening.
	{
		// Silenced to hide the distro's login message
		ret := C.pam_open_session(handle, 1)
		if ret != C.PAM_SUCCESS {
			C.pam_setcred(handle, C.PAM_DELETE_CRED)
			return errors.New("pam_open_session: " + pamReasonToString(handle, ret))
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

func pamReasonToString(handle *C.struct_pam_handle, err C.int) string {
	return C.GoString(C.pam_strerror(handle, err))
}

func closePamSession(handle *C.struct_pam_handle) {
	C.pam_close_session(handle, 0)
	result := C.pam_setcred(handle, C.PAM_DELETE_CRED)
	C.pam_end(handle, result)
}
