package main

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"time"

	"github.com/sedarasecurity/ossec-config-schema"
	"github.com/sedarasecurity/ossecsvc"
	"github.com/zerklabs/auburn/log"
)

var (
	manager        string
	port           string
	agentname      string
	listen         bool
	restartService bool

	configfile string
	clientkeys string
)

func init() {
	flag.StringVar(&manager, "manager", "", "Manager IP Address")
	flag.BoolVar(&listen, "listen", false, "Enables running in server mode")
	flag.BoolVar(&restartService, "controlsvc", false, "Enable or disable ossec-agent service control")
	flag.StringVar(&port, "port", "1515", "Manager port")

	h, err := os.Hostname()
	if err != nil {
		log.Error(err)
	}

	if h == "" {
		h = "localhost"
	}

	flag.StringVar(&agentname, "name", h, "Agent name")
	flag.StringVar(&configfile, "config", "", "Path to OSSEC config file (ossec.conf)")
	flag.StringVar(&clientkeys, "keyfile", "", "Path to OSSEC client keys file (client.keys)")
}

func main() {
	flag.Parse()

	// ossec.conf
	configFile := getOssecConfPath(configfile)
	if configFile == "" {
		log.Errorf("ossec.conf not found. Please specify OSSEC install directory")
		return
	}

	// client.keys
	keysFile := getClientKeysPath(clientkeys)
	if keysFile == "" {
		log.Infof("client.keys not found, creating an empty file")
		createDefaultClientKeys()
	}

	key, err := register(manager, port)
	if err != nil {
		log.Error(err)
		return
	}

	isclient := false

	if !listen {
		isclient = true
	}

	if err := appendkey(keysFile, key, isclient); err != nil {
		log.Error(err)
		return
	}

	if err := writeconfig(configFile, manager); err != nil {
		log.Error(err)
		return
	}

	if restartService {
		if err := ossecsvc.Stop(); err != nil {
			log.Error(err)
		} else {
			log.Infof("Service Stopped")
		}

		if err := ossecsvc.Start(); err != nil {
			log.Error(err)
		} else {
			log.Infof("Service Started")
		}
	}
}

func register(host string, port string) (string, error) {
	// see https://github.com/ossec/ossec-hids/blob/master/src/os_auth/main-server.c#L380 for buffer
	// allocated from ossec-auth server
	buf := make([]byte, 2048)

	// TODO(rch): this shouldn't ignore the certificate
	tc := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", net.JoinHostPort(host, port), tc)
	if err != nil {
		return "", fmt.Errorf("Error connecting to %s:%d. %v", host, port, err)
	}

	defer conn.Close()

	// timeout reading from the connection if we don't hear a response in 60 seconds
	err = conn.SetReadDeadline(time.Now().Add(time.Second * 60))
	if err != nil {
		return "", fmt.Errorf("timeout waiting for response from server: %v", err)
	}

	// send the request
	n, err := conn.Write([]byte(fmt.Sprintf("OSSEC A:'%s'\n", agentname)))
	if err != nil {
		return "", fmt.Errorf("Error writing to TCP connection: %v", err)
	}

	log.Debugf("Wrote %d bytes to TCP connection", n)

	n, err = conn.Read(buf)
	if err != nil {
		return "", fmt.Errorf("Error reading from TCP connection: %v", err)
	}

	log.Debugf("Read %d bytes to TCP connection", n)

	buffer := bytes.NewBuffer(nil)

	for _, v := range buf {
		if v != byte(0) {
			buffer.WriteByte(v)
		}
	}

	key := buffer.String()
	key = key[9:len(key)-3] + "\n"

	log.Debugf("Received response from ossec-authd: %s", key)
	return key, nil
}

func appendkey(file, line string, isclient bool) error {
	var fh *os.File
	var err error

	if file == "" {
		return fmt.Errorf("client.keys not found")
	}

	if isclient {
		fh, err = os.OpenFile(file, os.O_CREATE|os.O_RDWR|os.O_SYNC, 0700)
		if err != nil {
			fh.Close()
			if os.IsNotExist(err) {
				createDefaultClientKeys()
			} else {
				return fmt.Errorf("error appending key to client.keys: %v", err)
			}
		}
		defer fh.Close()
	} else {
		fh, err = os.OpenFile(file, os.O_APPEND|os.O_WRONLY|os.O_SYNC, 0700)
		if err != nil {
			fh.Close()
			return err
		}
		defer fh.Close()
	}

	log.Infof("Writing key to %s", file)

	if _, err := fh.Write([]byte(line)); err != nil {
		return fmt.Errorf("Error writing to client.keys: %v", err)
	}

	log.Infof("Key appended successfully")

	return nil
}

func readconfig(file string) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)
	fh, err := os.Open(file)
	if err != nil {
		fh.Close()
		return buf, fmt.Errorf("error reading config at %s: %v", file, err)
	}
	defer fh.Close()
	buf.ReadFrom(fh)

	return buf, nil
}

func writeconfig(file, ip string) error {
	var conf ossecconf.Ossec_config
	var nbuf bytes.Buffer

	buf, err := readconfig(file)
	if err != nil {
		return err
	}

	if err := xml.Unmarshal(buf.Bytes(), &conf); err != nil {
		return err
	}

	conf.ServerIp = ip

	fh, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR|os.O_SYNC, 0700)
	if err != nil {
		fh.Close()
		return err
	}
	defer fh.Close()

	enc := xml.NewEncoder(&nbuf)
	enc.Indent("", "  ")

	if err := enc.Encode(&conf); err != nil {
		return err
	}

	if err := ioutil.WriteFile(file, nbuf.Bytes(), 0700); err != nil {
		return err
	}

	return nil
}
