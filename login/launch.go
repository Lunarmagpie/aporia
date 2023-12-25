package login

// #cgo CFLAGS: -g
// #cgo LDFLAGS: -lpam -lpam_misc
// #include <unistd.h>
// #include <sys/wait.h>
// #include <utmp.h>
import "C"
import (
	"aporia/config"
	"aporia/constants"
	"os"
	"syscall"
)

func launch(session config.Session, pam_handle *C.struct_pam_handle, pwnam *C.struct_passwd) {
	pid := C.fork()

	if pid == 0 {
		// Child
		becomeUser(pwnam)
		shell := C.GoString(pwnam.pw_shell)
		env := makeEnv(pam_handle, pwnam, session.Name, string(session.SessionType))

		switch session.SessionType {
		case config.ShellSession:
			launchShell(env, shell)
		case config.X11Session:
			launchX11(env, shell, session.Exec, session.Filepath)
		case config.WaylandSession:
			launchWayland(env, shell, session.Exec, session.Filepath)
		}
	}

	// Parent
	utmpEntry := C.struct_utmp{}
	addUtmpEntry(&utmpEntry, pwnam, pid)

	// Wait for all child processes to finish
	var status C.int
	for C.waitpid(-1, &status, 0) > 0 {}

	closePamSession(pam_handle)
	removeUtmpEntry(&utmpEntry)
}

func launchShell(env []string, shell string) {
	syscall.Exec(shell, []string{shell}, env)
	os.Exit(0)
}

func launchX11(env []string, shell string, exec *string, filepath *string) {
	if filepath != nil {
		env = append(env, constants.AporiaStartxPath+"="+*filepath)
	} else {
		env = append(env, constants.AporiaExec+"="+*exec)
	}
	syscall.Exec(shell, []string{shell, "-c", constants.X11StartupCommand}, env)
	os.Exit(0)
}

func launchWayland(env []string, shell string, exec *string, filepath *string) {
	if filepath != nil {
		syscall.Exec(shell, []string{shell, "-c", *filepath}, env)
	} else {
		syscall.Exec(shell, []string{shell, "-c", *exec}, env)
	}
	os.Exit(0)
}
