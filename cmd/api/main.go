package main

import (
	"gitlab.com/EysteinnSig/stackmap-api/internal/api"
)

//go:generate swagger generate spec

func main() {
	api.Run(":3000")
}
