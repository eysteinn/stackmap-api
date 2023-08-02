package main

import (
	"gitlab.com/EysteinnSig/stackmap-api/internal/api"
)

//go:generate swagger generate spec

func main() {
	api.Run("0.0.0.0:3000")
}
