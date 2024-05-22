package env

import (
	"encoding/json"
	"fmt"
	"os"
)

type jsonHandler struct {
	loadedKeys map[Key]string
}

// NewJsonHandler returns a env-handler from loaded file
// the file MUST contain only string values (no support for nested object)
func NewJsonHandler(filePath string) (*jsonHandler, error) {
	buf, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("load env-json: %w", err)
	}

	var allKeys map[Key]string
	if err = json.Unmarshal(buf, &allKeys); err != nil {
		return nil, fmt.Errorf("unmarshall env-json: %w", err)
	}

	return &jsonHandler{allKeys}, nil
}

func (j *jsonHandler) Env(key Key, def string) string {
	if v, ok := j.loadedKeys[key]; ok {
		return v
	}

	return def
}

func (j *jsonHandler) MustEnv(key Key) string {
	if v, ok := j.loadedKeys[key]; ok {
		return v
	}

	panic("key not found")
}
