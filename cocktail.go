package main

import (
	"encoding/json"
	"fmt"
	"github.com/stianeikeland/go-rpio"
	"net/http"
	"path/filepath"
	"simonwaldherr.de/go/golibs/as"
	"simonwaldherr.de/go/golibs/cachedfile"
	"simonwaldherr.de/go/golibs/file"
	"simonwaldherr.de/go/golibs/gopath"
    //"simonwaldherr.de/go/golibs/xmath"
	"simonwaldherr.de/go/gwv"
	"strings"
	"time"
	waage "github.com/MichaelS11/go-hx711"
)

type Rezepte []struct {
	Name    string `json:"Name"`
	Zutaten []struct {
		Name  string `json:"Name"`
		Menge int    `json:"Menge"`
	} `json:"Zutaten"`
	Kommentar string `json:"Kommentar"`
}

type MultiStruct []int

var pins map[int]int
var zutaten map[string]int
//var multiplikator map[int]int
var multiplikator []int
var rezepte Rezepte
const aufladedauer =640
const zeitmultiplikator = 120

func init() {
    
    /*
    p := gpioreg.ByName("GPIO6")

    if err := p.Out(gpio.High); err != nil {
        fmt.Println(err)
    }
    fmt.Printf("%s is %s\n", p, p.Read())
    */
	
	pins = map[int]int{
		1:  2,
		2:  3,
		3:  4,
		4:  17,
		5:  27,
		6:  22,
		7:  10,
		8:  9,
		9:  11,
		10: 5,
		11: 6,
		12: 13,
		13: 19,
		14: 26,
		15: 21,
		16: 20,
	}

	zutaten = map[string]int{
		"Whisky":       2,
		"Zuckersirup":  3,
		"Gin":          4,
		"Zitronensaft": 5,
		"Tonic Water":  6,
		"Soda":         7,
		"Rum":          8,
	}
	
	/*
	multiplikator = map[int]int{
		1:  120,
		2:  120,
		3:  120,
		4:  120,
		5:  120,
		6:  120,
		7:  180,
	}*/
	
	err := rpio.Open()
	if err != nil {
		panic(fmt.Sprint("unable to open gpio", err.Error()))
	}

	str, _ := file.Read("./multiplikator.json")
	err = json.Unmarshal([]byte(str), &multiplikator)
	
	if err != nil {
		fmt.Println(err)
	}
	
	str, _ = file.Read("./rezepte.json")
	err = json.Unmarshal([]byte(str), &rezepte)
	
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Rezepte geladen:\n")
		//pumpe := rpio.Pin(pins[16])
		//pumpe.Output()
		//time.Sleep(time.Second * 2)
		/*
		for _, cocktail := range rezepte {
			fmt.Printf("Cocktail: %#v\n", cocktail.Name)
			fmt.Printf("  Zutaten:\n", )
			for _, zut := range cocktail.Zutaten {
				fmt.Printf("    %v: %v\n", zut.Name, zut.Menge)
				zutatPin := rpio.Pin(pins[zutaten[zut.Name]])
				zutatPin.Output()
				time.Sleep(time.Millisecond * aufladedauer)
				time.Sleep(time.Millisecond * time.Duration(zut.Menge) / 20)
				zutatPin.Input()
			}
			
			fmt.Printf("  Kommentar: %#v\n\n", cocktail.Kommentar)
		}
		*/
		//time.Sleep(time.Second * 2)
		//pumpe.Input()
		fmt.Printf("Ende.\n\n")
	}
}

func main() {
	dir := gopath.Dir()
	fmt.Println("DIR 1:", gopath.WD())
	fmt.Println("DIR 2:", dir)
	HTTPD := gwv.NewWebServer(8080, 60)

	fmt.Println("opening gpio")
	
	pumpe := rpio.Pin(pins[16])
    master := rpio.Pin(pins[15])
    entluft := rpio.Pin(pins[14])
    pumpe.Input()
    master.Input()
    entluft.Input()
    
    //waage.Reset()
    
    err := waage.HostInit()
    if err != nil {
        fmt.Println("HostInit error:", err)
        return
    }
        
    hx711, err := waage.NewHx711("GPIO6", "GPIO5")
    //hx711.Shutdown()
    //hx711, err = waage.NewHx711("GPIO6", "GPIO5")
    if err != nil {
        fmt.Println("NewHx711 error:", err)
        return
    }
    
    defer hx711.Shutdown()
    
    err = hx711.Reset()
    if err != nil {
        fmt.Println("Reset error:", err)
        return
    }
    /*
    var weight1 float64
    var weight2 float64
    
    weight1 = 200
    weight2 = 1000
    
    //hx711.GetAdjustValues(weight1, weight2)
    */
    hx711.AdjustZero = 128663
    hx711.AdjustScale = 385.000000
    
    //var tara = []int{}
    
    var data float64
    for i := 0; i < 10; i++ {
        time.Sleep(200 * time.Microsecond)
    
        data, err = hx711.ReadDataMedian(11)
        if err != nil {
            fmt.Println("ReadDataRaw error:", err)
            continue
        }
    
        fmt.Println(data)
        //tara = append(tara, data)
    }
/*
    taraAvg := xmath.Round(xmath.Arithmetic(tara))
    
    for i := 0; i < 200; i++ {
        time.Sleep(200 * time.Microsecond)
    
        data, err = hx711.ReadDataRaw()
        if err != nil || data < 10 {
            fmt.Println("ReadDataRaw error:", err)
            continue
        }
    
        fmt.Printf("data: %v \tavg: %v \tout: %v\n", data,taraAvg, data-taraAvg)
    }
   */ 
    return

/*
	for i := 1; i < 17; i++ {
		pinCola := rpio.Pin(pins[i])
		pinCola.Output()
		time.Sleep(time.Second / 25)
		pinCola.Input()
	}
*/

	defer rpio.Close()

	pinCola := rpio.Pin(17)
	//pinCola.Output()

	HTTPD.URLhandler(
		gwv.URL("^/toggle/?$", func(rw http.ResponseWriter, req *http.Request) (string, int) {
			for i := 1; i < 17; i++ {
				pinCola := rpio.Pin(pins[i])
				pinCola.Output()
				time.Sleep(time.Second / 5)
				pinCola.Input()
			}
			return "/", http.StatusFound
		}, gwv.HTML),
		gwv.URL("^/ein/?$", func(rw http.ResponseWriter, req *http.Request) (string, int) {
			pinCola.Output()
			return "/", http.StatusFound
		}, gwv.HTML),
		gwv.URL("^/aus/?$", func(rw http.ResponseWriter, req *http.Request) (string, int) {
			pinCola.Input()
			return "/", http.StatusFound
		}, gwv.HTML),
		gwv.URL("^/select/?.*$", func(rw http.ResponseWriter, req *http.Request) (string, int) {
			pin := strings.Replace(req.RequestURI, "/select/", "", 1)
			pinCola = rpio.Pin(pins[int(as.Int(pin))])
			return "", http.StatusOK
		}, gwv.HTML),
		gwv.URL("^/list/?$", func(rw http.ResponseWriter, req *http.Request) (string, int) {
			var ret string
			for _, cocktail := range rezepte {
				ctname := strings.Replace(cocktail.Name, " ", "", -1)
				//if ctname == wunschCocktail {
					ret += fmt.Sprintf("<a href=\"../ozapftis/%v\">%v</a>\n", ctname, cocktail.Name)
					ret += fmt.Sprintf("<p>%v</p>\n\n", cocktail.Kommentar)
				//}
			}
			return ret, http.StatusOK
		}, gwv.HTML),
		gwv.URL("^/test/\\d*/\\d*$", func(rw http.ResponseWriter, req *http.Request) (string, int) {
			str, _ := file.Read("./multiplikator.json")
			err := json.Unmarshal([]byte(str), &multiplikator)
			
			if err != nil {
				fmt.Println(err)
			}
			
			
			pumpe.Output()
			time.Sleep(time.Second * 1)
			
			testStr := strings.Replace(req.RequestURI, "/test/", "", 1)
			testArr := strings.Split(testStr, "/")
			testPin := rpio.Pin(pins[int(as.Int(testArr[0]))])
			vorlaufdauer := time.Millisecond * aufladedauer
			ansteuerdauer := time.Millisecond * time.Duration(int(as.Int(testArr[1]))*int(as.Int(multiplikator[int(as.Int(testArr[0]))])))
			fmt.Printf("vorlaufdauer: %v\tansteuerdauer: %v\n", vorlaufdauer, ansteuerdauer)
			testPin.Output()
			time.Sleep(vorlaufdauer)
			time.Sleep(ansteuerdauer)
			testPin.Input()
			time.Sleep(time.Second * 1)
			pumpe.Input()
			return "", http.StatusOK
		}, gwv.HTML),
		gwv.URL("^/ozapftis/?.*$", func(rw http.ResponseWriter, req *http.Request) (string, int) {
			wunschCocktail := strings.Replace(req.RequestURI, "/ozapftis/", "", 1)
			pumpe := rpio.Pin(pins[16])
			pumpe.Output()
			time.Sleep(time.Second * 2)
			for _, cocktail := range rezepte {
				if strings.Replace(cocktail.Name, " ", "", -1) == wunschCocktail {
					fmt.Printf("Cocktail: %#v\n", cocktail.Name)
					fmt.Printf("  Zutaten:\n", )
					for _, zut := range cocktail.Zutaten {
						fmt.Printf("    %v: %v\n", zut.Name, zut.Menge)
						zutatPin := rpio.Pin(pins[zutaten[zut.Name]])
						vorlaufdauer := time.Millisecond * aufladedauer
						ansteuerdauer := time.Millisecond * time.Duration(zut.Menge * zeitmultiplikator)//time.Duration(zut.Menge*int(as.Int(multiplikator[pins[zutaten[zut.Name]]])))
						fmt.Printf("vorlaufdauer: %v\tansteuerdauer: %v\n", vorlaufdauer, ansteuerdauer)
						zutatPin.Output()
						time.Sleep(vorlaufdauer)
						time.Sleep(ansteuerdauer)
						//time.Sleep(time.Millisecond * aufladedauer)
						//time.Sleep(time.Millisecond * time.Duration(zut.Menge) * zeitmultiplikator)
						zutatPin.Input()
						time.Sleep(time.Second * 1)
					}
					fmt.Printf("  Kommentar: %#v\n\n", cocktail.Kommentar)
				}
			}
			time.Sleep(time.Second * 2)
			pumpe.Input()
			fmt.Printf("Ende.\n\n")
			return "/", http.StatusFound
		}, gwv.HTML),
		gwv.URL("^/$", func(rw http.ResponseWriter, req *http.Request) (string, int) {
			return as.String(cachedfile.Read(filepath.Join(dir, "index.html"))), http.StatusOK
		}, gwv.HTML),
		gwv.Robots(as.String(cachedfile.Read(filepath.Join(dir, "..", "static", "robots.txt")))),
		gwv.Favicon(filepath.Join(dir, "..", "static", "favicon.ico")),
		gwv.StaticFiles("/", dir),
	)

	HTTPD.Start()
	HTTPD.WG.Wait()
}