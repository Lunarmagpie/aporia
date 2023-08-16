package constants

const ConfigDir = "/etc/aporia/"
const AsciiFileExt = "ascii"
const LastSessionFile = "/etc/aporia/.last-session"
const PamService = "aporia"
const PamConfDir = "/etc/pam.d"
const XSessionsPath = "/etc/X11/Xsession.d"

const X11SessionsDir = "/usr/share/xsessions"
const WaylandSessionsDir = "/usr/share/wayland-sessions"

const AporiaStartxPath = "APORIA_STARTX_PATH"
const AporiaExec = "APORIA_EXEC"
const X11StartupCommand = "/etc/aporia/.scripts/startx.sh /etc/aporia/.scripts/xsetup.sh"

// Ascii art used when there is no config
const DefaultAsciiArt = ``

func DefaultMessages() []string {
	return []string{"Login:"}
}
