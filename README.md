# ufwkd75 · D75 Studio

Configurador de escritorio para el teclado **Weikav D75** (USB `0603:0351`, dongle `0603:0352`), para Linux y Windows.
Hecho con [Wails 3](https://v3.wails.io) (Go + Svelte).

- **Teclas**: reasignación en las 3 capas (Normal, Fn, Fn2), combinaciones con modificadores, multimedia, ratón…
- **Macros**: hasta 10 secuencias grabadas tecleando, con retardos editables y repetición (N veces, mientras se mantiene o hasta volver a pulsar).
- **Iluminación**: los 20 efectos del firmware y modo personalizado con RGB tecla a tecla, degradados y arcoíris.
- **Perfiles**: guardados como JSON en `~/.config/d75/profiles` (Linux) o `%AppData%\d75\profiles` (Windows). Se pueden exportar e importar como archivo para compartirlos.
- **Detectar**: pulsa una tecla del teclado real y se selecciona en pantalla.

La configuración se guarda en el propio teclado: una vez aplicada funciona en cualquier ordenador sin software.

## Uso

Descarga el ejecutable desde [Releases](../../releases).

| Sistema | Ejecutable | Notas |
|---|---|---|
| Linux | `d75-studio-linux-amd64` | La primera vez, *Ajustes → Permisos* instala la regla udev (pide contraseña con pkexec) |
| Windows | `d75-studio-windows-amd64.exe` | Necesita WebView2 (viene con Windows 10/11) |

## Desarrollo

```sh
wails3 dev                  # modo desarrollo con recarga
wails3 build                # binario para el sistema actual
wails3 build GOOS=windows   # .exe de Windows desde Linux (sin CGO)
go test ./internal/...      # tests del protocolo
```

## Estructura

- `internal/d75/protocol.go`: matriz del teclado y construcción de los paquetes HID
- `internal/d75/device_linux.go`: hidraw + ioctl `HIDIOCSFEATURE`
- `internal/d75/device_windows.go`: `hid.dll` (`HidD_SetFeature`), una colección por usage page
- `internal/d75/profile.go`: perfiles en disco
- `services.go`: servicios expuestos a la interfaz y eventos (`d75:key`, `d75:status`)
- `frontend/`: interfaz en Svelte 5

El protocolo está documentado en [`docs/PROTOCOLO_D75.md`](docs/PROTOCOLO_D75.md).
