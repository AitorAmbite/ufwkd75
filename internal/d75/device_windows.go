//go:build windows

package d75

import (
	"context"
	"fmt"
	"sync"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

// En Windows cada colección de nivel superior de la interfaz es un dispositivo distinto,
// así que cada report va a la colección de su usage page (como DeviceId y
// DeviceId_Set8Bytes en la app oficial).
const (
	pageEvents = 0xFF01 // report 4 (input): pulsaciones físicas
	pageConfig = 0xFF02 // report 7 (feature, 620 bytes)
	pageShort  = 0xFF03 // report 8 (feature, 8 bytes)
)

var reportPage = map[byte]uint16{7: pageConfig, 8: pageShort}

var (
	hidDLL               = windows.NewLazySystemDLL("hid.dll")
	procGetAttributes    = hidDLL.NewProc("HidD_GetAttributes")
	procGetPreparsedData = hidDLL.NewProc("HidD_GetPreparsedData")
	procFreePreparsed    = hidDLL.NewProc("HidD_FreePreparsedData")
	procGetCaps          = hidDLL.NewProc("HidP_GetCaps")
	procSetFeature       = hidDLL.NewProc("HidD_SetFeature")

	hidGUID = windows.GUID{Data1: 0x4D1E55B2, Data2: 0xF16F, Data3: 0x11CF,
		Data4: [8]byte{0x88, 0xCB, 0x00, 0x11, 0x11, 0x00, 0x00, 0x30}}

	deviceMu sync.Mutex
)

type hidAttributes struct {
	Size          uint32
	VendorID      uint16
	ProductID     uint16
	VersionNumber uint16
}

type hidpCaps struct {
	Usage                     uint16
	UsagePage                 uint16
	InputReportByteLength     uint16
	OutputReportByteLength    uint16
	FeatureReportByteLength   uint16
	Reserved                  [17]uint16
	NumberLinkCollectionNodes uint16
	NumberInputButtonCaps     uint16
	NumberInputValueCaps      uint16
	NumberInputDataIndices    uint16
	NumberOutputButtonCaps    uint16
	NumberOutputValueCaps     uint16
	NumberOutputDataIndices   uint16
	NumberFeatureButtonCaps   uint16
	NumberFeatureValueCaps    uint16
	NumberFeatureDataIndices  uint16
}

type collection struct {
	path string
	caps hidpCaps
}

func openPath(path string, access uint32) (windows.Handle, error) {
	p, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return 0, err
	}
	return windows.CreateFile(p, access, windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE,
		nil, windows.OPEN_EXISTING, 0, 0)
}

// findCollections devuelve las colecciones HID del D75 indexadas por usage page.
func findCollections() (map[uint16]collection, error) {
	paths, err := windows.CM_Get_Device_Interface_List("", &hidGUID, windows.CM_GET_DEVICE_INTERFACE_LIST_PRESENT)
	if err != nil {
		return nil, err
	}
	found := map[uint16]collection{}
	for _, path := range paths {
		// Acceso 0: basta para leer atributos y capacidades, incluso de colecciones
		// que Windows tiene abiertas en exclusiva (el teclado estándar)
		h, err := openPath(path, 0)
		if err != nil {
			continue
		}
		attr := hidAttributes{Size: uint32(unsafe.Sizeof(hidAttributes{}))}
		ok, _, _ := procGetAttributes.Call(uintptr(h), uintptr(unsafe.Pointer(&attr)))
		if ok == 0 || attr.VendorID != vendorID || !isD75(attr.ProductID) {
			windows.CloseHandle(h)
			continue
		}
		var preparsed uintptr
		if ok, _, _ := procGetPreparsedData.Call(uintptr(h), uintptr(unsafe.Pointer(&preparsed))); ok != 0 {
			var caps hidpCaps
			const hidpStatusSuccess = 0x00110000
			if st, _, _ := procGetCaps.Call(preparsed, uintptr(unsafe.Pointer(&caps))); st == hidpStatusSuccess {
				found[caps.UsagePage] = collection{path, caps}
			}
			procFreePreparsed.Call(preparsed)
		}
		windows.CloseHandle(h)
	}
	if _, ok := found[pageConfig]; !ok {
		return nil, errNotFound
	}
	return found, nil
}

func isD75(pid uint16) bool {
	for _, p := range productIDs {
		if int(pid) == p {
			return true
		}
	}
	return false
}

func DeviceStatus() Status {
	cols, err := findCollections()
	if err != nil {
		return Status{Detail: err.Error()}
	}
	h, err := openPath(cols[pageConfig].path, windows.GENERIC_READ|windows.GENERIC_WRITE)
	if err != nil {
		return Status{Detail: newError(ErrOpen, "", err).Error()}
	}
	windows.CloseHandle(h)
	return Status{Connected: true, Detail: "USB HID"}
}

func sendReports(reports ...[]byte) error {
	deviceMu.Lock()
	defer deviceMu.Unlock()
	cols, err := findCollections()
	if err != nil {
		return err
	}
	handles := map[uint16]windows.Handle{}
	defer func() {
		for _, h := range handles {
			windows.CloseHandle(h)
		}
	}()
	for _, report := range reports {
		page, ok := reportPage[report[0]]
		col, found := cols[page]
		if !ok || !found {
			return newError(ErrNoCollection, fmt.Sprint(report[0]), nil)
		}
		h, open := handles[page]
		if !open {
			if h, err = openPath(col.path, windows.GENERIC_READ|windows.GENERIC_WRITE); err != nil {
				return newError(ErrOpen, "", err)
			}
			handles[page] = h
		}
		// HidD_SetFeature exige exactamente la longitud declarada por el dispositivo
		buf := make([]byte, col.caps.FeatureReportByteLength)
		copy(buf, report)
		if ok, _, callErr := procSetFeature.Call(uintptr(h), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf))); ok == 0 {
			return newError(ErrSend, fmt.Sprintf("report %d cmd %d", report[0], report[1]), callErr)
		}
		time.Sleep(reportDelay)
	}
	return nil
}

// ReadKeyEvents lee pulsaciones físicas hasta que se cancele ctx; reconecta solo.
func ReadKeyEvents(ctx context.Context, fn func(KeyEvent)) {
	for ctx.Err() == nil {
		cols, err := findCollections()
		col, ok := cols[pageEvents]
		if err != nil || !ok {
			sleepCtx(ctx, time.Second)
			continue
		}
		h, err := openPath(col.path, windows.GENERIC_READ)
		if err != nil {
			sleepCtx(ctx, time.Second)
			continue
		}
		stop := context.AfterFunc(ctx, func() { windows.CancelIoEx(h, nil) })
		buf := make([]byte, max(int(col.caps.InputReportByteLength), 64))
		for {
			var n uint32
			if err := windows.ReadFile(h, buf, &n, nil); err != nil {
				break
			}
			if ev, ok := parseKeyEvent(buf[:n]); ok {
				fn(ev)
			}
		}
		stop()
		windows.CloseHandle(h)
		sleepCtx(ctx, time.Second)
	}
}

// InstallPermissions no hace falta en Windows: el acceso HID no requiere permisos.
func InstallPermissions() error {
	return newError(ErrNoPermsNeeded, "", nil)
}
