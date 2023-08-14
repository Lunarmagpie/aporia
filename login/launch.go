package login

// #cgo CFLAGS: -g
// #cgo LDFLAGS: -lpam -lpam_misc
// #include <unistd.h>
// #include <sys/wait.h>
// #include <utmp.h>
import "C"
import (
	"aporia/ansi"
	"os"
	"strings"
	"syscall"
)

type SessionType int

const (
	shellSession SessionType = iota
	x11Session
	waylandSession
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

	getCommand := func() string {
		command, _ := os.ReadFile(*session.filepath)
		return strings.TrimSpace(string(command))
	}

	if pid == 0 {
		// Child
		becomeUser(pwnam)
		shell := C.GoString(pwnam.pw_shell)
		env := makeEnv(pam_handle, pwnam, session.Name)

		switch session.sessionType {
		case shellSession:
			launchShell(env, shell)
		case x11Session:
			launchX11(env, shell, getCommand())
		case waylandSession:
			launchWayland(env, shell, getCommand())
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
	ansi.Clear()
}

func launchShell(env []string, shell string) {
	syscall.Exec(shell, []string{shell}, env)
}

func launchX11(env []string, shell string, command string) {
	// Run all x11 xsession.d scripts
	// syscall.Exec(shell, []string{shell, "-c", "/etc/aporia/xsetup.sh \"startx" + startxPath + "\""}, env)
}

func launchWayland(env []string, shell string, command string) {
	command = "/etc/aporia/wsetup.sh " + command
	syscall.Exec(shell, []string{shell, "-c", command}, env)
}
