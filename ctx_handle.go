package olog

import (
	"context"
	"sync/atomic"
	"unsafe"
)

// Pointer to the default context handle
var dch unsafe.Pointer

// Initializes the default context handle with an empty handle
func init() {
	SetDefCtxHandle(emptyHandle)
}

// CtxHandle type that returns a slice of Fields based on the provided context
type CtxHandle func(ctx context.Context) []Field

// ctxHandler hold a context handle function
type ctxHandler struct {
	handle CtxHandle
}

// SetDefCtxHandle Sets the default context handle to the given handle function
func SetDefCtxHandle(handle CtxHandle) {
	atomic.StorePointer(&dch, unsafe.Pointer(&ctxHandler{handle: handle}))
}

// getDefCtxHandle returns the default context handle function
func getDefCtxHandle() CtxHandle {
	return (*ctxHandler)(atomic.LoadPointer(&dch)).handle
}

// emptyHandle is default context handle function that returns an empty slice of fields
func emptyHandle(ctx context.Context) []Field {
	return nil
}
