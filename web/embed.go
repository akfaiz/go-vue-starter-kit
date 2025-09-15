package web

import "embed"

//go:embed dist/* dist/**/*
var Dist embed.FS
