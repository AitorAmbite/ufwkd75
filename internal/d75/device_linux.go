//go:build linux

package d75

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

// En Linux todo el canal de configuración (reports 4, 7 y 8) está en un único hidraw:
// el de la interfaz que declara usage page 0xFF02 con report ID 7.
var configSignature = []byte{0x06, 0x02, 0xff, 0x09, 0x01, 0xa1, 0x01, 0x85, 0x07}

var deviceMu sync.Mutex

const udevRule = `# Weikav D75: acceso a hidraw para el usuario con sesión activa (sin root)
KERNEL=="hidraw*", ATTRS{idVendor}=="0603", ATTRS{idProduct}=="0351", TAG+="uaccess"
KERNEL=="hidraw*", ATTRS{idVendor}=="0603", ATTRS{idProduct}=="0352", TAG+="uaccess"
`

const udevRulePath = "/etc/udev/rules.d/70-weikav-d75.rules"

func findDevice() (string, error) {
	dirs, _ := filepath.Glob("/sys/class/hidraw/hidraw*")
	for _, dir := range dirs {
		uevent, err := os.ReadFile(filepath.Join(dir, "device/uevent"))
		if err != nil {
			continue
		}
		desc, err := os.ReadFile(filepath.Join(dir, "device/report_descriptor"))
		if err != nil || !bytes.Contains(desc, configSignature) {
			continue
		}
		for _, pid := range productIDs {
			if strings.Contains(string(uevent), fmt.Sprintf(":%08X:%08X", vendorID, pid)) {
				return "/dev/" + filepath.Base(dir), nil
			}
		}
	}
	return "", errNotFound
}

func DeviceStatus() Status {
	path, err := findDevice()
	if err != nil {
		return Status{Detail: err.Error()}
	}
	f, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return Status{Detail: path, NeedsPermission: errors.Is(err, os.ErrPermission)}
	}
	f.Close()
	return Status{Connected: true, Detail: path}
}

func hidiocsfeature(n int) uintptr {
	const iocWrite, iocRead = 1, 2
	return uintptr((iocWrite|iocRead)<<30 | n<<16 | 'H'<<8 | 0x06)
}

func sendReports(reports ...[]byte) error {
	deviceMu.Lock()
	defer deviceMu.Unlock()
	path, err := findDevice()
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return newError(ErrOpen, path, err)
	}
	defer f.Close()
	for _, buf := range reports {
		_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), hidiocsfeature(len(buf)), uintptr(unsafe.Pointer(&buf[0])))
		if errno != 0 {
			return newError(ErrSend, fmt.Sprintf("report %d cmd %d", buf[0], buf[1]), errno)
		}
		time.Sleep(reportDelay)
	}
	return nil
}

// ReadKeyEvents lee pulsaciones físicas hasta que se cancele ctx; reconecta solo.
func ReadKeyEvents(ctx context.Context, fn func(KeyEvent)) {
	buf := make([]byte, 64)
	for ctx.Err() == nil {
		path, err := findDevice()
		var f *os.File
		if err == nil {
			f, err = os.Open(path)
		}
		if err != nil {
			sleepCtx(ctx, time.Second)
			continue
		}
		stop := context.AfterFunc(ctx, func() { f.Close() })
		for {
			n, err := f.Read(buf)
			if err != nil {
				break
			}
			if ev, ok := parseKeyEvent(buf[:n]); ok {
				fn(ev)
			}
		}
		stop()
		f.Close()
		sleepCtx(ctx, time.Second)
	}
}

// InstallPermissions instala la regla udev con pkexec (pide la contraseña con el diálogo del sistema).
func InstallPermissions() error {
	cmd := exec.Command("pkexec", "sh", "-c",
		fmt.Sprintf("cat > %s && udevadm control --reload && udevadm trigger --subsystem-match=hidraw", udevRulePath))
	cmd.Stdin = strings.NewReader(udevRule)
	if out, err := cmd.CombinedOutput(); err != nil {
		return newError(ErrPermission, strings.TrimSpace(string(out)), err)
	}
	time.Sleep(time.Second) // udev tarda un momento en aplicar los permisos
	return nil
}
