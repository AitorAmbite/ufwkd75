//go:build !linux && !windows

package d75

import "context"

var errUnsupported = newError(ErrUnsupported, "", nil)

func DeviceStatus() Status                                { return Status{Detail: errUnsupported.Error()} }
func sendReports(...[]byte) error                         { return errUnsupported }
func ReadKeyEvents(ctx context.Context, _ func(KeyEvent)) { <-ctx.Done() }
func InstallPermissions() error                           { return errUnsupported }
