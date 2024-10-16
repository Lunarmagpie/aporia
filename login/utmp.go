package login 
// #cgo CFLAGS: -g
// #cgo LDFLAGS: -lpam -lpam_misc
// #include <stdlib.h>
// #include <unistd.h>
// #include <string.h>
// #include <sys/time.h>
// #include <utmp.h>
// #include <utmpx.h>
import "C"
import (
	"strings"
	"unsafe"
)

func addUtmpEntry(entry *C.struct_utmp, pwnam *C.struct_passwd, pid C.int) {
	entry.ut_type = C.LOGIN_PROCESS
	entry.ut_pid = pid

	setEntryTime(entry)
	ttynameString := C.GoString(C.ttyname(C.STDIN_FILENO))

	C.strcpy((*C.char)(unsafe.Pointer(&entry.ut_line)), C.CString(strings.TrimPrefix(ttynameString, "/dev/")))
	C.strcpy((*C.char)(unsafe.Pointer(&entry.ut_id)), C.CString(strings.TrimPrefix(ttynameString, "/dev/tty")))
	C.strcpy((*C.char)(unsafe.Pointer(&entry.ut_user)), pwnam.pw_name)
	C.memset(unsafe.Pointer(&entry.ut_host), 0, C.UT_HOSTSIZE)

	entry.ut_addr_v6[0] = 0

	C.setutent()
	C.pututline(entry)
}

func setEntryTime(entry *C.struct_utmp) {
	var tv C.struct_timeval
	C.gettimeofday(&tv, C.NULL)
	entry.ut_tv.tv_sec = C.uint(tv.tv_sec)
	entry.ut_tv.tv_usec = C.int(tv.tv_usec)
}

func removeUtmpEntry(entry *C.struct_utmp) {
	entry.ut_type = C.DEAD_PROCESS
	setEntryTime(entry)

	C.memset(unsafe.Pointer(&entry.ut_line), 0, C.UT_LINESIZE)
	C.memset(unsafe.Pointer(&entry.ut_host), 0, C.UT_NAMESIZE)

	C.setutent()
	C.pututline(entry)
	C.endutent()
}
