package clubhandler

type PersistClub struct {
	Name string `json:"name" binding:"min=3,max=75,required"`
}

type ClupResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
