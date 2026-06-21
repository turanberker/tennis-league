package playerpoint

type AddPlayerPoint struct {
	PlayerId    string
	EarnedPoint int32
	ScoreType   SCORE_TYPE
}

type PlayerPoints struct {
	ID          string
	DoublePoint int
	SinglePoint int
}

type matchParticipant struct {
	PlayerID    string
	DoublePoint int
	IsWinner    bool
}

type SCORE_TYPE string

const (
	SCORE_TYPE_SINGLE SCORE_TYPE = "SINGLE"
	SCORE_TYPE_DOUBLE SCORE_TYPE = "DOUBLE"
)
