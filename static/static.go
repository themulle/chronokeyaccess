package static

import "embed"

//go:embed frontend/*
//go:embed templates/*
var StaticFs embed.FS


