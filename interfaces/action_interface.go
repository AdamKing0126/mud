package interfaces

type Action interface {
	GetPlayer() Player
	GetCommand() string
	GetArguments() []string
}
