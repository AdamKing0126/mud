package areas

type Area struct {
	UUID        string
	Name        string
	Description string
	Rooms       []Room
	Channel     chan Action
}

func (a *Area) GetUUID() string {
	return a.UUID
}

func (a *Area) GetName() string {
	return a.Name
}

func (a *Area) GetDescription() string {
	return a.Description
}

type AreaInfo struct {
	UUID        string
	Name        string
	Description string
}

func NewArea(uuid string, name string, description string) *Area {
	return &Area{UUID: uuid, Name: name, Description: description}
}
