package main

import (
    "fmt"
    "encoding/json"
    "net/http"
    "log"
    "bytes"
    "time"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

func msteamsReport(w http.ResponseWriter, r *http.Request) {
	database, err := sql.Open("mysql", "rw_annotation:ikljnwenkrkw9023nefi2904@tcp(10.202.60.40:3306)/annotations")
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	database.SetConnMaxLifetime(time.Minute * 3)
	database.SetMaxOpenConns(10)
	database.SetMaxIdleConns(10)
	//Check if path correct
        if r.URL.Path != "/msteams" {
                http.Error(w, "404 not found.", http.StatusNotFound)
                return
        }
	w.Header().Add("Content-Type", "application/json")
        switch r.Method {
        case "GET":
		 w.Header().Add("Content-Type", "text/html")
                 http.ServeFile(w, r, "/var/html/msteams-index.html")
        case "POST":
		//Set up headers to provide json
		w.Header().Add("Content-Type", "application/json")
		//Set up variables
		var messageData map[string]interface{}
		rLog := new(bytes.Buffer)
		//Make sure we handle non json request body
		err := json.NewDecoder(r.Body).Decode(&messageData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		messageAsText := fmt.Sprintf("%v", messageData["text"])
		reporterAsText := messageData["from"].(map[string]interface{})["name"].(string)
		_, sqlErr := database.Exec("INSERT INTO chreports(report) VALUES('" + messageAsText + " by " + reporterAsText + "')")
                if sqlErr != nil {
			fmt.Println("SQL error: ")
			log.Fatal(sqlErr)
                }
		json.NewEncoder(w).Encode(reportResponds(messageData["from"].(map[string]interface{})["name"].(string)))
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

