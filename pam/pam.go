package pam

// #cgo CFLAGS: -g
// #cgo LDFLAGS: -lpam -lpam_misc
// #include <stdlib.h>
// #include <unistd.h>
// #include <sys/wait.h>
// #include <grp.h>
// #include <security/pam_appl.h>
// #include <pwd.h>
// #include <string.h>
// #include <utmp.h>
// #include <utmpx.h>
// #include <sys/stat.h>
// #include <login.h>
// #include <utils.h>
import "C"
import (
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

func Authenticate(username string, password string) error {
	var handle *C.struct_pam_handle
	usernameStr := C.CString(username)
	serviceStr := C.CString("aporia")
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
	fmt.Print("\033[H\033[0J")

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

	launch(handle, pwnam)

	return nil
}

func launch(pam_handle *C.struct_pam_handle, pwnam *C.struct_passwd) {
	pid := C.fork()

	if pid == 0 {
		// Child
		becomeUser(pwnam)
		shell := C.GoString(pwnam.pw_shell)
		initEnv(pam_handle, pwnam)
		syscall.Exec(shell, []string{shell}, os.Environ())
	}

	// Parent
	utmpEntry := C.struct_utmp{}
	addUtmpEntry(&utmpEntry, pwnam, pid)

	C.pam_close_session(pam_handle, 0)

	var status C.int
	C.waitpid(pid, &status, 0)

	removeUtmpEntry(&utmpEntry)
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

func becomeUser(pwnam *C.struct_passwd) error {
	if 0 != C.chdir(pwnam.pw_dir) {
		return errors.New("chdir")
	}
	if 0 != C.setgid(pwnam.pw_gid) {
		return errors.New("setgid")
	}
	if 0 != C.setuid(pwnam.pw_uid) {
		return errors.New("setuid")
	}
	if 0 != C.initgroups(pwnam.pw_name, pwnam.pw_gid) {
		return errors.New("initgroups")
	}
	return nil
}

func foregroundProcessGroup() error {
	pid := C.getpid()
	if 0 != C.setpgid(pid, pid) {
		return errors.New("setgpid")
	}
	if 0 != C.tcsetpgrp(C.STDIN_FILENO, pid) {
		return errors.New("tcsetpgrp")
	}
	return nil
}


func initEnv(pam_handle *C.struct_pam_handle, pwnam *C.struct_passwd) {
	homeDir := C.GoString(pwnam.pw_dir)

	os.Clearenv()

	os.Setenv("HOME", homeDir)
	os.Setenv("PWD", homeDir)
	os.Setenv("SHELL", C.GoString(pwnam.pw_shell))
	os.Setenv("USER", C.GoString(pwnam.pw_name))
	os.Setenv("LOGNAME", C.GoString(pwnam.pw_name))

	_, found := os.LookupEnv("TERM")
	if !found {
		os.Setenv("TERM", "linux")
	}

	os.Setenv("PATH", "/usr/local/sbin:/usr/local/bin:/usr/bin:/usr/sbin:/sbin")

	pamEnvList := C.pam_getenvlist(pam_handle)

	for _, v := range cArrayToGoSlice(pamEnvList) {
		l := strings.Split(v, "=")
		os.Setenv(l[0], l[1])
	}

	C.free(unsafe.Pointer(pamEnvList))
}

func cArrayToGoSlice(arr **C.char) []string {
	var envs []string

	for i := C.int(0); C.index_string_array(arr, i) != nil; i++ {
		nextString := C.GoString(C.index_string_array(arr, i))
		envs = append(envs, nextString)
	}

	return envs
}

func addUtmpEntry(entry *C.struct_utmp, pwnam *C.struct_passwd, pid C.int) {
	entry.ut_type = C.USER_PROCESS
	entry.ut_pid = pid

	ttynameString := C.GoString(C.ttyname(C.STDIN_FILENO))

	C.strcpy((*C.char)(unsafe.Pointer(&entry.ut_line)), C.CString(strings.TrimPrefix(ttynameString, "/dev/")))
	C.strcpy((*C.char)(unsafe.Pointer(&entry.ut_id)), C.CString(strings.TrimPrefix(ttynameString, "/dev/tty")))
	C.strcpy((*C.char)(unsafe.Pointer(&entry.ut_user)), pwnam.pw_name)
	C.memset(unsafe.Pointer(&entry.ut_host), 0, C.UT_HOSTSIZE)

	C.setutent()
	C.pututline(entry)
}

func removeUtmpEntry(entry *C.struct_utmp) {
	entry.ut_type = C.DEAD_PROCESS
	C.memset(unsafe.Pointer(&entry.ut_line), 0, C.UT_LINESIZE)
	C.memset(unsafe.Pointer(&entry.ut_host), 0, C.UT_NAMESIZE)

	C.setutent()
	C.pututline(entry)
	C.endutent()
}
