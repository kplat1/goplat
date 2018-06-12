package main

import (
	"fmt"
	"net/http"
	"log"

	bolt "github.com/coreos/bbolt"
)

var Username = "kaio"

var KPlat = `
<head>
<title>KPlat</title>
</head>
<h1>Welcome to KPlat</h1>
<ul>
<li>Here is option 1</li>
<li>Option 2</li>
</ul/>
`

func main() {
  	db, err := bolt.Open("KPlat.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()


	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	  switch r.URL.Path {

    case "/KPlat.html":
	    fmt.Fprintf(w, "%v\nUsername is: %v\n", KPlat, Username)
	  }
	})
	http.ListenAndServe(":2083", nil)

}