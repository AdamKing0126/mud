package interfaces

type ColorProfileInterface interface {
	GetUUID() string
	GetColor(string) string
}
