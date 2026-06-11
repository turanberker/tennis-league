package api

import (
	"fmt"
	"tennis-league/common/lib/config"
)

func main() {
	serverConfig := config.LoadServerConfig()

	fmt.Printf("Starting server on port %s...\n", serverConfig.Port)
}
