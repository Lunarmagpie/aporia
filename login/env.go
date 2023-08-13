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
	"os"
	"errors"
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
