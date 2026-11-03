package d75

// Error lleva un código estable que la interfaz traduce al idioma del usuario,
// más un detalle técnico opcional (ruta, tecla, error del sistema…) que se muestra tal cual.
// Error() produce "código" o "código: detalle".
type Error struct {
	Code   string
	Detail string
	Err    error
}

const (
	ErrNotFound      = "not_found"
	ErrOpen          = "open_failed"
	ErrSend          = "send_failed"
	ErrPermission    = "permission_failed"
	ErrNoPermsNeeded = "no_permissions_needed"
	ErrUnsupported   = "unsupported_os"
	ErrInvalidName   = "invalid_name"
	ErrUnknownKey    = "unknown_key"
	ErrInvalidAssign = "invalid_assignment"
	ErrInvalidLED    = "invalid_led"
	ErrDebounceRange = "debounce_range"
	ErrNoCollection  = "no_collection"
)

func (e *Error) Error() string {
	s := e.Code
	detail := e.Detail
	if e.Err != nil {
		if detail != "" {
			detail += ": "
		}
		detail += e.Err.Error()
	}
	if detail != "" {
		s += ": " + detail
	}
	return s
}

func (e *Error) Unwrap() error { return e.Err }

func newError(code, detail string, err error) error {
	return &Error{Code: code, Detail: detail, Err: err}
}
