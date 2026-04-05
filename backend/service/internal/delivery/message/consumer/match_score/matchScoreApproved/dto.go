package matchScoreApproved

import "github.com/turanberker/tennis-league-service/internal/domain/match"

type AddPlayerPoint struct {
	PlayerId    string
	EarnedPoint int32
	MatchType   match.Match_TYPE
}
