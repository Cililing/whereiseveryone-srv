package env

type Key string

type Handler interface {
	Env(key Key, def string) string
	MustEnv(key Key) string
}

func NewHandler(filePath string) (Handler, error) {
	if filePath == "" {
		return NewOsHandler()
	}

	return NewJSONHandler(filePath)
}
