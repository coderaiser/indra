//go:build linux || darwin

package formatter_progress_bar

import (
	"os"
	"strconv"
	"syscall"
	"unsafe"
)

// winsizeReader returns the terminal row/column size. It is a variable so
// tests can inject a synthetic terminal size to exercise the resize path.
var winsizeReader = ioctlWinsize

func ioctlWinsize() (row, col uint16) {
	type winsize struct{ Row, Col, Xpixel, Ypixel uint16 }
	ws := &winsize{}
	syscall.Syscall(syscall.SYS_IOCTL, os.Stderr.Fd(), syscall.TIOCGWINSZ, uintptr(unsafe.Pointer(ws)))
	return ws.Row, ws.Col
}

// TermWidth returns the terminal width from INDRA_TERM_WIDTH env or TIOCGWINSZ.
func TermWidth() int {
	if v := os.Getenv("INDRA_TERM_WIDTH"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	_, col := winsizeReader()
	if col > 0 {
		return int(col)
	}
	return 80
}
