package env

import "os"

type osHandler struct {
	provider func(key Key, def string) string
}

// NewOsHandler returns a default, os-env based handler
func NewOsHandler() (*osHandler, error) {
	return &osHandler{}, nil
}

func (h *osHandler) Env(key Key, def string) string {
	if e, ok := os.LookupEnv(string(key)); ok {
		return e
	}

	return def
}

func (h *osHandler) MustEnv(key Key) string {
	if e, ok := os.LookupEnv(string(key)); ok {
		return e
	}

	panic("key not found")
}
