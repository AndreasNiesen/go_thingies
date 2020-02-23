package main

import(
	"fmt"
	"os"
	"strconv"
	"time"
	"strings"
	"net/http"
	"io/ioutil"
)

func main() {
	var written int = 0
	var str_to_write string = ""
	f, err := os.OpenFile("public_IP.log", os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0644)
	if(err != nil) {
		fmt.Println("Error opening file to write to.\nErrorcode:", err)
		os.Exit(0)
	}
	defer f.Close()
	var curTime string = ""
	var shortSleep bool = false
	var resp *http.Response
		
	for {
		//get current time as a string and empty str_to_write
		curTime = time.Now().String()
		str_to_write = ""
		
		//remove unnecessary info from the curTime string and append to str_to_write
		slcs := strings.Split(curTime, " ")
		for i := 0; i < len(slcs)-1; i++ {
			str_to_write = str_to_write + slcs[i] + " "
		}
		
		//get public ip via GET request to http://api.ipify.org 
		resp, err = http.Get("http://api.ipify.org")
		if(err != nil) {
			fmt.Println("Error getting public IP via \"http://api.ipify.org\".\nErrorcode:", err)
			fmt.Println("-- Will sleep for 1 sec and retry --")
			time.Sleep(1 * time.Second)
			continue
		}
		
		//read body (only containing public IP in this case)
		ip, err := ioutil.ReadAll(resp.Body)
		if(err != nil) {
			fmt.Println("Error while reading the response body.\nErrorcode:", err)
			fmt.Println("-- Will sleep for 1 sec and retry --")
			time.Sleep(1 * time.Second)
			resp.Body.Close()
			continue
		}
		resp.Body.Close()
		
		//readying str_to_write -- after these to lines = curTime + -- + publicIP + \n
		str_to_write = str_to_write + "-- " + string(ip)
		str_to_write += "\n"
		
		//write str_to_write to file
		written, err = f.WriteString(str_to_write)
		if(err != nil || written != len(str_to_write)) {
			fmt.Println("Error writing to file.")
			fmt.Println("Stringlength - Chars written: " + strconv.Itoa(len(str_to_write)) + " - " + strconv.Itoa(written))
			fmt.Println("Errorcode:", err)
			shortSleep = true
		}
		
		//short sleep if write failed -- maybe system had a weird hiccup -- , long sleep if write was successful
		if(shortSleep){
			time.Sleep(10 * time.Second)
			shortSleep = false
		} else {
			fmt.Println("Wrote \"" + strings.Replace(str_to_write,"\n", "",1) + "\" to file.")
			time.Sleep(10 * time.Minute)
		}
	}
}