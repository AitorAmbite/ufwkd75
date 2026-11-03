package d75

import (
	"errors"
	"os"
	"runtime"
)

// Meta es todo lo que la interfaz necesita para dibujar el teclado y sus opciones.
type Meta struct {
	Layout       []Key    `json:"layout"`
	Codes        []Code   `json:"codes"`
	Effects      []string `json:"effects"`
	Palette      []string `json:"palette"`
	CustomEffect int      `json:"customEffect"`
	Default      Profile  `json:"default"`
	Platform     string   `json:"platform"`
}

func GetMeta() Meta {
	return Meta{layout, codes, effects[:], palette, customLED, DefaultProfile(), runtime.GOOS}
}

// Status describe la conexión con el teclado.
type Status struct {
	Connected       bool   `json:"connected"`
	NeedsPermission bool   `json:"needsPermission"`
	Detail          string `json:"detail"`
}

// KeyEvent es una pulsación física. El teclado la envía siempre por la interfaz de
// configuración como report 4: [04 F8 posición pulsada], sea cual sea la asignación.
type KeyEvent struct {
	Index   int  `json:"index"`
	Pressed bool `json:"pressed"`
}

func parseKeyEvent(buf []byte) (KeyEvent, bool) {
	if len(buf) >= 4 && buf[0] == 4 && buf[1] == 0xF8 {
		return KeyEvent{Index: int(buf[2]), Pressed: buf[3] != 0}, true
	}
	return KeyEvent{}, false
}

func DefaultProfile() Profile { return defaultProfile() }

// Apply envía la configuración completa al teclado y la recuerda como "aplicada".
func Apply(p Profile) error {
	normalize(&p)
	reports, err := buildAll(&p)
	if err != nil {
		return err
	}
	if err := sendReports(reports...); err != nil {
		return err
	}
	return saveApplied(p)
}

// FactoryReset devuelve el teclado a su configuración original.
func FactoryReset() error {
	if err := sendReports(cmd8(8)); err != nil {
		return err
	}
	return saveApplied(defaultProfile())
}

// Applied devuelve lo último enviado al teclado (el firmware no permite leerlo).
func Applied() Profile {
	if path, err := lastAppliedPath(); err == nil {
		if p, err := readProfile(path); err == nil {
			return p
		}
	}
	return defaultProfile()
}

func saveApplied(p Profile) error {
	path, err := lastAppliedPath()
	if err != nil {
		return err
	}
	return writeProfile(path, p)
}

func ListProfiles() ([]string, error) { return listProfiles() }

func LoadProfile(name string) (Profile, error) {
	path, err := profilePath(name)
	if err != nil {
		return Profile{}, err
	}
	return readProfile(path)
}

func SaveProfile(name string, p Profile) error {
	normalize(&p)
	if _, err := buildAll(&p); err != nil {
		return err
	}
	path, err := profilePath(name)
	if err != nil {
		return err
	}
	return writeProfile(path, p)
}

func DeleteProfile(name string) error {
	path, err := profilePath(name)
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

func normalize(p *Profile) {
	for i := range p.Layers {
		if p.Layers[i] == nil {
			p.Layers[i] = map[string]Assign{}
		}
	}
}
