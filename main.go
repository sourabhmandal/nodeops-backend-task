package main

import (
	"fmt"
	"net/http"
	"taikox/docker"
	"text/template"
)

type PageData struct {
	Title string
}

func main() {

	// Define the handler function for the landing page.
	landingPageHandler := func(w http.ResponseWriter, r *http.Request) {
		// Load and parse the template from a separate file
		var tmplFile = "index.tmpl"
		t, err := template.New(tmplFile).ParseFiles(tmplFile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Execute the template with the data.
		data := PageData{
			Title: "Taiko Node Starter",
		}

		err = t.Execute(w, data)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	// Set up the server.
	http.HandleFunc("/", landingPageHandler)
	http.HandleFunc("/start", docker.StartContainer)
	http.HandleFunc("/wait", docker.WaitForService)
	http.HandleFunc("/check", docker.CheckAvailableStatus)
	http.HandleFunc("/down", docker.StopContainer)

	// Start the server on port 8080.
	port := 8080
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Server is running on http://localhost%s\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}
