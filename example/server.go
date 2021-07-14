package main

import "snap7go"

var (
	Server *snap7go.S7Server
	DB21   [512]byte
	DB103  [1280]byte
	DB3    [1024]byte
	cnt    byte = 0
)
