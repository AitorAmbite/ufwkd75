package d75

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestExportImportRoundTrip(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("APPDATA", dir)

	p := defaultProfile()
	p.Layers[1]["KeyQ"] = Assign{Code: 0xF0}
	p.Macros = []Macro{{Name: "hola", Mode: MacroCount, Count: 2, Events: []MacroEvent{{Code: 0x0B, Down: true, Delay: 50}}}}
	p.LED.Effect = 3

	file := filepath.Join(dir, "Mi perfil: juegos")
	if err := ExportProfile(file, p); err != nil {
		t.Fatal(err)
	}
	file += ".json" // sin extensión se añade .json
	// Dos importaciones del mismo archivo no se pisan
	for _, want := range []string{"Mi perfil juegos", "Mi perfil juegos 2"} {
		name, got, err := ImportProfile(file)
		if err != nil {
			t.Fatal(err)
		}
		if name != want || !reflect.DeepEqual(got, p) {
			t.Fatalf("importado %q %+v", name, got)
		}
		if _, err := LoadProfile(name); err != nil {
			t.Fatal("no se guardó como perfil:", err)
		}
	}
}

func TestImportRejectsInvalidFiles(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("APPDATA", dir)
	for name, data := range map[string]string{
		"no-json.json":     "hola",
		"otro-json.json":   `{"version": 1}`,
		"macro-falta.json": `{"layers": [{"KeyQ": {"code": 240}}, {}, {}], "led": {"effect": 0, "brightness": 4, "speed": 5, "color": 7}, "debounce": 5}`,
		"tecla-rara.json":  `{"layers": [{"Nope": {"code": 4}}, {}, {}], "led": {"effect": 0, "brightness": 4, "speed": 5, "color": 7}, "debounce": 5}`,
	} {
		path := filepath.Join(dir, name)
		os.WriteFile(path, []byte(data), 0o644)
		if _, _, err := ImportProfile(path); err == nil {
			t.Errorf("%s: debería fallar", name)
		}
	}
	if names, _ := ListProfiles(); len(names) != 0 {
		t.Fatalf("no debería guardar nada: %v", names)
	}
}
