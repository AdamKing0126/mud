package areas

type Room struct {
	UUID        string
	AreaUUID    string
	Name        string
	Description string
	ExitNorth   string
	ExitSouth   string
	ExitWest    string
	ExitEast    string
	ExitUp      string
	ExitDown    string
}

type RoomImport struct {
	UUID        string            `yaml:"uuid"`
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Exits       map[string]string `yaml:"exits"`
}
