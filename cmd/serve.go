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
	"gopkg.in/yaml.v2"
	"io"
	"k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	kyaml "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
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

var app internal.App

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringP("port", "p", "8080", "Port to run the server on")
}

func serve(port string) {
	kubeClient, err := internal.GetClientset()
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v\n", err)
	}

	app.KubeClient = kubeClient

	// Set up CORS
	c := cors.New(cors.Options{
		AllowOriginVaryRequestFunc: func(r *http.Request, origin string) (bool, []string) {
			// Implement your dynamic origin check here
			host := r.Host // Extract the host from the request
			allowedOrigins := []string{"http://localhost:" + port, "https://" + host, "http://localhost:3000"}
			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin {
					return true, []string{"Origin"}
				}
			}
			return false, nil
		},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	})

	// Set up statik filesystem
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	// Set up handlers
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveStaticFiles(statikFS, w, r, "index.html")
	})
	http.HandleFunc("/what-if", func(w http.ResponseWriter, r *http.Request) {
		serveStaticFiles(statikFS, w, r, "what-if.html")
	})
	http.HandleFunc("/api/data", dataHandler)
	http.HandleFunc("/api/what-if", whatIfHandler)

	handler := c.Handler(http.DefaultServeMux)

	// Start the server
	startupMessage := fmt.Sprintf("Starting rbac-wizard on %s", fmt.Sprintf("http://localhost:%s", port))
	fmt.Println(startupMessage)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}

func serveStaticFiles(statikFS http.FileSystem, w http.ResponseWriter, r *http.Request, defaultFile string) {
	// Set cache control headers
	cacheControllers(w)

	path := r.URL.Path
	if path == "/" {
		path = "/" + defaultFile
	}

	file, err := statikFS.Open(path)
	if err != nil {
		// If the file is not found, serve the default file (index.html)
		file, err = statikFS.Open("/" + defaultFile)
		if err != nil {
			http.NotFound(w, r)
			return
		}
	}
	defer file.Close()

	// Get the file information
	fileInfo, err := file.Stat()
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.ServeContent(w, r, path, fileInfo.ModTime(), file)
}

func dataHandler(w http.ResponseWriter, _ *http.Request) {
	// Set cache control headers
	cacheControllers(w)

	// Get the bindings
	bindings, err := internal.Generator(app).GetBindings()
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
	_, err = w.Write(byteData)
	if err != nil {
		http.Error(w, "Failed to write data", http.StatusInternalServerError)
		return
	}
}

func cacheControllers(w http.ResponseWriter) {
	// Set cache control headers
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
}

func whatIfHandler(w http.ResponseWriter, r *http.Request) {
	cacheControllers(w)

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	var input struct {
		Yaml string `json:"yaml"`
	}

	if err := json.Unmarshal(body, &input); err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	var obj interface{}
	if err := yaml.Unmarshal([]byte(input.Yaml), &obj); err != nil {
		http.Error(w, "Invalid YAML format", http.StatusBadRequest)
		return
	}

	var responseData struct {
		Nodes []internal.Node `json:"nodes"`
		Links []internal.Link `json:"links"`
	}

	if obj == nil {
		http.Error(w, "Empty object", http.StatusBadRequest)
		return
	}

	decode := kyaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	uObj := &unstructured.Unstructured{}
	_, _, err = decode.Decode([]byte(input.Yaml), nil, uObj)
	if err != nil {
		http.Error(w, "Failed to parse ClusterRoleBinding", http.StatusBadRequest)
		fmt.Printf("Failed to parse ClusterRoleBinding: %v\n", err)
		return
	}

	switch obj.(map[interface{}]interface{})["kind"] {
	case "ClusterRoleBinding":
		crb := &v1.ClusterRoleBinding{}
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(uObj.UnstructuredContent(), crb)
		if err != nil {
			http.Error(w, "Failed to convert to ClusterRoleBinding", http.StatusBadRequest)
			fmt.Printf("Failed to convert to ClusterRoleBinding: %v\n", err)
			return
		}

		responseData = internal.WhatIfGenerator(app).ProcessClusterRoleBinding(crb)
	case "RoleBinding":
		rb := &v1.RoleBinding{}
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(uObj.UnstructuredContent(), rb)
		if err != nil {
			http.Error(w, "Failed to convert to ClusterRoleBinding", http.StatusBadRequest)
			fmt.Printf("Failed to convert to ClusterRoleBinding: %v\n", err)
			return
		}
		responseData = internal.WhatIfGenerator(app).ProcessRoleBinding(rb)
	default:
		http.Error(w, "Unsupported resource type", http.StatusBadRequest)
		fmt.Println("Unsupported resource type ", obj, " of type ", fmt.Sprintf("%T", obj))
		return
	}

	respData, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, "Failed to marshal response data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(respData)
	if err != nil {
		http.Error(w, "Failed to write response data", http.StatusInternalServerError)
		return
	}
}
