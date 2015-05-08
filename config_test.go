package main

import (
	"os"
	"testing"
)

func bootstrapTests() {
	os.Remove("/tmp/ossec.conf")
	os.Remove("/tmp/client.keys")

	os.Create("/tmp/ossec.conf")
	os.Create("/tmp/client.keys")
}

func moveFiles(start bool) {
	if start {
		os.Rename("/var/ossec/etc/client.keys", "/var/ossec/etc/client.keys.fortesting")
		os.Rename("/var/ossec/etc/ossec.conf", "/var/ossec/etc/ossec.conf.fortesting")
	} else {
		os.Rename("/var/ossec/etc/client.keys.fortesting", "/var/ossec/etc/client.keys")
		os.Rename("/var/ossec/etc/ossec.conf.fortesting", "/var/ossec/etc/ossec.conf")
	}
}

func TestGetClientKeysPath(t *testing.T) {
	bootstrapTests()

	res := getClientKeysPath("/tmp/client.keys")
	if res == "" {
		t.Errorf("res should not be empty")
	}

	if res != "/tmp/client.keys" {
		t.Errorf("res should have been /tmp/client.keys")
	}
}

func TestGetClientKeysPath_Default(t *testing.T) {
	res := getClientKeysPath("")
	if res == "" {
		t.Errorf("res should not be empty")
	}

	if res != "/var/ossec/etc/client.keys" {
		t.Errorf("res should have been /var/ossec/etc/client.keys, but was %s", res)
	}
}

func TestGetClientKeysPath_Missing(t *testing.T) {
	moveFiles(true)
	res := getClientKeysPath("")
	if res != "" {
		t.Errorf("res should be empty")
	}
	moveFiles(false)
}

func TestGetOssecConfPath(t *testing.T) {
	bootstrapTests()

	res := getOssecConfPath("/tmp/ossec.conf")
	if res == "" {
		t.Errorf("res should not be empty")
	}

	if res != "/tmp/ossec.conf" {
		t.Errorf("res should have been /tmp/ossec.conf")
	}
}

func TestGetOssecConfPath_Default(t *testing.T) {
	res := getOssecConfPath("")
	if res == "" {
		t.Errorf("res should not be empty")
	}

	if res != "/var/ossec/etc/ossec.conf" {
		t.Errorf("res should have been /var/ossec/etc/ossec.conf, but was %s", res)
	}
}

func TestGetOssecConfPath_Missing(t *testing.T) {
	moveFiles(true)
	res := getOssecConfPath("")
	if res != "" {
		t.Errorf("res should be empty")
	}
	moveFiles(false)
}

func TestCreateDefaultClientKeys(t *testing.T) {
	err := createDefaultClientKeys()
	if err != nil {
		t.Error(err)
	}
}
