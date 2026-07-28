// Command server is the bootstrap entry point for the authoritative match host.
//
// Runtime transports, persistence adapters, and gameplay orchestration are added
// only after their contracts and follow-up technology decisions are accepted.
package main

import "log"

func main() {
	log.Print("unmatched authoritative server bootstrap: runtime not configured")
}
