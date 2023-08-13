package login

// #include <stdlib.h>
// #include <unistd.h>
// #include <string.h>
// #include <utmp.h>
// #include <utmpx.h>
import "C"
import (
	"strings"
	"unsafe"
)

func addUtmpEntry(entry *C.struct_utmp, pwnam *C.struct_passwd, pid C.int) {
	entry.ut_type = C.USER_PROCESS
	entry.ut_pid = pid

	ttynameString := C.GoString(C.ttyname(C.STDIN_FILENO))

	C.strcpy((*C.char)(unsafe.Pointer(&entry.ut_line)), C.CString(strings.TrimPrefix(ttynameString, "/dev/")))
	C.strcpy((*C.char)(unsafe.Pointer(&entry.ut_id)), C.CString(strings.TrimPrefix(ttynameString, "/dev/tty")))
	C.strcpy((*C.char)(unsafe.Pointer(&entry.ut_user)), pwnam.pw_name)
	C.memset(unsafe.Pointer(&entry.ut_host), 0, C.UT_HOSTSIZE)

	C.setutent()
	C.pututline(entry)
}

func removeUtmpEntry(entry *C.struct_utmp) {
	entry.ut_type = C.DEAD_PROCESS
	C.memset(unsafe.Pointer(&entry.ut_line), 0, C.UT_LINESIZE)
	C.memset(unsafe.Pointer(&entry.ut_host), 0, C.UT_NAMESIZE)

	C.setutent()
	C.pututline(entry)
	C.endutent()
}
