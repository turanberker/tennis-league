package customerror

var CONSTRAINT_VALIDATION = "CONTS_100"

var INSUFFICIENT_PERMISSIONS = "AUTH_100"
var INVALID_CREDENTIAL = "AUTH_101"
var ErrSessionExpired = "AUTH_102"
var ErrCodeEmailAlreadyExists = "AUTH_102"

var ErrLeagueAlreadyExists = "LEAGUE_100"
var ErrLeagueAlreadyFixtureCreated = "LEAGUE_101"

var ErrUserDoesNotHavePlayer = "USER_100"
var ErrInvalidPassword = "USER_101"
var ErrUserNotExists = "USER_102"

var ErrMatchApprovedCanNotUpdateScore = "MATCH_100"
