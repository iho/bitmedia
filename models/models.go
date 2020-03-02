package models

type User struct {
	Email     string `json:"email" faker:"email"`
	LastName  string `json:"last_name" faker:"last_name"`
	Country   string `json:"country"`
	City      string `json:"city"`
	Gender    string `json:"gender"`
	BirthDate string `json:"birth_date"`
}

type GameResult struct {
	PointsGained string `json:"points_gained"`
	WinStatus    string `json:"win_status"`
	GameType     string `json:"game_type"`
	Created      string `json:"created"`
}
