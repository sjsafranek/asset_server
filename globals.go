package main

import (
	"github.com/sjsafranek/ligneous"
)

const (
	DEFAULT_PORT = 1111
)

var (
	logger                  = ligneous.NewLogger()
	PORT             int    = DEFAULT_PORT
	ASSETS_DIRECTORY string = "assets/"
)
