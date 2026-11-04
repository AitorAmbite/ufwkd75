package d75

import (
	"bytes"
	"testing"
)

// Paquete probado en el teclado real con d75.py (ArrowUp = B): la flecha arriba escribió "b".
func TestKeymapMatchesHardwareTest(t *testing.T) {
	p := defaultProfile()
	p.Layers[0]["ArrowUp"] = Assign{Code: 0x05}
	buf, err := buildKeymap(&p)
	if err != nil {
		t.Fatal(err)
	}
	if len(buf) != 620 || !bytes.Equal(buf[:3], []byte{7, 3, 1}) {
		t.Fatalf("cabecera incorrecta: % x", buf[:8])
	}
	if got := buf[8+81*2 : 8+81*2+2]; !bytes.Equal(got, []byte{0x80, 0x05}) {
		t.Fatalf("ArrowUp = % x", got)
	}
	// Las demás teclas, en las tres capas, llevan flags 0 y el código de fábrica
	for l := 0; l < numLayers; l++ {
		for i := 0; i < matrixLen; i++ {
			if l == 0 && i == 81 {
				continue
			}
			off := 8 + (l*matrixLen+i)*2
			if buf[off] != 0 || buf[off+1] != matrixDefaults[i] {
				t.Fatalf("capa %d tecla %d = % x", l, i, buf[off:off+2])
			}
		}
	}
}

func TestLayoutKeysExistInMatrix(t *testing.T) {
	if len(layout) != 80 {
		t.Fatalf("el D75 tiene 80 teclas, layout tiene %d", len(layout))
	}
	for _, k := range layout {
		if keyIndex(k.Name) < 0 {
			t.Errorf("tecla %q no está en la matriz", k.Name)
		}
	}
}

func TestLEDCustom(t *testing.T) {
	l := LED{Effect: customLED, Brightness: 4, Speed: 3, Color: 0, Custom: map[string][3]int{"KeyA": {1, 2, 3}}}
	out, err := buildLED(l)
	if err != nil || len(out) != 2 {
		t.Fatalf("esperaba 2 reports: %v", err)
	}
	if out[0][9] != 19 || out[0][10] != 4 || out[0][11] != 3 {
		t.Fatalf("efecto: % x", out[0][:12])
	}
	i := keyIndex("KeyA")
	if !bytes.Equal(out[1][8+i*3:8+i*3+3], []byte{1, 2, 3}) {
		t.Fatal("color de KeyA mal colocado")
	}
}

func TestRejectsBadInput(t *testing.T) {
	p := defaultProfile()
	p.Layers[1]["NoExiste"] = Assign{Code: 4}
	if _, err := buildKeymap(&p); err == nil {
		t.Fatal("debería rechazar teclas desconocidas")
	}
	if _, err := buildLED(LED{Effect: 20, Speed: 1}); err == nil {
		t.Fatal("debería rechazar efecto 20")
	}
}

// Posiciones medidas en el teclado real con el report 4 ([04 F8 posición pulsada]).
func TestMatrixMatchesHardware(t *testing.T) {
	measured := map[string]int{"MetaLeft": 86, "AltLeft": 87, "AltRight": 95, "Fn": 96, "KeyQ": 35, "ArrowUp": 81, "ControlLeft": 85}
	for name, pos := range measured {
		if got := keyIndex(name); got != pos {
			t.Errorf("%s: posición %d, medida en el teclado %d", name, got, pos)
		}
	}
	if matrixDefaults[95] != 0xE6 || matrixDefaults[96] != 0xFE {
		t.Errorf("AltGr/Fn por defecto: %#x %#x", matrixDefaults[95], matrixDefaults[96])
	}
}

// Formato de MacroToData en JupengSeries.js: repeticiones BE y eventos [delay>>8 | 0x80 pulsar, delay, usage].
func TestMacroReport(t *testing.T) {
	p := defaultProfile()
	p.Macros = []Macro{
		{Name: "x", Mode: MacroToggle, Events: []MacroEvent{{Code: 0x04, Down: true, Delay: 5}}},
		{Name: "hola", Mode: MacroCount, Count: 3, Events: []MacroEvent{
			{Code: 0x0B, Down: true, Delay: 300}, {Code: 0x0B, Down: false, Delay: 20},
		}},
	}
	p.Layers[1]["KeyQ"] = Assign{Code: 0xF1}
	reports, err := buildAll(&p)
	if err != nil {
		t.Fatal(err)
	}
	m := reports[2]
	if !bytes.Equal(m[:4], []byte{7, 5, 1, 1}) {
		t.Fatalf("cabecera % x", m[:8])
	}
	if want := []byte{0, 3, 0x81, 0x2C, 0x0B, 0x00, 0x14, 0x0B, 0}; !bytes.Equal(m[8:17], want) {
		t.Fatalf("datos % x", m[8:17])
	}
	// Retardo mínimo de 20 ms y modo "hasta volver a pulsar" = 0x0000
	if want := []byte{0, 0, 0x80, 0x14, 0x04}; !bytes.Equal(reports[1][8:13], want) {
		t.Fatalf("toggle % x", reports[1][8:13])
	}
	if off := 8 + matrixLen*2 + 35*2; !bytes.Equal(reports[0][off:off+2], []byte{0x80, 0xF1}) {
		t.Fatalf("KeyQ capa Fn = % x", reports[0][off:off+2])
	}

	p.Macros[1].Mode = MacroHold
	if buf, _ := buildMacro(1, p.Macros[1]); !bytes.Equal(buf[8:10], []byte{0xFF, 0xFF}) {
		t.Fatalf("hold % x", buf[8:10])
	}
	p.Layers[0]["KeyW"] = Assign{Code: 0xF5}
	if _, err := buildAll(&p); err == nil {
		t.Fatal("debería rechazar una macro que no existe")
	}
	delete(p.Layers[0], "KeyW")
	p.Macros[0].Events = make([]MacroEvent, maxMacroEvents+1)
	if _, err := buildAll(&p); err == nil {
		t.Fatal("debería rechazar macros demasiado largas")
	}
}
