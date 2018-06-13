package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	bolt "github.com/coreos/bbolt"
	"github.com/goki/gi"
	"github.com/goki/gi/oswin"
	"github.com/goki/gi/oswin/driver"
	"github.com/goki/gi/units"
	"github.com/goki/ki"
)

type LoginRec struct {
	Username string
	Password string
	Points   float64
}

// master kplat database, so we don't have to pass it around as an argument everywhere
var KPlatDB *bolt.DB

// todo: you could check errors on all of these!

func SaveNewLogin(rec *LoginRec) {
	KPlatDB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("LoginTable"))
		jb, err := json.Marshal(rec) // converts rec to json, as bytes jb
		err = b.Put([]byte(rec.Username), jb)
		return err
	})
}

func LoadLoginTable() []*LoginRec {
	lt := make([]*LoginRec, 0, 100) // 100 is the starting capacity of slice -- increase if you expect more users.
	KPlatDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("LoginTable"))
		if b != nil {
			b.ForEach(func(k, v []byte) error {
				if v != nil {
					rec := LoginRec{}
					json.Unmarshal(v, &rec) // loads v value as json into rec
					lt = append(lt, &rec)   // adds record to login table
				}
				return nil
			})
		}
		return nil
	})
	return lt
}

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
	var err error
	KPlatDB, err = bolt.Open("KPlat.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer KPlatDB.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {

		case "/KPlat.html":
			fmt.Fprintf(w, "%v\n", KPlat)
		}
	})

	// go = run as a separate goroutine -- main flow keeps going
	go http.ListenAndServe(":2083", nil)

	driver.Main(func(app oswin.App) {
		gogirun()
	})
}

// gogirun runs a local gogi viewer of your databases..
func gogirun() {
	width := 1024
	height := 768

	rec := ki.Node{}          // receiver for events
	rec.InitName(&rec, "rec") // this is essential for root objects not owned by other Ki tree nodes

	win := gi.NewWindow2D("GoGi Kplat Viewer", width, height, true) // pixel sizes

	vp := win.WinViewport2D()
	updt := vp.UpdateStart()
	vp.Fill = true

	// style sheet!
	var css = ki.Props{
		"button": ki.Props{
			"background-color": gi.Color{255, 240, 240, 255},
		},
		"#combo": ki.Props{
			"background-color": gi.Color{240, 255, 240, 255},
		},
		".hslides": ki.Props{
			"background-color": gi.Color{240, 225, 255, 255},
		},
	}
	vp.CSS = css

	vlay := vp.AddNewChild(gi.KiT_Frame, "vlay").(*gi.Frame)
	vlay.Lay = gi.LayoutCol

	trow := vlay.AddNewChild(gi.KiT_Layout, "trow").(*gi.Layout)
	trow.Lay = gi.LayoutRow
	trow.SetStretchMaxWidth()

	trow.AddNewChild(gi.KiT_Stretch, "str1")
	title := trow.AddNewChild(gi.KiT_Label, "title").(*gi.Label)
	title.Text = "KPlat System Controls"
	title.SetProp("text-align", gi.AlignTop)
	title.SetProp("align-vert", gi.AlignTop)
	title.SetProp("font-family", "Times New Roman, serif")
	title.SetProp("font-weight", "bold")
	title.SetProp("font-size", units.NewValue(24, units.Pt))
	trow.AddNewChild(gi.KiT_Stretch, "str2")

	irow := vlay.AddNewChild(gi.KiT_Layout, "irow").(*gi.Layout)
	irow.Lay = gi.LayoutRow
	irow.SetStretchMaxWidth()
	irow.AddNewChild(gi.KiT_Stretch, "str1")
	instr := irow.AddNewChild(gi.KiT_Label, "instr").(*gi.Label)
	instr.Text = "Shortcuts: Control+Alt+P = Preferences, Control+Alt+E = Editor, Command +/- = zoom"
	instr.SetProp("text-align", gi.AlignTop)
	instr.SetProp("align-vert", gi.AlignTop)
	// instr.SetMinPrefWidth(units.NewValue(30, units.Ch))
	irow.AddNewChild(gi.KiT_Stretch, "str2")

	//////////////////////////////////////////
	//      Buttons

	vlay.AddNewChild(gi.KiT_Space, "blspc")
	blrow := vlay.AddNewChild(gi.KiT_Layout, "blrow").(*gi.Layout)
	blab := blrow.AddNewChild(gi.KiT_Label, "blab").(*gi.Label)
	blab.Text = "Actions:"

	brow := vlay.AddNewChild(gi.KiT_Layout, "brow").(*gi.Layout)
	brow.Lay = gi.LayoutRow
	brow.SetProp("align-horiz", gi.AlignLeft)
	// brow.SetProp("align-horiz", gi.AlignJustify)
	brow.SetStretchMaxWidth()

	viewlogins := brow.AddNewChild(gi.KiT_Button, "viewlogins").(*gi.Button)
	viewlogins.SetText("View LoginTable")
	viewlogins.ButtonSig.Connect(rec.This, func(recv, send ki.Ki, sig int64, data interface{}) {
		if sig == int64(gi.ButtonClicked) {
			lt := LoadLoginTable()
			gi.StructTableViewDialog(vp, &lt, true, nil, "Login Table", "", nil, nil, nil)
		}
	})

	addlogin := brow.AddNewChild(gi.KiT_Button, "addlogin").(*gi.Button)
	addlogin.SetText("Add Login")
	addlogin.ButtonSig.Connect(rec.This, func(recv, send ki.Ki, sig int64, data interface{}) {
		if sig == int64(gi.ButtonClicked) {
			rec := LoginRec{}
			gi.StructViewDialog(vp, &rec, nil, "Enter Login Info", "", recv, func(recv, send ki.Ki, sig int64, data interface{}) {
				if sig == int64(gi.DialogAccepted) {
					SaveNewLogin(&rec)
				}
			})
		}
	})

	quit := brow.AddNewChild(gi.KiT_Button, "quit").(*gi.Button)
	quit.SetText("Quit")
	quit.ButtonSig.Connect(rec.This, func(recv, send ki.Ki, sig int64, data interface{}) {
		if sig == int64(gi.ButtonClicked) {
			gi.PromptDialog(vp, "Quit", "Quit: Are You Sure?", true, true, recv, func(recv, send ki.Ki, sig int64, data interface{}) {
				if sig == int64(gi.DialogAccepted) {
					KPlatDB.Close()
					vp.Win.Quit()
				}
			})
		}
	})

	vp.UpdateEndNoSig(updt)

	win.StartEventLoop()
}
