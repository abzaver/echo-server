package main

import (
	"testing"
)

// Tests for the first example using the Proxy pattern (nginx and application)

func TestNginxProxy_RequestAllowed(t *testing.T) {
	nginxServer := newNginxServer()

	httpCode, body := nginxServer.handleRequest("/app/status", "GET")
	if httpCode != 200 || body != "Ok" {
		t.Errorf("Expected 200 OK, got %d %s", httpCode, body)
	}
}

func TestNginxProxy_RequestBlocked(t *testing.T) {
	nginxServer := newNginxServer()

	// Make 3 requests to exceed the limit
	nginxServer.handleRequest("/app/status", "GET")
	nginxServer.handleRequest("/app/status", "GET")
	httpCode, body := nginxServer.handleRequest("/app/status", "GET")

	if httpCode != 403 || body != "Not Allowed" {
		t.Errorf("Expected 403 Not Allowed, got %d %s", httpCode, body)
	}
}

func TestNginxProxy_InvalidEndpoint(t *testing.T) {
	nginxServer := newNginxServer()

	httpCode, body := nginxServer.handleRequest("/invalid/url", "GET")
	if httpCode != 404 || body != "Not Ok" {
		t.Errorf("Expected 404 Not Ok, got %d %s", httpCode, body)
	}
}

func TestNginxProxy_CreateUser(t *testing.T) {
	nginxServer := newNginxServer()

	// Test the creation of a user
	httpCode, body := nginxServer.handleRequest("/create/user", "POST")
	if httpCode != 201 || body != "User Created" {
		t.Errorf("Expected 201 User Created, got %d %s", httpCode, body)
	}
}
