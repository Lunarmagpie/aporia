package login

// #cgo CFLAGS: -g
// #cgo LDFLAGS: -lpam -lpam_misc
// #include <unistd.h>
// #include <sys/wait.h>
// #include <utmp.h>
import "C"
import (
	"aporia/ansi"
	"aporia/config"
	"aporia/constants"
	"strings"
	"fmt"
	"syscall"
)

func launch(session config.Session, pam_handle *C.struct_pam_handle, pwnam *C.struct_passwd) {
	pid := C.fork()

	if pid == 0 {
		// Child
		fmt.Println("Launching " + session.SessionType + "...")
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
		default:
			fmt.Println("Invalid session type. This should never happen.")
		}
	}

	// Parent
	utmpEntry := C.struct_utmp{}
	addUtmpEntry(&utmpEntry, pwnam, pid)

	// Wait for all child processes to finish
	var status C.int
	for C.waitpid(-1, &status, 0) > 0 {}

	if (status >= 1) {
		// There was an error, so lets wait until the user reads it to continue.
		fmt.Println()
		fmt.Println("[ Script exited with code", int(status), "]", " There was an error with your Window Manager or configuration. Press enter to continue...")
		fmt.Scanln()
	} else {
		fmt.Println("Child process has closed, beginning cleanup...")
	}

	removeUtmpEntry(&utmpEntry)
	closePamSession(pam_handle)
}

func launchShell(env []string, shell string) {
	// We don't want logs showing up on shell login.
	ansi.Clear()
	syscall.Exec(shell, []string{shell}, env)
}

func launchX11(env []string, shell string, exec *string, filepath *string) {
	if filepath != nil {
		env = append(env, constants.AporiaStartxPath+"="+*filepath)
	} else {
		env = append(env, constants.AporiaExec+"="+*exec)
	}
	syscall.Exec(shell, []string{shell, "-c", constants.X11StartupCommand}, env)
}

func launchWayland(env []string, shell string, exec *string, filepath *string) {
	if filepath != nil {
		syscall.Exec(shell, []string{shell, "-c", strings.Replace(*filepath, " ", "\\ ", -1)}, env)
	} else {
		syscall.Exec(shell, []string{shell, "-c", *exec}, env)
	}
}
