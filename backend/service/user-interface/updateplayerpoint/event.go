package updateplayerpoint

import "tennis-league/common/consumer"

type MatchPlayers struct {
	WinnerPlayerIds []string `json:"winnerPlayerIds"`
	LooserPlayerIds []string `json:"loserPlayerIds"`
}

const RoutingName_MatchApproved consumer.RoutingName = "MatchApproved"
