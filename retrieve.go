package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
)

func main() {
	resp, err := http.Get("https://www.pathe-thuis.nl/movie/index/browse?mainURL=collectie&subURL=81&page=1&amount=30")

	fmt.Println("http transport error:", err)
	
	body, err := ioutil.ReadAll(resp.Body)
	
	fmt.Println("read error is:", err)
	
	fmt.Println(string(body))
	
	

}
