package pam

// #cgo CFLAGS: -g -Wall
// #cgo LDFLAGS: -lpam -lpam_misc
// #include <stdlib.h>
// #include <security/pam_appl.h>
// #include <pwd.h>
// #include <login.h>
import "C"
import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"unsafe"
)

const service = "aporia"

func Authenticate(username string, password string) error {
	var handle *C.struct_pam_handle
	usernameStr := C.CString(username)
	serviceStr := C.CString("aporia")
	passwordStr := C.CString(password)
	conv := C.new_conv(passwordStr)

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
			return errors.New("pam_open_session")
		}
	}

	pwnam := C.getpwnam(usernameStr)

	initEnv(handle, pwnam)

	C.free(unsafe.Pointer(usernameStr))
	C.free(unsafe.Pointer(serviceStr))
	C.free(unsafe.Pointer(passwordStr))

	launch(pwnam)

	return nil
}

func launch(pwnam *C.struct_passwd) {
	shell := C.GoString(pwnam.pw_shell)

	cmd := exec.Command(shell, "-c", "/usr/bin/tput reset")
	if err := cmd.Run(); err != nil {
		log.Fatal(cmd)
	}

	cmd.Wait()
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
}

func initEnv(handle *C.struct_pam_handle, pwnam *C.struct_passwd) {
	homeDir := C.GoString(pwnam.pw_dir)
	xauthority := fmt.Sprintf(homeDir, "/", ".Xauthority")

	pamSetEnv(handle, "HOME", homeDir)
	pamSetEnv(handle, "PWD", C.GoString(pwnam.pw_dir))
	pamSetEnv(handle, "SHELL", C.GoString(pwnam.pw_shell))
	pamSetEnv(handle, "USER", C.GoString(pwnam.pw_name))
	pamSetEnv(handle, "LOGNAME", C.GoString(pwnam.pw_name))
	pamSetEnv(handle, "PATH", "/usr/local/sbin:/usr/local/bin:/usr/bin")
	pamSetEnv(handle, "XAUTHORITY", xauthority)
}

func pamSetEnv(handle *C.struct_pam_handle, k string, v string) {
	set := C.CString(fmt.Sprint(k, "=", v))
	C.pam_putenv(handle, set)
	C.free(unsafe.Pointer(set))
}
