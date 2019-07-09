package main
/*
Implement a small small webserver (preferably in Go, Python, Ruby, or Java; extra effort to do that in Go will be recognised):
/ping - returns response code 200 and string OK when file /tmp/ok is present, if file is not present returns 503 service unavailable
/img - returns a 1x1 gif image, and log the request in apache common log format
Server needs to scale to many concurrent users and be efficient. Propose improvements you'd like to work on and - time permits - implement.
 */
import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// main page
func mainPage (w http.ResponseWriter, r * http.Request){
	fmt.Fprintf(w, "Welcome")
}

//if /temp/ok is present, set statusCode 200 and return "OK". Otherwise, return 503 along with Service Not available
func ping(w http.ResponseWriter, r *http.Request){
		_, err := ioutil.ReadFile("tmp/ok")
		if err!=nil{
			log.Println(err)
			w.WriteHeader(503)
			fmt.Fprintf(w, "Service Not available")
		}else{
			w.WriteHeader(200)
			fmt.Fprintf(w, "OK")
		}
}
// return image.gif and log request info.
func image(w http.ResponseWriter, r *http.Request){

	t := time.Now()
	img, err := os.Open("sojern.gif")
	if err != nil {
		log.Fatal(err) // perhaps handle this nicer
	}
	defer img.Close()
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "image/gif") // <-- set the content-type header
	_ , err = io.Copy(w, img)

	log.Println(w.Header())
	log.Printf("\n%s\n%s\n%s %s\n%v\n%d",
		r.RemoteAddr,
		r.RequestURI,
		r.Method,
		r.Proto,
		t, r.ContentLength)

}
func server()  {
	mux := http.NewServeMux()
	mux.HandleFunc("/", mainPage)
	mux.HandleFunc("/ping", ping)
	mux.HandleFunc("/img", image)
	http.ListenAndServe("127.0.0.1:8000", mux)
}

func main(){
server()

}
