# Protocolo HID del Weikav D75 (0603:0351)

Obtenido del código de la app oficial "D75 Keyboard" (Electron), desofuscado:
`backend/protocol/device/keyborad/JupengSeries.js`, `backend/protocol/device/Device.js`,
`backend/others/SKLocation.js`, `backend/others/SupportData.js` y la base de datos
`userappdata/D75_3M_Series/data/SupportDevice.db`. Los tamaños se han comprobado con el
report descriptor real del teclado.

- Extracción del asar: `~/d75-asar/`
- Firmware "JP" = Jupeng (la clase del driver se llama `JupengSeries`)
- VID/PID: `0603:0351` por cable, `0603:0352` con dongle 2.4G

## Canales

Todo el canal de configuración está en la **interfaz 1** (hidraw8). La interfaz 2 (hidraw9,
usage page 0x8C) no la usa la app; probablemente es para actualizar el firmware.

| Report ID | Usage page / usage | Tipo | Tamaño (con ID) | Uso en la app |
|---|---|---|---|---|
| 2 | Keyboard | Input | 16 | NKRO |
| 3 | Consumer | Input | 3 | Multimedia |
| 4 | 0xFF01 / 0x01 | Input + Output | 9 | Eventos de tecla `[04 F8 idx pulsada]` (modo test) |
| 5 | Generic Desktop / Mouse | Input | 6 | Ratón (teclas con función de ratón) |
| 6 | 0xFF55 / 0x0202 | Input + Output | 64 | Streaming de LEDs (modo sync) |
| **7** | **0xFF02 / 0x01** | **Feature** | **620** | **Keymap, macros, efectos LED, colores por tecla** |
| **8** | **0xFF03 / 0x01** | **Feature** | **8** | **Comandos cortos: reset, debounce, apagado** |

En Linux: `ioctl(fd, HIDIOCSFEATURE(620), buf)` con `buf[0] = 7`, o `HIDIOCSFEATURE(8)`
con `buf[0] = 8`. La app nunca lee la configuración del teclado; solo escribe. No se puede
hacer copia de seguridad del estado actual: la red de seguridad es el **reset de fábrica**
(report 8, comando 8).

## Report 7 (620 bytes): `[07, cmd, ...]`, datos a partir del byte 8

### cmd 3: keymap (las 3 capas a la vez)

```
[0]=07 [1]=03 [2]=fnModeindex+1 [3..7]=00
[8 ..   211]  capa 0 (normal)   102 teclas x 2 bytes
[212 .. 415]  capa 1 (Fn)       102 x 2
[416 .. 619]  capa 2 (Fn 2)     102 x 2
```

Cada tecla tiene 2 bytes, `[flags, código]`, y se indexa por **posición en la matriz**
(ver la tabla de abajo):

| flags | código | Significado |
|---|---|---|
| `0x00` | código por defecto | Sin cambios (la app pone el valor de `KeyMatrixNormal_Default`) |
| `0x80` | usage HID | Tecla reasignada (código HID estándar o especial, ver abajo) |
| `0x80 \| mods` | usage HID | Combinación: mods = bit0 Ctrl, bit1 Shift, bit2 Alt, bit3 Win |
| `0x80` | `0xF0 + n` | Macro número n (0-9) |

Envía siempre las tres capas: lo que pongas a cero en una capa se queda a cero.

### cmd 5: macro

```
[0]=07 [1]=05 [2]=01 [3]=n (0-9)
[8..9]   repeticiones, uint16 big-endian
[10..]   eventos de 3 bytes: [ (delay>>8) | 0x80 si pulsar, delay & 0xFF, usage HID ]
         retardo mínimo 20 ms, máximo 256 bytes de datos (unos 84 eventos)
```

### cmd 2: efecto LED

```
[0]=07 [1]=02 [2]=01
[9]=efecto (0-19)  [10]=brillo (0-4)  [11]=velocidad (1-5)
[12+efecto]=índice de color (7 = multicolor)
```

Efectos: 0 Neon Stream, 1 Ripples Out, 2 Shadow Out, 3 Respiratory, 4 Radiant Splendour,
5 Drizzle, 6 Turns Around, 7 Rainbow Sunset, 8 Breathe, 9 Floating Light, 10 RippleBand,
11 Little Ripple, 12 Rainbow Wheel, 13 Superb Wheel, 14 Simplified Scanning,
15 Continuous Flow, 16 Constant On, 17 Glittering, 18 Ripple Back, 19 Customize.

### cmd 7: colores por tecla (para el efecto 19, Customize)

`[0]=07 [1]=07`, y luego `[8 + idx*3 + {0,1,2}] = R, G, B`, con idx = posición en la matriz.

## Report 8 (8 bytes): `[08, cmd, args...]`

| cmd | Argumentos | Función |
|---|---|---|
| 6 | — | Entra en modo streaming de LEDs (después se usa el report 6) |
| 8 | — | **Reset de fábrica** |
| 9 | `[2]` = ms | Debounce (5 por defecto) |
| 11 | `[2]` = lightOff+1, `[3]` = hibernate+1 | Apagado/hibernación (solo inalámbrico; índices en `[1,2,5,10,15,20,25,30]` min y `[3,5,10,15,20,25,30,60]` min) |

## Report 6 (64 bytes, output): streaming de LEDs

`[06, página 1-5, 62 bytes RGB]`, 5 páginas que cubren `idx*3` sobre 306 bytes.
Hay que enviar antes `report 8 cmd 6`.

## Matriz del D75 (102 posiciones, índice = posición en el report)

```
  0 Esc          1 -            2 F1           3 F2           4 F3           5 F4
  6 F5           7 F6           8 F7           9 F8          10 F9          11 F10
 12 F11         13 F12         14 Insert      15 -           16 -           17 `
 18 1           19 2           20 3           21 4           22 5           23 6
 24 7           25 8           26 9           27 0           28 -           29 =
 30 Backspace   31 Delete      32 -           33 -           34 Tab         35 Q
 36 W           37 E           38 R           39 T           40 Y           41 U
 42 I           43 O           44 P           45 [           46 ]           47 \
 48 PageUp      49 -           50 -           51 CapsLock    52 A           53 S
 54 D           55 F           56 G           57 H           58 J           59 K
 60 L           61 ;           62 '           63 ISO#(k42)   64 Enter       65 PageDown
 66 -           67 -           68 LShift      69 ISO\(k64)   70 Z           71 X
 72 C           73 V           74 B           75 N           76 M           77 ,
 78 .           79 /           80 RShift      81 Up          82 -           83 -
 84 -           85 LCtrl       86 LWin        87 LAlt        88 lspace      89 mspace
 90 -           91 Space       92 -           93 mspace      94 (no existe) 95 AltGr
 96 Fn          97 Left        98 Down        99 Right      100 -          101 -
```

Valores por defecto (capa 0), en usage HID:
`41,0,58..69,73,0,0,53,30..39,45,46,42,76,0,0,43,20,26,8,21,23,28,24,12,18,19,47,48,49,75,0,0,
57,4,22,7,9,10,11,13,14,15,51,52,49,40,78,0,0,225,0,29,27,6,25,5,17,16,54,55,56,229,82,0,0,0,
224,227,226,44,230,44,44,0,44,230,230,254,80,81,79,0,0`
(Fn = 254 = 0xFE)

**Corrección medida en el teclado real:** la app oficial pone Fn en la posición 95 y
Ctrl derecho en la 96, pero las teclas físicas AltGr y Fn están en 95 y 96. Con los datos
de la app, AltGr hace de Fn y Fn hace de Ctrl derecho. Los valores de arriba ya están corregidos.

## Report 4 (input): posición física de cada pulsación

El teclado envía siempre `[04 F8 posición pulsada(1/0)]` por la interfaz 1, sea cual sea
la asignación de la tecla. Sirve para identificar teclas y verificar la matriz.

## Códigos especiales (no estándar HID, propios del firmware)

| Código | Función | Código | Función |
|---|---|---|---|
| 0x00 | Desactivar | 0xA8 | Navegador: inicio |
| 0xA0 | Reproductor | 0xA9 | Navegador: recargar |
| 0xA1 | Play/Pausa | 0xAB | Navegador: atrás |
| 0xA2 | Silencio | 0xAC | Navegador: adelante |
| 0xA3 | Volumen + | 0xAD | Favoritos |
| 0xA4 | Volumen − | 0xAF | Mi PC |
| 0xA5 | Stop | 0xB1 | Email |
| 0xA6 | Anterior | 0xB2 / 0xB3 | Perfil + / − |
| 0xA7 | Siguiente | 0xB4-0xB8 | Ratón: izq, der, central, atrás, adelante |
| 0xFE | Fn | 0xFF | Fn2 |

## Pendiente

- Por confirmar en el hardware: en las capas Fn, con flags `0x00`, ¿el firmware usa su
  valor por defecto o el código que se envía? La app manda ahí el valor por defecto de la
  capa normal.
- `GetFWVersion` lo resuelve el módulo nativo `hiddevice.node`, probablemente leyendo
  bcdDevice.
