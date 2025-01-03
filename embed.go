package mediatools_embed

import "embed"

//go:embed filters
var Filters embed.FS

//go:embed localize
var Translations embed.FS
