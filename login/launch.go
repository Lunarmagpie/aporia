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
			launchX11(env, shell, *session.Filepath)
		case config.WaylandSession:
			launchWayland(env, shell, *session.Filepath)
		}

	}

	// Parent
	utmpEntry := C.struct_utmp{}
	addUtmpEntry(&utmpEntry, pwnam, pid)

	var status C.int
	C.waitpid(pid, &status, 0)

	closePamSession(pam_handle)
	removeUtmpEntry(&utmpEntry)

	// Clear the screen to prevent X11 weirdities
	// ansi.Clear()
}

func launchShell(env []string, shell string) {
	syscall.Exec(shell, []string{shell}, env)
	os.Exit(0)
}

func launchX11(env []string, shell string, filepath string) {
	env = append(env, constants.AporiaStartxPath+"="+filepath)
	syscall.Exec(shell, []string{shell, "-c", constants.X11StartupCommand}, env)
	os.Exit(0)
}

func launchWayland(env []string, shell string, filepath string) {
	syscall.Exec(shell, []string{shell, "-c", filepath}, env)
	os.Exit(0)
}
