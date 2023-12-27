package constants

const ConfigDir = "/etc/aporia/"
const ConfigFile = "config"
const AsciiFileExt = "ascii"
const LastSessionFile = "/etc/aporia/.last-session"
const PamService = "aporia"
const PamConfDir = "/etc/pam.d"
const XSessionsPath = "/etc/X11/Xsession.d"

const X11SessionsDir = "/usr/share/xsessions"
const WaylandSessionsDir = "/usr/share/wayland-sessions"

const AporiaStartxPath = "APORIA_STARTX_PATH"
const AporiaExec = "APORIA_EXEC"
const X11StartupCommand = "exec /bin/bash --login /etc/aporia/.scripts/startx.sh /etc/aporia/.scripts/xsetup.sh"

var ShutdownCommand = []string{"/bin/bash", "-c", "shutdown now"}
var RebootCommand = []string{"/bin/bash", "-c", "reboot"}

// Ascii art used when there is no config
const DefaultAsciiArt = `Check the docs on how to add an ascii art!`

func DefaultMessages() []string {
	return []string{"Login:"}
}
