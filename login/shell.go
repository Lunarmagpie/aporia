package login

// #include <unistd.h>
// #include <sys/wait.h>
// #include <utmp.h>
import "C"
import (
	"os"
	"syscall"
)

func launchShell(pam_handle *C.struct_pam_handle, pwnam *C.struct_passwd) {
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

	closePamSession(pam_handle)

	var status C.int
	C.waitpid(pid, &status, 0)

	removeUtmpEntry(&utmpEntry)
}
