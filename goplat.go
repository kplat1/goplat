package main

import (
	"fmt"
	"net/http"
	"log"
	    "encoding/binary"
    "bytes"

	bolt "github.com/coreos/bbolt"
)

type LoginRec struct {
  Username string
  Password string
  Points float64
}

func (lr *LoginRec) Bytes() {
 
}

func SaveNewLogin(db *bolt.DB, rec *LoginRec) {
  db.Update(func(tx *bolt.Tx) error {
    b, err := tx.CreateBucketIfNotExists([]byte("LoginTable"))
  	var bin_buf bytes.Buffer
    binary.Write(&bin_buf, binary.BigEndian, rec)
  	err = b.Put([]byte(rec.Username), bin_buf.Bytes())
  	return err
  })
}


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

  rec := LoginRec{Username: "kai", Password: "kk", Points: 100.2}
  
  SaveNewLogin(db, &rec) // & = get the pointer to the record

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	  switch r.URL.Path {

    case "/KPlat.html":
	    fmt.Fprintf(w, "%v\nUsername is: %v\n", KPlat, Username)
	  }
	})
	http.ListenAndServe(":2083", nil)

}
