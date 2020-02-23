package main

import (
	"fmt"
	"os"
	"flag"
	"syscall"
	"encoding/hex"
	"unicode/utf8"
	"crypto/sha256"
	"golang.org/x/crypto/ssh/terminal"
)

// older version of the c++ password thingy reimagined in go.
// primary to get a bit more familiar with go
// far from good/useful
// same input -> same output
// not sorry for the name/pun

func btof(to_trans bool) (string) {
	if to_trans {
		return("set")
	} else {
		return("not set")
	}
}

func getNum(input string) (int){
	var output int = 0
	
	for c := 0; c < utf8.RuneCountInString(input); c++ {
		switch input[c] {
			case '0':
				output+=0
			case '1':
				output+=1
			case '2':
				output+=2
			case '3':
				output+=3
			case '4':
				output+=4
			case '5':
				output+=5
			case '6':
				output+=6
			case '7':
				output+=7
			case '8':
				output+=8
			case '9':
				output+=9
			case 'a':
				output+=10
			case 'b':
				output+=11
			case 'c':
				output+=12
			case 'd':
				output+=13
			case 'e':
				output+=14
			case 'f':
				output+=15
		}
	}
	
	return(output)
}

func getVec(input string, value int, flagnum uint) (int) {
	var output int = 0
	var flagdec int = 0
	
	switch input[1] {
		case '0':
			output+=0
		case '1':
			output+=1
		case '2':
			output+=2
		case '3':
			output+=3
		case '4':
			output+=4
		case '5':
			output+=5
		case '6':
			output+=6
		case '7':
			output+=7
		case '8':
			output+=8
		case '9':
			output+=9
		case 'a':
			output+=10
		case 'b':
			output+=11
		case 'c':
			output+=12
		case 'd':
			output+=13
		case 'e':
			output+=14
		case 'f':
			output+=15
	}
	
	output += value
	
	switch flagnum {
		case 8:
			flagdec = 2
		case 12:
			flagdec = 4
		case 13:
			flagdec = 5
		case 15:
			flagdec = 7
		default:
			flagdec = 2
	}
	
	for (output - flagdec) >= 0 {
		output -= flagdec
	}
	
	return(output)
}

func createPW(input string, flagnum uint) (string) {
	var alphanumshitz [7]string
	alphanumshitz[0] = "abcdefghijklm"
	alphanumshitz[1] = "nopqrstuvwxyz"
	alphanumshitz[2] = "ABCDEFGHIJKLM"
	alphanumshitz[3] = "NOPQRSTUVWXYZ"
	alphanumshitz[4] = "0123456789"
	alphanumshitz[5] = "()= !\\$%&/"
	alphanumshitz[6] = "@#°?\"{[]}*+'~"
		
	var output string = ""
	var buffer string = ""
	var chooseVec int = 0
	var substring int = 0
	
	for c := 0; c < utf8.RuneCountInString(input); c += 4 {
		buffer = input[c:c+4]
		substring = getNum(buffer)
		chooseVec = getVec(buffer, substring, flagnum)
		for substring >= utf8.RuneCountInString(alphanumshitz[chooseVec]) {
			substring -= utf8.RuneCountInString(alphanumshitz[chooseVec])
		}
		output += string([]rune(alphanumshitz[chooseVec])[substring])
	}
	
	return(output)
}

func main() {
	pName := flag.String("u", "", "your username")
	pSite := flag.String("s", "", "site to log into")
	pfNum := flag.Bool("nonum", false, "deactivate numbers for result-password")
	pfSpc := flag.Bool("nospc", false, "deactivate special characters for result-password")
	var flagnum uint = 12
	
	flag.Parse()
	
	if *pName == "" || *pSite == ""	{
		fmt.Println("\nYou need to specify a usernamename AND a site!")
		fmt.Println("---------------------------------------------------------------")
		fmt.Println(os.Args[0], "-u username -s site (optional flags:) -nonum -nospc")
		fmt.Println("---------------------------------------------------------------")
		os.Exit(0)
	}
	
	fmt.Print("Your password: ")
	bytePW, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		fmt.Println("Error in ReadPassword secretly!")
	} else if len(bytePW) == 0 {
		fmt.Println("\nPassword is not allowed to be empty!")
		os.Exit(0)
	}
	strPW := string(bytePW)
	
	compStr := *pName + strPW + *pSite
	sha_res := sha256.Sum256([]byte(compStr))
	if *pfNum == false {
		flagnum += 1
	}
	if *pfSpc == false {
		flagnum += 2
	}
	
	src := sha_res[:]
	dest := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dest, src)
	
	fmt.Printf("\nName = %s | Site = %s | Password = **** | NoNum-Flag = %s | NoSpecial-Flag = %s\n", *pName, *pSite, btof(*pfNum), btof(*pfSpc))
	final := createPW(string(dest), flagnum)

	fmt.Printf("%s\n", final)
}