package main

import (
    "fmt"
    "encoding/json"
    "net/http"
    "log"
    "bytes"
)

func msteamsReport(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/msteams" {
                http.Error(w, "404 not found.", http.StatusNotFound)
                return
        }
	w.Header().Add("Content-Type", "application/json")
        switch r.Method {
        case "GET":
		 w.Header().Add("Content-Type", "text/html")
                 http.ServeFile(w, r, "html/query.html")
        case "POST":
		w.Header().Add("Content-Type", "application/json")
		var messageData map[string]interface{}
		rLog := new(bytes.Buffer)
		err := json.NewDecoder(r.Body).Decode(&messageData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var msteamsRespond = make(map[string]interface{})
		msteamsRespond["text"] = "Acknowledged! Thanks " + messageData["from"].(map[string]interface{})["name"].(string)
                msteamsRespond["type"] = "message"
		json.NewEncoder(w).Encode(msteamsRespond)
		json.NewEncoder(rLog).Encode(messageData)
		fmt.Printf("recevide message: %v" ,rLog)
	        default:
                fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
        }
}

func ansibleIndex(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/" {
                http.Error(w, "404 not found.", http.StatusNotFound)
                return
        }

        switch r.Method {
        case "GET":
                 http.ServeFile(w, r, "html/index.html")
        case "POST":
                if err := r.ParseForm(); err != nil {
                        fmt.Fprintf(w, "ParseForm() err: %v", err)
                        return
                }
                fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
                name := r.FormValue("name")
                address := r.FormValue("address")
                fmt.Fprintf(w, "Name = %s\n", name)
                fmt.Fprintf(w, "Address = %s\n", address)
        default:
                fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
        }
}

func main() {
	argsPort := "80"
        http.HandleFunc("/msteams", msteamsReport)
	http.HandleFunc("/", msteamsReport)

        fmt.Printf("Starting server...\n")
        if err := http.ListenAndServe(":"+ argsPort, nil); err != nil {
                log.Fatal(err)
        }
}

