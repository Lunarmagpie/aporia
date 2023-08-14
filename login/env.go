package login

// #cgo CFLAGS: -g
// #cgo LDFLAGS: -lpam -lpam_misc
// #include <stdlib.h>
// #include <unistd.h>
// #include <grp.h>
// #include <security/pam_appl.h>
// #include <login.h>
// #include <utils.h>
import "C"
import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unsafe"
)

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

func makeEnv(pam_handle *C.struct_pam_handle, pwnam *C.struct_passwd, desktopName string) []string {
	homeDir := C.GoString(pwnam.pw_dir)

	envMap := map[string]string{}

	setEnv := func(k string, v string) {
		envMap[k] = v
	}

	setEnv("HOME", homeDir)
	setEnv("PWD", homeDir)
	setEnv("SHELL", C.GoString(pwnam.pw_shell))
	setEnv("USER", C.GoString(pwnam.pw_name))
	setEnv("LOGNAME", C.GoString(pwnam.pw_name))
	setEnv("XAUTHORITY", filepath.Join(homeDir, ".Xauthority"))
	
	termValue, found := os.LookupEnv("TERM")
	if found {
		setEnv("TERM", termValue)
	} else {
		setEnv("TERM", "linux")
	}

	setEnv("PATH", "/usr/local/sbin:/usr/local/bin:/usr/bin:/usr/sbin:/sbin")

	// XDG Variables
	user := fmt.Sprint("/run/user/", pwnam.pw_uid)
	setEnv("XDG_RUNTIME_DIR", user)
	setEnv("XDG_SESSION_CLASS", "user")
	setEnv("XDG_SESSION_ID", "1")
	setEnv("XDG_SESSION_DESKTOP", desktopName)
	setEnv("XDG_SEAT", "seat0")
	setEnv("XDG_VTNR", "1")
	setEnv("XDG_SESSION_TYPE", "x11")

	setEnv("DBUS_SESSION_BUS_ADDRESS", fmt.Sprint("unix:pth=/run/user/", pwnam.pw_uid, "/bus"))

	pamEnvList := C.pam_getenvlist(pam_handle)

	for _, v := range cArrayToGoSlice(pamEnvList) {
		l := strings.Split(v, "=")
		setEnv(l[0], l[1])
	}

	C.free(unsafe.Pointer(pamEnvList))

	envList := []string{}

	for k, v := range envMap {
		envList = append(envList, k+"="+v)
	}

	return envList
}

func cArrayToGoSlice(arr **C.char) []string {
	var envs []string

	for i := C.int(0); C.index_string_array(arr, i) != nil; i++ {
		nextString := C.GoString(C.index_string_array(arr, i))
		envs = append(envs, nextString)
	}

	return envs
}
