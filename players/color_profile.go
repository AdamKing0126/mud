package players

import "github.com/jmoiron/sqlx"

type ColorProfile struct {
	UUID        string
	Name        string
	Primary     string
	Secondary   string
	Warning     string
	Danger      string
	Title       string
	Description string
}

func (c *ColorProfile) GetUUID() string {
	return c.UUID
}

func (c *ColorProfile) GetColor(colorUse string) string {
	switch colorUse {
	case "primary":
		return c.Primary
	case "secondary":
		return c.Secondary
	case "warning":
		return c.Warning
	case "danger":
		return c.Danger
	case "title":
		return c.Title
	default:
		return "\033[0m"
	}
}

func NewColorProfileFromDB(db *sqlx.DB, uuid string) (*ColorProfile, error) {
	var colorProfile ColorProfile
	err := db.QueryRow("SELECT uuid, name, primary_color, secondary_color, warning_color, danger_color, title_color, description_color FROM color_profiles WHERE uuid = ?", uuid).
		Scan(&colorProfile.UUID, &colorProfile.Name, &colorProfile.Primary, &colorProfile.Secondary, &colorProfile.Warning, &colorProfile.Danger, &colorProfile.Title, &colorProfile.Description)
	if err != nil {
		return nil, err
	}
	return &colorProfile, nil
}
