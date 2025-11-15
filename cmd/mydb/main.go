package main

import (
	"log"

	// internal/api exposes SetupRouter which registers routes, middleware and handlers.
	"github.com/networks-security2526/lab3-base/internal/api"
)

func main() {
	// Initialize the router with all routes and middleware.
	router := api.SetupRouter()

	// Informational startup message (address shown for convenience).
	log.Println("Servidor escuchando en https://myserver.local:5000")

	// Start the HTTPS server using the provided certificate and private key.
	// RunTLS blocks until the server stops; a non-nil error means startup failed.
	err := router.RunTLS("localhost:5000", "cert.crt", "cert-priv.pem")
	if err != nil {
		// Fatal log to ensure the process exits on startup failure.
		log.Fatalf("Error al iniciar el servidor HTTPS: %v", err)
	}
}
