package playerhandler

type PlayerResponse struct {
	ID      int64  `json:"id"`
	UUID    string `json:"uuid"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	UserId  *int64 `json:"user_id"`
}
