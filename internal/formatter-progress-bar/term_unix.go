//go:build linux || darwin

package formatter_progress_bar

import (
	"os"
	"strconv"
	"syscall"
	"unsafe"
)

// TermWidth returns the terminal width from INDRA_TERM_WIDTH env or TIOCGWINSZ.
func TermWidth() int {
	if v := os.Getenv("INDRA_TERM_WIDTH"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	type winsize struct{ Row, Col, Xpixel, Ypixel uint16 }
	ws := &winsize{}
	syscall.Syscall(syscall.SYS_IOCTL, os.Stderr.Fd(), syscall.TIOCGWINSZ, uintptr(unsafe.Pointer(ws)))
	if ws.Col > 0 {
		return int(ws.Col)
	}
	return 80
}
