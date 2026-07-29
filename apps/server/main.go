// Command server is the assembly point for the authoritative match host.
// Wave 1 intentionally leaves network transport and durable adapters unselected.
package main

import "log"

func main() {
	log.Print("unmatched core runtime Wave 1 is available; transport assembly is deferred")
}
