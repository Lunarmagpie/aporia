package constants

const ConfigDir = "/etc/aporia/"
const AsciiFileExt = "ascii"

const PamService = "aporia"
const PamConfDir = "/etc/pam.d"

// Ascii art used when there is no config
const DefaultAsciiArt = ``
func DefaultMessages() []string {
	return []string{"Login:"}
}
