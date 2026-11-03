package d75

import (
	"context"
	"time"
)

const (
	vendorID = 0x0603
	// Pausa entre reports, la misma que usa la app oficial
	reportDelay = 100 * time.Millisecond
)

var productIDs = []int{0x0351, 0x0352} // cable, dongle 2.4G

var errNotFound = newError(ErrNotFound, "", nil)

func sleepCtx(ctx context.Context, d time.Duration) {
	select {
	case <-ctx.Done():
	case <-time.After(d):
	}
}
