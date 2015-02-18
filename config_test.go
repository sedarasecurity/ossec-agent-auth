package main

import "testing"

func TestGetClientKeysPath(t *testing.T) {
	res := getClientKeysPath("/tmp/")
	if res != "" {
		t.Errorf("expected res to be empty")
	}
}

func TestCreateDefaultClientKeys(t *testing.T) {
	err := createDefaultClientKeys()
	if err != nil {
		t.Error(err)
	}
}
