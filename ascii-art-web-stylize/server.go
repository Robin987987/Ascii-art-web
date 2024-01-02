package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func formHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil { //parse the raw query
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	name := r.FormValue("name")
	selectedFont := r.FormValue("fonts")
	selectedFontInt, _ := strconv.Atoi(selectedFont)
	var txtFile string
	if selectedFontInt == 2 {
		txtFile = "banner/shadow.txt"
	} else if selectedFontInt == 3 {
		txtFile = "banner/thinkertoy.txt"
	} else {
		txtFile = "banner/standard.txt"
	}

	// Create Ascii Art
	if name == "\\n" { // convert "\n" to newline
		fmt.Fprintf(w, "\n")
	} else if name != "" {
		fontBytes, err := os.ReadFile(txtFile) // read Font File
		if err != nil {
			fmt.Fprint(w, "Error reading font file: ")
			return
		}
		lines := strings.Split(string(fontBytes), "\n") // split Font File
		name = strings.ReplaceAll(name, "\\n", "\n")    // handle "\n" input
		parts := strings.Split(name, "\n")              // split Input
		for _, line := range parts {
			for i := 1; i < 9; i++ {
				for _, ascii := range line {
					fmt.Fprintf(w, lines[(ascii-32)*9+rune(i)])
				}
				if line == "" { // handle empty line ("\n")
					fmt.Fprintln(w)
					break
				}
				fmt.Fprintln(w)
			}
		}
	}
}
func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/ascii-art" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
}
func main() {
	fileServer := http.FileServer(http.Dir("./static")) // create the file server
	http.Handle("/", fileServer)                        // accepts a path and the fileserver
	http.HandleFunc("/ascii-art", formHandler)
	http.HandleFunc("/hello", helloHandler)
	fmt.Printf("Starting server at port 8080\n") // starts the server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
