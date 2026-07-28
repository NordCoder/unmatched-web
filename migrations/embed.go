// Package migrations embeds the ordered PostgreSQL schema used by Core Runtime.
package migrations

import "embed"

// Files contains every SQL migration in lexical order.
//
//go:embed *.sql
var Files embed.FS
