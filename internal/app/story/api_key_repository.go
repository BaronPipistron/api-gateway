package story

type APIKeysRepository interface {
	GetApiKeys(key string) (APIKey, error)
}
type APIKey struct {
	key     string
	enabled bool
}
