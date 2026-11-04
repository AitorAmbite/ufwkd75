package d75

import "fmt"

const (
	maxMacros      = 10
	macroCodeBase  = 0xF0 // en el keymap, 0xF0 + n asigna la macro n
	macroDataLen   = 256  // bytes de datos por macro: 2 de repeticiones + 3 por evento
	maxMacroEvents = (macroDataLen - 2) / 3
	minMacroDelay  = 20
	maxMacroDelay  = 0x7FFF // el bit alto del primer byte indica pulsar/soltar

	// Valores especiales de repeticiones (los de la app oficial)
	repeatHold   = 0xFFFF // se repite mientras la tecla esté pulsada
	repeatToggle = 0x0000 // se repite hasta volver a pulsar la tecla
	maxRepeat    = 0xFFFE
)

// Modos de repetición de una macro.
const (
	MacroCount  = "count"  // Count veces
	MacroHold   = "hold"   // mientras se mantenga pulsada
	MacroToggle = "toggle" // hasta volver a pulsarla
)

// MacroEvent pulsa o suelta una tecla (usage HID) y espera Delay ms antes del siguiente evento.
type MacroEvent struct {
	Code  int  `json:"code"`
	Down  bool `json:"down"`
	Delay int  `json:"delay"`
}

type Macro struct {
	Name   string       `json:"name"`
	Mode   string       `json:"mode"`
	Count  int          `json:"count,omitempty"`
	Events []MacroEvent `json:"events"`
}

// macroIndex indica si un código del keymap asigna una macro y cuál.
func macroIndex(code int) (int, bool) {
	n := code - macroCodeBase
	return n, n >= 0 && n < maxMacros
}

func (m Macro) repeat() (int, error) {
	switch m.Mode {
	case MacroHold:
		return repeatHold, nil
	case MacroToggle:
		return repeatToggle, nil
	case MacroCount, "":
		if m.Count < 1 || m.Count > maxRepeat {
			return 0, newError(ErrInvalidMacro, fmt.Sprintf("%s: count %d", m.Name, m.Count), nil)
		}
		return m.Count, nil
	}
	return 0, newError(ErrInvalidMacro, fmt.Sprintf("%s: mode %q", m.Name, m.Mode), nil)
}

// buildMacro genera el report 7 cmd 5 de la macro n:
// datos desde el byte 8 = [repeticiones u16 BE, eventos de 3 bytes [delay>>8 | 0x80 si pulsa, delay & 0xFF, usage]].
func buildMacro(n int, m Macro) ([]byte, error) {
	if len(m.Events) == 0 || len(m.Events) > maxMacroEvents {
		return nil, newError(ErrInvalidMacro, fmt.Sprintf("%s: %d eventos (1-%d)", m.Name, len(m.Events), maxMacroEvents), nil)
	}
	rep, err := m.repeat()
	if err != nil {
		return nil, err
	}
	buf := make([]byte, 620)
	buf[0], buf[1], buf[2], buf[3] = 7, 5, 1, byte(n)
	buf[8], buf[9] = byte(rep>>8), byte(rep)
	for i, ev := range m.Events {
		if ev.Code < 0x04 || ev.Code > 0xE7 || ev.Delay < 0 || ev.Delay > maxMacroDelay {
			return nil, newError(ErrInvalidMacro, fmt.Sprintf("%s: evento %d %+v", m.Name, i+1, ev), nil)
		}
		delay := max(ev.Delay, minMacroDelay)
		off := 10 + i*3
		buf[off], buf[off+1], buf[off+2] = byte(delay>>8), byte(delay), byte(ev.Code)
		if ev.Down {
			buf[off] |= 0x80
		}
	}
	return buf, nil
}

func buildMacros(p *Profile) ([][]byte, error) {
	if len(p.Macros) > maxMacros {
		return nil, newError(ErrInvalidMacro, fmt.Sprintf("%d macros (máximo %d)", len(p.Macros), maxMacros), nil)
	}
	var out [][]byte
	for n, m := range p.Macros {
		buf, err := buildMacro(n, m)
		if err != nil {
			return nil, err
		}
		out = append(out, buf)
	}
	return out, nil
}
