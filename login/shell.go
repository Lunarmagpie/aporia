package login

// #cgo CFLAGS: -g
// #cgo LDFLAGS: -lpam -lpam_misc
// #include <unistd.h>
// #include <sys/wait.h>
// #include <utmp.h>
import "C"
import "syscall"

func launchShell(pam_handle *C.struct_pam_handle, pwnam *C.struct_passwd) {
	pid := C.fork()

	if pid == 0 {
		// Child
		becomeUser(pwnam)
		shell := C.GoString(pwnam.pw_shell)
		syscall.Exec(shell, []string{shell}, makeEnv(pam_handle, pwnam))
	}

	// Parent
	utmpEntry := C.struct_utmp{}
	addUtmpEntry(&utmpEntry, pwnam, pid)

	var status C.int
	C.waitpid(pid, &status, 0)

	closePamSession(pam_handle)
	removeUtmpEntry(&utmpEntry)
}
