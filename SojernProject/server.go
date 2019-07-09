package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

/*
Application - 'math api'
Implement a web service (preferably in Go, Python, Ruby or Java;
extra effort to do that in Go will be recognised; using a framework or not):
/min - given list of numbers and a quantifier (how many) provides min number(s)
/max - given list of numbers and a quantifier (how many) provides max number(s)
/avg - given list of numbers calculates their average
/median - given list of numbers calculates their median
/percentile - given list of numbers and quantifier 'q', compute the qth percentile of the list elements
No need to be concerned with resources, we're assuming there's plenty enough of memory, etc.
 */


type op struct {
	Quantifier string `json: quantifier`
	List []int `json: list`
}

func max(w http.ResponseWriter, r *http.Request){
	max := op{}
	_ = json.NewDecoder(r.Body).Decode(&max)
	log.Printf("Received: %v\n", max)
	sort.Ints(max.List)
	num,_ := strconv.ParseInt(max.Quantifier, 10, 32)
	aws := max.List[len(max.List) - int(num) : len(max.List)]
	str := strings.Trim(strings.Replace(fmt.Sprint(aws), " ", ",", -1), "[]")
	log.Println(str)
	s := fmt.Sprintf("%s", str)
	fmt.Fprintf(w, s)

}

func min(w http.ResponseWriter, r *http.Request){
	min := op{}
	_ = json.NewDecoder(r.Body).Decode(&min)
	log.Printf("Received: %v\n", min)
	sort.Ints(min.List)
	num,_ := strconv.ParseInt(min.Quantifier, 10, 32)
	aws := min.List[0 : num]
	str := strings.Trim(strings.Replace(fmt.Sprint(aws), " ", ",", -1), "[]")
	log.Println(str)
	s := fmt.Sprintf("%s", str)
	fmt.Fprintf(w, s)

}

func median(w http.ResponseWriter, r *http.Request){
	median := op{}
	_ = json.NewDecoder(r.Body).Decode(&median)
	log.Printf("Received: %v\n", median)
	sort.Ints(median.List)
	var aws float64
	if len(median.List) %2 ==0{
		aws = (float64(median.List[len(median.List)/2]) + float64(median.List[len(median.List)/2 -1]))/2

	} else{
		aws = float64(median.List[len(median.List)/2])
	}

	log.Println(aws)
	s := fmt.Sprintf("%f", aws)
	fmt.Fprintf(w, s)

}

func percentile(w http.ResponseWriter, r *http.Request){
	percentile := op{}
	_ = json.NewDecoder(r.Body).Decode(&percentile)
	log.Printf("Received: %v\n", percentile)
	sort.Ints(percentile.List)
	num,_ := strconv.ParseInt(percentile.Quantifier, 10, 32)
	index := math.Ceil(float64(int(num)* (len(percentile.List)-1)) / 100)
	//log.Println(index)
	aws  := percentile.List[int(index)]
	log.Println(aws)
	s := fmt.Sprintf("%d", aws)
	fmt.Fprintf(w, s)

}

func avg(w http.ResponseWriter, r *http.Request) {

	avg := op{}
	_ = json.NewDecoder(r.Body).Decode(&avg)
	log.Printf("Received: %v\n", avg)

	sum :=0
	for _, val := range avg.List{
		sum+= val
	}
	num  := len(avg.List)
	aws  := float64(sum)/float64(num)
	log.Println(aws)
	s := fmt.Sprintf("%f", aws)
	fmt.Fprintf(w, s)


}


func server() {
		mux := http.NewServeMux()
		mux.HandleFunc("/avg", avg)
		mux.HandleFunc("/max", max)
		mux.HandleFunc("/min", min)
		mux.HandleFunc("/median", median)
		mux.HandleFunc("/percentile", percentile)
		http.ListenAndServe(":8088", mux)
}


// q is  a quantifier corresponding to each function
// list is input, which in this implementation is []int
// operation can be /max, /min....etc
// responding answer is String

func client() {

	q := "2"
	list := []int{8,2,35,5}
	operation := "/percentile"
	url := "http://localhost:8088"+ operation

	b := new(bytes.Buffer)
	input := op{Quantifier: q, List: list}
	json.NewEncoder(b).Encode(input)
	req, _ := http.NewRequest("POST", url, b)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err :=client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("Response: ", string(body))
	resp.Body.Close()

}
func main() {
	go server()
	client()

}
