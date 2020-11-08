package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
)

var targetHost string
var gopherType string
var useBool bool

func init() {
	const (
		defaultGopher      = "pocket"
		defaultGopherUsage = "the variety of gopher"
	)
	flag.StringVar(&gopherType, "gopher-type", defaultGopher, defaultGopherUsage)
	flag.StringVar(&gopherType, "g", defaultGopher, defaultGopherUsage+" (shorthand)")
	flag.StringVar(&targetHost, "target", "defaultValue", "usage target")
	flag.BoolVar(&useBool, "use-bool", true, "how to use Use Bool parameter")
}

func main() {
	fs := flag.NewFlagSet("ExampleFunc", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	var ip net.IP
	fs.Func("ip", "`IP address` to parse", func(s string) error {
		ip = net.ParseIP(s)
		if ip == nil {
			return errors.New("could not parse IP")
		}
		return nil
	})
	fs.Parse([]string{"-ip", "127.0.0.1"})
	fmt.Printf("{ip: %v, loopback: %t}\n\n", ip, ip.IsLoopback())

	// 256 is not a valid IPv4 component
	fs.Parse([]string{"-ip", "256.0.0.1"})
	fmt.Printf("{ip: %v, loopback: %t}\n\n", ip, ip.IsLoopback())
	/*
		flag.Parse()

		// loop over all arguments by index and value
		for i, arg := range os.Args {
			// print index and value
			fmt.Println("item", i, "is", arg)
		}

		// is targetHost defaultValue - then it wasn't set on the command line
		if targetHost == "defaultValue" {
			fmt.Println("target not set, using default value")
		} else {
			fmt.Println("target set to: ", targetHost)
		}

		if gopherType == "pocket" {
			fmt.Println("gopher type not set, using default value: ", gopherType)
		} else {
			fmt.Println("gopher type set to: ", gopherType)
		}

		if useBool {
			fmt.Println("use bool not set, using default value: ", useBool)
		} else {
			fmt.Println("guse bool set to: ", useBool)
		}
	*/
}
