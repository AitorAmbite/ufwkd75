package d75

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Assign reasigna una tecla: código HID (o especial del firmware) más modificadores
// (bit0 Ctrl, bit1 Shift, bit2 Alt, bit3 Win). Las teclas sin entrada mantienen su valor de fábrica.
type Assign struct {
	Code int `json:"code"`
	Mods int `json:"mods,omitempty"`
}

type LED struct {
	Effect     int               `json:"effect"`
	Brightness int               `json:"brightness"`
	Speed      int               `json:"speed"`
	Color      int               `json:"color"`
	Custom     map[string][3]int `json:"custom,omitempty"`
}

// Profile es la configuración completa del teclado. Se aplica siempre entera,
// porque el firmware no permite leer la configuración actual.
type Profile struct {
	Layers   [numLayers]map[string]Assign `json:"layers"`
	LED      LED                          `json:"led"`
	Debounce int                          `json:"debounce"`
}

// defaultProfile coincide con el estado de fábrica (SupportDevice.db).
func defaultProfile() Profile {
	p := Profile{LED: LED{Effect: 0, Brightness: 4, Speed: 5, Color: 7}, Debounce: 5}
	for i := range p.Layers {
		p.Layers[i] = map[string]Assign{}
	}
	return p
}

var validName = regexp.MustCompile(`^[\p{L}\p{N} _-]{1,64}$`)

func configDir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "d75"), nil
}

func profilePath(name string) (string, error) {
	if !validName.MatchString(name) || strings.TrimSpace(name) == "" {
		return "", newError(ErrInvalidName, name, nil)
	}
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "profiles", name+".json"), nil
}

func listProfiles() ([]string, error) {
	dir, err := configDir()
	if err != nil {
		return nil, err
	}
	files, _ := filepath.Glob(filepath.Join(dir, "profiles", "*.json"))
	names := []string{}
	for _, f := range files {
		names = append(names, strings.TrimSuffix(filepath.Base(f), ".json"))
	}
	sort.Strings(names)
	return names, nil
}

func readProfile(path string) (Profile, error) {
	p := defaultProfile()
	data, err := os.ReadFile(path)
	if err != nil {
		return p, err
	}
	if err := json.Unmarshal(data, &p); err != nil {
		return p, err
	}
	for i := range p.Layers {
		if p.Layers[i] == nil {
			p.Layers[i] = map[string]Assign{}
		}
	}
	return p, nil
}

func writeProfile(path string, p Profile) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// lastAppliedPath guarda lo último enviado al teclado, para mostrarlo al abrir la app.
func lastAppliedPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "aplicado.json"), nil
}
