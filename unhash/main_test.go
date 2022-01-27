package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// Utility function to execute the main function with given arguments and capture the output
func runMainWithArgs(args []string) (string, error) {
	cmd := exec.Command(os.Args[0], args...)
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

// Test function for main with dictionary generate
func TestMainFunctionGenerateDict(t *testing.T) {
	// Create a temporary dictionary file
	tmpfile, err := os.CreateTemp("", "dict")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	// Write sample passwords to the dictionary file
	passwords := []string{"password1", "password2", "password3"}
	for _, password := range passwords {
		if _, err := tmpfile.WriteString(password + "\n"); err != nil {
			t.Fatal(err)
		}
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Calculate hashes for one of the passwords
	md5Hash := fmt.Sprintf("%x", md5.Sum([]byte("password2")))
	sha256Hash := fmt.Sprintf("%x", sha256.Sum256([]byte("password2")))

	tests := []struct {
		args     []string
		expected string
	}{
		{
			args:     []string{"-md5-hash", md5Hash, "-f", tmpfile.Name()},
			expected: "[+] Password found (MD5) on line 2: password2",
		},
		{
			args:     []string{"-sha256-hash", sha256Hash, "-f", tmpfile.Name()},
			expected: "[+] Password found (SHA-256) on line 2: password2",
		},
		{
			args:     []string{"-md5-hash", md5Hash, "-sha256-hash", sha256Hash, "-f", tmpfile.Name()},
			expected: "[+] Password found (MD5) on line 2: password2\n[+] Password found (SHA-256) on line 2: password2",
		},
	}

	for _, tt := range tests {
		t.Run(strings.Join(tt.args, " "), func(t *testing.T) {
			out, err := runMainWithArgs(tt.args)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(out, tt.expected) {
				t.Fatalf("expected %q to contain %q", out, tt.expected)
			}
		})
	}
}

func TestMainFunctionExistedDict(t *testing.T) {
	// Existed hashes for one of the passwords
	md5Hash := "f04aaafc54829f13626593888dd3d858"
	sha256Hash := "3ce628ef40fd52229649b2657f8eb42740dca2a2790b1af986e430ffa68629eb"

	tests := []struct {
		args     []string
		expected string
	}{
		{
			args:     []string{"--md5-hash", md5Hash, "-f", "wordlist.txt"},
			expected: "[+] Password found (MD5)",
		},
		{
			args:     []string{"--sha256-hash", sha256Hash, "-f", "wordlist.txt"},
			expected: "[+] Password found (SHA-256)",
		},
		{
			args:     []string{"--md5-hash", md5Hash, "--sha256-hash", sha256Hash, "-f", "wordlist.txt"},
			expected: "[+] Password found (MD5) on line 5: vaapylvc\n[+] Password found (SHA-256) on line 5: vaapylvc",
		},
	}

	for _, tt := range tests {
		t.Run(strings.Join(tt.args, " "), func(t *testing.T) {
			out, err := runMainWithArgs(tt.args)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(out, tt.expected) {
				t.Fatalf("expected %q to contain %q", out, tt.expected)
			}
		})
	}
}

// Necessary for running exec.Command in tests
func TestMain(m *testing.M) {
	flag.Parse()
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		os.Exit(m.Run())
	}
	main()
	os.Exit(0)
}
