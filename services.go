package main

import (
	"context"
	"sync"
	"time"

	"d75studio/internal/d75"

	"github.com/wailsapp/wails/v3/pkg/application"
)

const (
	eventKey    = "d75:key"
	eventStatus = "d75:status"
)

func init() {
	application.RegisterEvent[d75.KeyEvent](eventKey)
	application.RegisterEvent[d75.Status](eventStatus)
}

// KeyboardService habla con el teclado.
type KeyboardService struct {
	cancel context.CancelFunc
}

func (s *KeyboardService) ServiceStartup(ctx context.Context, _ application.ServiceOptions) error {
	ctx, s.cancel = context.WithCancel(ctx)
	app := application.Get()
	go d75.ReadKeyEvents(ctx, func(ev d75.KeyEvent) { app.Event.Emit(eventKey, ev) })
	go watchStatus(ctx, app)
	return nil
}

func (s *KeyboardService) ServiceShutdown() error {
	if s.cancel != nil {
		s.cancel()
	}
	return nil
}

// watchStatus avisa a la interfaz cuando el teclado se conecta o desconecta.
func watchStatus(ctx context.Context, app *application.App) {
	var last d75.Status
	for first := true; ctx.Err() == nil; first = false {
		st := d75.DeviceStatus()
		if first || st != last {
			app.Event.Emit(eventStatus, st)
			last = st
		}
		select {
		case <-ctx.Done():
		case <-time.After(2 * time.Second):
		}
	}
}

func (s *KeyboardService) Meta() d75.Meta            { return d75.GetMeta() }
func (s *KeyboardService) Status() d75.Status        { return d75.DeviceStatus() }
func (s *KeyboardService) Applied() d75.Profile      { return d75.Applied() }
func (s *KeyboardService) InstallPermissions() error { return d75.InstallPermissions() }

var applyMu sync.Mutex

// Apply envía la configuración completa al teclado.
func (s *KeyboardService) Apply(p d75.Profile) error {
	applyMu.Lock()
	defer applyMu.Unlock()
	return d75.Apply(p)
}

// FactoryReset restaura el teclado de fábrica.
func (s *KeyboardService) FactoryReset() error {
	applyMu.Lock()
	defer applyMu.Unlock()
	return d75.FactoryReset()
}

// ProfileService gestiona los perfiles guardados en disco.
type ProfileService struct{}

func (ProfileService) List() ([]string, error)               { return d75.ListProfiles() }
func (ProfileService) Load(name string) (d75.Profile, error) { return d75.LoadProfile(name) }
func (ProfileService) Save(name string, p d75.Profile) error { return d75.SaveProfile(name, p) }
func (ProfileService) Delete(name string) error              { return d75.DeleteProfile(name) }

const profileFilter = "*.json"

// Export pide dónde guardar y escribe el perfil; devuelve false si el usuario cancela.
func (ProfileService) Export(title, filename string, p d75.Profile) (bool, error) {
	path, err := application.Get().Dialog.SaveFile().
		SetMessage(title).
		SetFilename(filename + ".json").
		AddFilter("D75 Studio (JSON)", profileFilter).
		CanCreateDirectories(true).
		PromptForSingleSelection()
	if err != nil || path == "" {
		return false, err
	}
	return true, d75.ExportProfile(path, p)
}

// Imported es un perfil leído de un archivo y guardado con Name. Name vacío: el usuario canceló.
type Imported struct {
	Name    string      `json:"name"`
	Profile d75.Profile `json:"profile"`
}

// Import pide un archivo, lo valida y lo guarda como perfil nuevo.
func (ProfileService) Import(title string) (Imported, error) {
	path, err := application.Get().Dialog.OpenFile().
		SetTitle(title).
		AddFilter("D75 Studio (JSON)", profileFilter).
		CanChooseFiles(true).
		PromptForSingleSelection()
	if err != nil || path == "" {
		return Imported{}, err
	}
	name, p, err := d75.ImportProfile(path)
	return Imported{name, p}, err
}
