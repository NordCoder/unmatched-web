// Command server runs the local two-tab playable slice while keeping the
// approved PostgreSQL Core adapter available and unchanged.
package main

import (
	"log"
	"net/http"
	"os"

	webapp "github.com/NordCoder/unmatched-web/apps/web"
	"github.com/NordCoder/unmatched-web/internal/playableslice/content"
	sliceserver "github.com/NordCoder/unmatched-web/internal/playableslice/server"
)

func main() {
	registry, err := content.Load()
	if err != nil {
		log.Fatalf("load playable-slice content: %v", err)
	}
	address := sliceserver.RunAddress(os.Getenv("PORT"))
	log.Printf("Robin Hood vs Bigfoot playable slice listening on http://localhost%s", address)
	if err := http.ListenAndServe(address, sliceserver.NewHandler(registry, webapp.Static())); err != nil {
		log.Fatal(err)
	}
}
