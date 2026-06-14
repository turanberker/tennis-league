package api

import (
	"fmt"
	"tennis-league/common/http/router"
)

func main() {
	serverConfig := router.LoadServerConfig()

	fmt.Printf("Starting server on port %s...\n", serverConfig.Port)
}
