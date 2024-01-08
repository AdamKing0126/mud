package interfaces

type ColorProfile interface {
	GetUUID() string
	GetColor(string) string
}
