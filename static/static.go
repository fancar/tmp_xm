package static

import "embed"

// FS contains the static FS.
// integrations/* logo/* static/* swagger/* vendor/* *.json *.png *.html *.js
//go:embed  swagger/* vendor/* *.html
var FS embed.FS
