/*
Copyright © 2024 Furkan Pehlivan <furkanpehlivan34@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/pehlicd/rbac-wizard/internal"
	_ "github.com/pehlicd/rbac-wizard/internal/statik"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"log"
	"net/http"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the server for the rbac-wizard",
	Long:  `Start the server for the rbac-wizard. This will start the server on the specified port and serve the frontend.`,
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		serve(port)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringP("port", "p", "8080", "Port to run the server on")
}

func serve(port string) {
	// Set up CORS
	c := cors.New(cors.Options{
		AllowOriginRequestFunc: func(r *http.Request, origin string) bool {
			// Implement your dynamic origin check here
			host := r.Host // Extract the host from the request
			allowedOrigins := []string{"http://localhost:" + port, "https://" + host, "http://localhost:3000"}
			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin {
					return true
				}
			}
			return false
		},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	})

	// Set up handlers
	http.HandleFunc("/", dashboardHandler)
	http.HandleFunc("/api/data", dataHandler)

	handler := c.Handler(http.DefaultServeMux)

	// Start the server
	serverURL := fmt.Sprintf("http://localhost:%s", port)
	startupMessage := fmt.Sprintf("Starting dashboard server on %s", serverURL)
	fmt.Println(startupMessage)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Set cache control headers

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	// Serve the embedded frontend
	http.FileServer(statikFS).ServeHTTP(w, r)
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	// Set cache control headers
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	// Get the bindings
	bindings, err := internal.GetBindings()
	if err != nil {
		http.Error(w, "Failed to get bindings", http.StatusInternalServerError)
		return
	}

	data := internal.GenerateData(bindings)
	byteData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Failed to marshal data", http.StatusInternalServerError)
		return
	}

	// Write the bindings to the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(byteData)
}
