package d75

import "fmt"

const (
	matrixLen   = 102
	numLayers   = 3
	effectCount = 20
	customLED   = 19 // efecto "Customize": colores por tecla
)

// Posiciones de la matriz del firmware (SKLocation.js, entrada "0x06030x0351").
// Corregido con el teclado real: la app oficial pone Fn en 95 y Ctrl der. en 96, pero las
// teclas físicas AltGr y Fn están en 95 y 96 (la 94 no existe).
var matrixKeys = [matrixLen]string{
	"Escape", "", "F1", "F2", "F3", "F4", "F5", "F6", "F7", "F8", "F9", "F10", "F11", "F12",
	"Insert", "", "", "Backquote", "Digit1", "Digit2", "Digit3", "Digit4", "Digit5", "Digit6",
	"Digit7", "Digit8", "Digit9", "Digit0", "Minus", "Equal", "Backspace", "Delete", "", "",
	"Tab", "KeyQ", "KeyW", "KeyE", "KeyR", "KeyT", "KeyY", "KeyU", "KeyI", "KeyO", "KeyP",
	"BracketLeft", "BracketRight", "Backslash", "PageUp", "", "", "CapsLock", "KeyA", "KeyS",
	"KeyD", "KeyF", "KeyG", "KeyH", "KeyJ", "KeyK", "KeyL", "Semicolon", "Quote", "k42",
	"Enter", "PageDown", "", "", "ShiftLeft", "k64", "KeyZ", "KeyX", "KeyC", "KeyV", "KeyB",
	"KeyN", "KeyM", "Comma", "Period", "Slash", "ShiftRight", "ArrowUp", "", "", "",
	"ControlLeft", "MetaLeft", "AltLeft", "lspace", "mspace", "", "Space", "", "mspace",
	"", "AltRight", "Fn", "ArrowLeft", "ArrowDown", "ArrowRight", "", "",
}

var matrixDefaults = [matrixLen]byte{
	41, 0, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 73, 0, 0, 53, 30, 31, 32, 33, 34,
	35, 36, 37, 38, 39, 45, 46, 42, 76, 0, 0, 43, 20, 26, 8, 21, 23, 28, 24, 12, 18, 19, 47,
	48, 49, 75, 0, 0, 57, 4, 22, 7, 9, 10, 11, 13, 14, 15, 51, 52, 49, 40, 78, 0, 0, 225, 0,
	29, 27, 6, 25, 5, 17, 16, 54, 55, 56, 229, 82, 0, 0, 0, 224, 227, 226, 44, 230, 44, 44,
	0, 44, 230, 230, 254, 80, 81, 79, 0, 0,
}

func keyIndex(name string) int {
	for i, k := range matrixKeys {
		if k == name {
			return i
		}
	}
	return -1
}

// Key es una tecla física para dibujar el teclado (posición y tamaño en unidades de tecla).
type Key struct {
	Name    string  `json:"name"`
	Index   int     `json:"index"`
	Label   string  `json:"label"`
	Default int     `json:"default"`
	X       float64 `json:"x"`
	Y       float64 `json:"y"`
	W       float64 `json:"w"`
}

// Distribución 75% ANSI del D75; etiquetas según la distribución española.
var layout = func() []Key {
	type k struct {
		name, label string
		x, w        float64
	}
	rows := []struct {
		y    float64
		keys []k
	}{
		{0, []k{{"Escape", "Esc", 0, 1}, {"F1", "F1", 1.25, 1}, {"F2", "F2", 2.25, 1}, {"F3", "F3", 3.25, 1},
			{"F4", "F4", 4.25, 1}, {"F5", "F5", 5.5, 1}, {"F6", "F6", 6.5, 1}, {"F7", "F7", 7.5, 1},
			{"F8", "F8", 8.5, 1}, {"F9", "F9", 9.75, 1}, {"F10", "F10", 10.75, 1}, {"F11", "F11", 11.75, 1},
			{"F12", "F12", 12.75, 1}, {"Insert", "Insert", 15.5, 1}}},
		{1.25, []k{{"Backquote", "º", 0, 1}, {"Digit1", "1", 1, 1}, {"Digit2", "2", 2, 1}, {"Digit3", "3", 3, 1},
			{"Digit4", "4", 4, 1}, {"Digit5", "5", 5, 1}, {"Digit6", "6", 6, 1}, {"Digit7", "7", 7, 1},
			{"Digit8", "8", 8, 1}, {"Digit9", "9", 9, 1}, {"Digit0", "0", 10, 1}, {"Minus", "'", 11, 1},
			{"Equal", "¡", 12, 1}, {"Backspace", "⌫", 13, 2}, {"Delete", "Supr", 15.5, 1}}},
		{2.25, []k{{"Tab", "Tab", 0, 1.5}, {"KeyQ", "Q", 1.5, 1}, {"KeyW", "W", 2.5, 1}, {"KeyE", "E", 3.5, 1},
			{"KeyR", "R", 4.5, 1}, {"KeyT", "T", 5.5, 1}, {"KeyY", "Y", 6.5, 1}, {"KeyU", "U", 7.5, 1},
			{"KeyI", "I", 8.5, 1}, {"KeyO", "O", 9.5, 1}, {"KeyP", "P", 10.5, 1}, {"BracketLeft", "`", 11.5, 1},
			{"BracketRight", "+", 12.5, 1}, {"Backslash", "ç", 13.5, 1.5}, {"PageUp", "RePág", 15.5, 1}}},
		{3.25, []k{{"CapsLock", "Bloq Mayús", 0, 1.75}, {"KeyA", "A", 1.75, 1}, {"KeyS", "S", 2.75, 1},
			{"KeyD", "D", 3.75, 1}, {"KeyF", "F", 4.75, 1}, {"KeyG", "G", 5.75, 1}, {"KeyH", "H", 6.75, 1},
			{"KeyJ", "J", 7.75, 1}, {"KeyK", "K", 8.75, 1}, {"KeyL", "L", 9.75, 1}, {"Semicolon", "Ñ", 10.75, 1},
			{"Quote", "´", 11.75, 1}, {"Enter", "Intro", 12.75, 2.25}, {"PageDown", "AvPág", 15.5, 1}}},
		{4.25, []k{{"ShiftLeft", "Shift", 0, 2.25}, {"KeyZ", "Z", 2.25, 1}, {"KeyX", "X", 3.25, 1},
			{"KeyC", "C", 4.25, 1}, {"KeyV", "V", 5.25, 1}, {"KeyB", "B", 6.25, 1}, {"KeyN", "N", 7.25, 1},
			{"KeyM", "M", 8.25, 1}, {"Comma", ",", 9.25, 1}, {"Period", ".", 10.25, 1}, {"Slash", "-", 11.25, 1},
			{"ShiftRight", "Shift", 12.25, 1.75}}},
		{5.25, []k{{"ControlLeft", "Ctrl", 0, 1.25}, {"MetaLeft", "Win", 1.25, 1.25}, {"AltLeft", "Alt", 2.5, 1.25},
			{"Space", "Espacio", 3.75, 6.25}, {"AltRight", "AltGr", 10, 1.25}, {"Fn", "Fn", 11.25, 1.25}}},
	}
	var out []Key
	for _, r := range rows {
		for _, c := range r.keys {
			i := keyIndex(c.name)
			out = append(out, Key{c.name, i, c.label, int(matrixDefaults[i]), c.x, r.y, c.w})
		}
	}
	// Flechas desplazadas un cuarto de tecla hacia abajo, como en el teclado real
	for _, a := range []struct {
		name, label string
		x, y        float64
	}{{"ArrowUp", "↑", 14.5, 4.5}, {"ArrowLeft", "←", 13.5, 5.5}, {"ArrowDown", "↓", 14.5, 5.5}, {"ArrowRight", "→", 15.5, 5.5}} {
		i := keyIndex(a.name)
		out = append(out, Key{a.name, i, a.label, int(matrixDefaults[i]), a.x, a.y, 1})
	}
	return out
}()

// Code es una función asignable a una tecla. Group es un identificador que traduce la interfaz;
// Label es la etiqueta por defecto (la interfaz traduce las que son palabras).
type Code struct {
	Code  int    `json:"code"`
	Label string `json:"label"`
	Group string `json:"group"`
}

var codes = func() []Code {
	var out []Code
	add := func(group string, pairs ...any) {
		for i := 0; i < len(pairs); i += 2 {
			out = append(out, Code{pairs[i].(int), pairs[i+1].(string), group})
		}
	}
	add("special", 0x00, "Desactivar", 0xFE, "Fn", 0xFF, "Fn2")
	for i := 0; i < 26; i++ {
		add("letters", 4+i, string(rune('A'+i)))
	}
	add("letters", 0x33, "Ñ")
	for i := 1; i <= 10; i++ {
		add("numbers", 0x1D+i, fmt.Sprint(i%10))
	}
	add("symbols", 0x35, "º", 0x2D, "'", 0x2E, "¡", 0x2F, "`", 0x30, "+", 0x31, "ç", 0x34, "´",
		0x36, ",", 0x37, ".", 0x38, "-", 0x64, "<")
	for i := 0; i < 12; i++ {
		add("function", 0x3A+i, fmt.Sprintf("F%d", i+1))
	}
	for i := 0; i < 12; i++ {
		add("function", 0x68+i, fmt.Sprintf("F%d", i+13))
	}
	add("editing", 0x29, "Esc", 0x28, "Intro", 0x2A, "Retroceso", 0x2B, "Tab", 0x2C, "Espacio",
		0x39, "Bloq Mayús", 0x49, "Insert", 0x4C, "Supr", 0x46, "Impr Pant", 0x47, "Bloq Despl",
		0x48, "Pausa", 0x65, "Menú")
	add("navigation", 0x52, "↑", 0x51, "↓", 0x50, "←", 0x4F, "→", 0x4A, "Inicio", 0x4D, "Fin",
		0x4B, "RePág", 0x4E, "AvPág")
	add("modifiers", 0xE0, "Ctrl", 0xE1, "Shift", 0xE2, "Alt", 0xE3, "Win", 0xE4, "Ctrl der",
		0xE5, "Shift der", 0xE6, "AltGr", 0xE7, "Win der")
	add("media", 0xA1, "Play/Pausa", 0xA5, "Stop", 0xA6, "Anterior", 0xA7, "Siguiente",
		0xA2, "Silencio", 0xA3, "Vol +", 0xA4, "Vol −", 0xA0, "Reproductor")
	add("apps", 0xA8, "Nav. inicio", 0xAB, "Nav. atrás", 0xAC, "Nav. adelante",
		0xA9, "Nav. recargar", 0xAD, "Favoritos", 0xAF, "Mi PC", 0xB1, "Email")
	add("mouse", 0xB4, "Clic izq", 0xB5, "Clic der", 0xB6, "Clic central", 0xB7, "Atrás", 0xB8, "Adelante")
	add("numpad", 0x53, "Bloq Num", 0x54, "Num /", 0x55, "Num *", 0x56, "Num −", 0x57, "Num +",
		0x58, "Num Intro", 0x63, "Num .")
	for i := 1; i <= 10; i++ {
		add("numpad", 0x58+i, fmt.Sprintf("Num %d", i%10))
	}
	return out
}()

var effects = [effectCount]string{
	"Neon Stream", "Ripples Out", "Shadow Out", "Respiratory", "Radiant Splendour", "Drizzle",
	"Turns Around", "Rainbow Sunset", "Breathe", "Floating Light", "Ripple Band", "Little Ripple",
	"Rainbow Wheel", "Superb Wheel", "Simplified Scanning", "Continuous Flow", "Constant On",
	"Glittering", "Ripple Back", "Personalizado",
}

// Paleta fija del firmware para el índice de color; 7 = multicolor.
var palette = []string{"#ff0000", "#00ff00", "#0000ff", "#ffff00", "#ff00ff", "#00ffff", "#ffffff"}

// buildKeymap genera el report 7 cmd 3: las tres capas completas, 102 teclas x [flags, código].
func buildKeymap(p *Profile) ([]byte, error) {
	buf := make([]byte, 620)
	buf[0], buf[1], buf[2] = 7, 3, 1
	for l := 0; l < numLayers; l++ {
		base := 8 + l*matrixLen*2
		for i := 0; i < matrixLen; i++ {
			buf[base+i*2+1] = matrixDefaults[i]
		}
		for name, a := range p.Layers[l] {
			i := keyIndex(name)
			if i < 0 {
				return nil, newError(ErrUnknownKey, name, nil)
			}
			if a.Code < 0 || a.Code > 0xFF || a.Mods < 0 || a.Mods > 0x0F {
				return nil, newError(ErrInvalidAssign, fmt.Sprintf("%s %+v", name, a), nil)
			}
			if n, ok := macroIndex(a.Code); ok && (n >= len(p.Macros) || a.Mods != 0) {
				return nil, newError(ErrInvalidAssign, fmt.Sprintf("%s: macro %d", name, n+1), nil)
			}
			buf[base+i*2] = 0x80 | byte(a.Mods)
			buf[base+i*2+1] = byte(a.Code)
		}
	}
	return buf, nil
}

// buildLED genera el report 7 cmd 2 y, para el efecto personalizado, el cmd 7 con los colores.
func buildLED(l LED) ([][]byte, error) {
	if l.Effect < 0 || l.Effect >= effectCount || l.Brightness < 0 || l.Brightness > 4 ||
		l.Speed < 1 || l.Speed > 5 || l.Color < 0 || l.Color > 7 {
		return nil, newError(ErrInvalidLED, fmt.Sprintf("%+v", l), nil)
	}
	buf := make([]byte, 620)
	buf[0], buf[1], buf[2] = 7, 2, 1
	buf[9], buf[10], buf[11] = byte(l.Effect), byte(l.Brightness), byte(l.Speed)
	buf[12+l.Effect] = byte(l.Color)
	out := [][]byte{buf}
	if l.Effect == customLED {
		custom := make([]byte, 620)
		custom[0], custom[1] = 7, 7
		for name, rgb := range l.Custom {
			i := keyIndex(name)
			if i < 0 {
				return nil, newError(ErrUnknownKey, name, nil)
			}
			for c := 0; c < 3; c++ {
				custom[8+i*3+c] = byte(rgb[c])
			}
		}
		out = append(out, custom)
	}
	return out, nil
}

// cmd8 genera un report 8 (comandos cortos de 8 bytes).
func cmd8(cmd byte, args ...byte) []byte {
	buf := make([]byte, 8)
	buf[0], buf[1] = 8, cmd
	copy(buf[2:], args)
	return buf
}

func buildAll(p *Profile) ([][]byte, error) {
	keymap, err := buildKeymap(p)
	if err != nil {
		return nil, err
	}
	macros, err := buildMacros(p)
	if err != nil {
		return nil, err
	}
	led, err := buildLED(p.LED)
	if err != nil {
		return nil, err
	}
	if p.Debounce < 1 || p.Debounce > 50 {
		return nil, newError(ErrDebounceRange, fmt.Sprint(p.Debounce), nil)
	}
	reports := append(append([][]byte{keymap}, macros...), led...)
	return append(reports, cmd8(9, byte(p.Debounce))), nil
}
