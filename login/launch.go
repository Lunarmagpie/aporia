package login

// #cgo CFLAGS: -g
// #cgo LDFLAGS: -lpam -lpam_misc
// #include <unistd.h>
// #include <sys/wait.h>
// #include <utmp.h>
import "C"
import (
	_ "aporia/ansi"
	"aporia/constants"
	"os"
	"syscall"
)

type SessionType string

const (
	shellSession   = "tty"
	x11Session     = "x11"
	waylandSession = "wayland"
)

type Session struct {
	Name        string
	sessionType SessionType
	// The filepath to the launcher file for the session, or null if its a shell session.
	filepath *string
}

func X11Session(name string, startxPath string) Session {
	return Session{
		Name:        name,
		sessionType: x11Session,
		filepath:    &startxPath,
	}
}

func WaylandSession(name string, filepath string) Session {
	return Session{
		Name:        name,
		sessionType: waylandSession,
		filepath:    &filepath,
	}
}

func ShellSession() Session {
	return Session{
		Name:        "shell",
		sessionType: shellSession,
		filepath:    nil,
	}
}

func launch(session Session, pam_handle *C.struct_pam_handle, pwnam *C.struct_passwd) {
	pid := C.fork()

	if pid == 0 {
		// Child
		becomeUser(pwnam)
		shell := C.GoString(pwnam.pw_shell)
		env := makeEnv(pam_handle, pwnam, session.Name, string(session.sessionType))

		switch session.sessionType {
		case shellSession:
			launchShell(env, shell)
		case x11Session:
			launchX11(env, shell, *session.filepath)
		case waylandSession:
			launchWayland(env, shell, *session.filepath)
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
