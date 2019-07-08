package main
import (        
	"flag"
	"fmt"
	"net/http"
	"io/ioutil"
	"os"
	"regexp"
)               

func main() {
	flag.Parse()

	args := flag.Args()
	fmt.Println(args)	
	
	if len(args) < 1 {
		fmt.Println("Please specify start page.")
		os.Exit(1)
	}
	retrieve(args[0])
}

	
func retrieve(uri string) {
	resp, err := http.Get(uri)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	
	body, _ := ioutil.ReadAll(resp.Body)
	
//	match, _ := regexp.MatchString("<h4>([a-z]+)</h4>", string(body))
//	fmt.Println(match)

	r, _ := regexp.Compile("<h4>([a-z]+)</h4>")
//	fmt.Println(r.FindAllString("<h4>joshua</h4> <h4>moesa</h4>"+ string(body), -1))
	fmt.Printf("%#v\n", r.FindStringSubmatch(string(body)))

//	fmt.Println(string(body))
	
}