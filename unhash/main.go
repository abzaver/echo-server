package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha256"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

var (
	md5Hash      string
	sha256Hash   string
	dictFilename string
)

func init() {
	flag.StringVar(&md5Hash, "md5-hash", "", "md5-hashed password")
	flag.StringVar(&sha256Hash, "sha256-hash", "", "sha256-hashed password")
	flag.StringVar(&dictFilename, "f", "", "dictionary file name")
}

func main() {
	flag.Parse()
	//md5Hash = "f04aaafc54829f13626593888dd3d858"
	//sha256Hash = "3ce628ef40fd52229649b2657f8eb42740dca2a2790b1af986e430ffa68629eb"

	if md5Hash == "" && sha256Hash == "" {
		fmt.Println("You need to specify any password hashes. Use --md5-hash or --sha256-hash parameters")
	}

	f, err := os.Open(dictFilename)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	startTime := time.Now()
	scanner := bufio.NewScanner(f)
	linenum := 0
	for scanner.Scan() {
		var hash string = ""
		linenum++
		password := scanner.Text()

		if md5Hash != "" {
			hash = fmt.Sprintf("%x", md5.Sum([]byte(password)))
			if hash == md5Hash {
				fmt.Printf("[+] Password found (MD5) on line %d: %s\n", linenum, password)
			}
		}

		if sha256Hash != "" {
			hash = fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
			if hash == sha256Hash {
				fmt.Printf("[+] Password found (SHA-256) on line %d: %s\n", linenum, password)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Time spended: %f sec, lines scanned: %d", time.Since(startTime).Seconds(), linenum)
}
