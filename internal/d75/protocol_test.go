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
