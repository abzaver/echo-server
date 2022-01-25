package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"log"
	"os"
)

var md5hash = "f04aaafc54829f13626593888dd3d858"
var sha256hash = "3ce628ef40fd52229649b2657f8eb42740dca2a2790b1af986e430ffa68629eb"

func main() {
	f, err := os.Open("wordlist.txt")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		password := scanner.Text()
		hash := fmt.Sprintf("%x", md5.Sum([]byte(password)))
		if hash == md5hash {
			fmt.Printf("[+] Password found (MD5): %s\n", password)
		}

		hash = fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
		if hash == sha256hash {
			fmt.Printf("[+] Password found (SHA-256): %s\n", password)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}
}
