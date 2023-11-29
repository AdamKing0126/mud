package interfaces

type ActionInterface interface {
	GetPlayer() PlayerInterface
	GetCommand() string
}
