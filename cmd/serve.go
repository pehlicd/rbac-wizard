/*
Modified by Alessio Greggi © 2025. Based on work by Furkan Pehlivan <furkanpehlivan34@gmail.com>.

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
	"io"
	"io/fs"
	"net/http"

	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	kyaml "k8s.io/apimachinery/pkg/runtime/serializer/yaml"

	embedfiles "github.com/alegrey91/rancher-rbac-wizard/internal/embed"
	"github.com/alegrey91/rancher-rbac-wizard/internal/logger"
	internal "github.com/alegrey91/rancher-rbac-wizard/internal/rrw"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the server for the rbac-wizard",
	Long:  `Start the server for the rbac-wizard. This will start the server on the specified port and serve the frontend.`,
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		enableLogging, _ := cmd.Flags().GetBool("logging")
		logLevel, _ := cmd.Flags().GetString("log-level")
		logFormat, _ := cmd.Flags().GetString("log-format")
		serve(port, enableLogging, logLevel, logFormat)
	},
}

var app internal.App

type Serve struct {
	App internal.App
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringP("port", "p", "8080", "Port to run the server on")
	serveCmd.Flags().BoolP("logging", "g", false, "Enable logging")
	serveCmd.Flags().StringP("log-level", "l", "info", "Log level")
	serveCmd.Flags().StringP("log-format", "f", "text", "Log format default is text [text, json]")
}

func serve(port string, logging bool, logLevel string, logFormat string) {
	// Set up logger if logging is enabled
	if logging {
		l := logger.New(logLevel, logFormat)
		app.Logger = l
	} else {
		l := logger.New("off", logFormat)
		app.Logger = l
	}

	kubeClient, dynamicClient, err := internal.GetClientset()
	if err != nil {
		app.Logger.Fatal().Err(err).Msg("Failed to create Kubernetes client")
	}

	// Set up clients
	app.KubeClient = kubeClient
	app.DynamicClient = dynamicClient
	serve := Serve{
		app,
	}

	// define subtree directory
	uiFS, err := fs.Sub(embedfiles.UIfs, "dist")
	if err != nil {
		app.Logger.Fatal().Err(err).Msg("Failed to find subtree in embedded directory")
	}

	// Set up embedded filesystem
	uiFiles := http.FS(uiFS)
	if err != nil {
		app.Logger.Fatal().Err(err).Msg("Failed to create embed filesystem")
	}

	// Set up CORS
	c := setupCors(port)

	// Create a new serve mux
	mux := http.NewServeMux()

	// Set up handlers
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveStaticFiles(uiFiles, w, r, "index.html")
	})
	mux.HandleFunc("/what-if", func(w http.ResponseWriter, r *http.Request) {
		serveStaticFiles(uiFiles, w, r, "what-if.html")
	})
	mux.HandleFunc("/api/data", serve.dataHandler)
	mux.HandleFunc("/api/what-if", serve.whatIfHandler)

	handler := c.Handler(serve.App.LoggerMiddleware(mux))

	// Start the server
	startupMessage := fmt.Sprintf("Starting rancher-rbac-wizard on %s", fmt.Sprintf("http://localhost:%s", port))
	fmt.Println(startupMessage)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		app.Logger.Fatal().Err(err).Msg("Failed to create embed filesystem")
	}
}

func setupCors(port string) *cors.Cors {
	return cors.New(cors.Options{
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
}

func serveStaticFiles(staticFS http.FileSystem, w http.ResponseWriter, r *http.Request, defaultFile string) {
	// Set cache control headers
	cacheControllers(w)

	path := r.URL.Path
	if path == "/" {
		path = "/" + defaultFile
	}

	file, err := staticFS.Open(path)
	if err != nil {
		// If the file is not found, serve the default file (index.html)
		file, err = staticFS.Open("/" + defaultFile)
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

func (s *Serve) dataHandler(w http.ResponseWriter, _ *http.Request) {
	// Set cache control headers
	cacheControllers(w)

	// Get the bindings
	bindings, err := internal.Generator(app).GetBindings()
	if err != nil {
		s.App.Logger.Error().Err(err).Msg("Failed to get bindings")
		http.Error(w, "Failed to get bindings", http.StatusInternalServerError)
		return
	}

	data := internal.GenerateData(bindings)
	byteData, err := json.Marshal(data)
	if err != nil {
		s.App.Logger.Error().Err(err).Msg("Failed to marshal data")
		http.Error(w, "Failed to marshal data", http.StatusInternalServerError)
		return
	}

	// Write the bindings to the response
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(byteData)
	if err != nil {
		s.App.Logger.Error().Err(err).Msg("Failed to write data")
		http.Error(w, "Failed to write data", http.StatusInternalServerError)
		return
	}
}

func (s *Serve) whatIfHandler(w http.ResponseWriter, r *http.Request) {
	cacheControllers(w)

	if r.Method != http.MethodPost {
		s.App.Logger.Error().Msg("Invalid request method")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.App.Logger.Error().Err(err).Msg("Failed to read request body")
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	var input struct {
		Yaml string `json:"yaml"`
	}

	if err := json.Unmarshal(body, &input); err != nil {
		s.App.Logger.Error().Err(err).Msg("Failed to parse JSON")
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	var obj interface{}
	if err := yaml.Unmarshal([]byte(input.Yaml), &obj); err != nil {
		s.App.Logger.Error().Err(err).Msg("Invalid YAML format")
		http.Error(w, "Invalid YAML format", http.StatusBadRequest)
		return
	}

	var responseData struct {
		Nodes []internal.Node `json:"nodes"`
		Links []internal.Link `json:"links"`
	}

	if obj == nil {
		s.App.Logger.Error().Msg("Empty object")
		http.Error(w, "Empty object", http.StatusBadRequest)
		return
	}

	decode := kyaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	uObj := &unstructured.Unstructured{}
	_, _, err = decode.Decode([]byte(input.Yaml), nil, uObj)
	if err != nil {
		s.App.Logger.Error().Err(err).Msg("Failed to parse object")
		http.Error(w, "Failed to parse object", http.StatusBadRequest)
		return
	}

	switch obj.(map[interface{}]interface{})["kind"] {
	case "ClusterRoleBinding":
		crb := &v1.ClusterRoleBinding{}
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(uObj.UnstructuredContent(), crb)
		if err != nil {
			s.App.Logger.Error().Err(err).Msg("Failed to convert to ClusterRoleBinding")
			http.Error(w, "Failed to convert to ClusterRoleBinding", http.StatusBadRequest)
			return
		}

		responseData = internal.WhatIfGenerator(app).ProcessClusterRoleBinding(crb)
	case "RoleBinding":
		rb := &v1.RoleBinding{}
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(uObj.UnstructuredContent(), rb)
		if err != nil {
			s.App.Logger.Error().Err(err).Msg("Failed to convert to RoleBinding")
			http.Error(w, "Failed to convert to RoleBinding", http.StatusBadRequest)
			return
		}
		responseData = internal.WhatIfGenerator(app).ProcessRoleBinding(rb)
	case "ProjectRoleTemplateBinding":
		prtb := &v3.ProjectRoleTemplateBinding{}
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(uObj.UnstructuredContent(), prtb)
		if err != nil {
			s.App.Logger.Error().Err(err).Msg("Failed to convert to ProjectRoleTemplateBinding")
			http.Error(w, "Failed to convert to ProjectRoleTemplateBinding", http.StatusBadRequest)
			return
		}
		responseData = internal.WhatIfGenerator(app).ProcessProjectRoleTemplateBinding(prtb)
	case "ClusterRoleTemplateBinding":
		crtb := &v3.ClusterRoleTemplateBinding{}
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(uObj.UnstructuredContent(), crtb)
		if err != nil {
			s.App.Logger.Error().Err(err).Msg("Failed to convert to ClusterRoleTemplateBinding")
			http.Error(w, "Failed to convert to ClusterRoleTemplateBinding", http.StatusBadRequest)
			return
		}
		responseData = internal.WhatIfGenerator(app).ProcessClusterRoleTemplateBinding(crtb)
	case "GlobalRoleBinding":
		grb := &v3.GlobalRoleBinding{}
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(uObj.UnstructuredContent(), grb)
		if err != nil {
			s.App.Logger.Error().Err(err).Msg("Failed to convert to GlobalRoleBinding")
			http.Error(w, "Failed to convert to GlobalRoleBinding", http.StatusBadRequest)
			return
		}
		responseData = internal.WhatIfGenerator(app).ProcessGlobalRoleBinding(grb)
	default:
		s.App.Logger.Error().Msg("Unsupported resource type")
		http.Error(w, "Unsupported resource type", http.StatusBadRequest)
		return
	}

	respData, err := json.Marshal(responseData)
	if err != nil {
		s.App.Logger.Error().Err(err).Msg("Failed to marshal response data")
		http.Error(w, "Failed to marshal response data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(respData)
	if err != nil {
		s.App.Logger.Error().Err(err).Msg("Failed to write response data")
		http.Error(w, "Failed to write response data", http.StatusInternalServerError)
		return
	}
}

func cacheControllers(w http.ResponseWriter) {
	// Set cache control headers
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
}
