package main

import (
	"io"
	"net"
	"strings"
	"time"
)

func setupDaemon(addr, password, user string) {
	for {
		conn, ipType, err := setupConnection(addr, password, user)

		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		setUpTimer(conn, ipType)
	}
}

func setupConnection(addr, password, user string) (net.Conn, int, error) {
	logf("Connecting...")

	ipType := 4
	conn, err := net.DialTimeout("tcp", addr, 10*time.Second)

	if err != nil {
		logger.Printf("Can not connection server: %s\n", err.Error())
		destroyConn(conn)
		return nil, 0, err
	}

	buf := make([]byte, 1024)
	_, err = conn.Read(buf)

	if err != nil && err != io.EOF {
		logger.Printf("Read server data faild: %s\n", err.Error())
		destroyConn(conn)
		return nil, 0, err
	}

	logf("Response: %s\n", string(buf))

	if strings.Contains(string(buf), "Authentication required") {
		_, err := conn.Write([]byte(user + ":" + password + "\n"))

		if err != nil {
			logger.Printf("Send authentication required faild: %s\n", err.Error())
			destroyConn(conn)
			return nil, 0, err
		}

		buf := make([]byte, 1024)
		_, err = conn.Read(buf)

		if err != nil && err != io.EOF {
			logger.Printf("Read authentication required faild: %s\n", err.Error())
			destroyConn(conn)
			return nil, 0, err
		}

		res := string(buf)

		if !strings.Contains(res, "Authentication successful") {
			logger.Printf("Authentication required faild: %s\n", res)
			destroyConn(conn)
			return nil, 0, err
		}

		logf("%s\n", res)

		if strings.Contains(res, "IPv4") {
			ipType = 4
		} else {
			ipType = 6
		}
	}

	return conn, ipType, nil
}

func destroyConn(conn net.Conn) {
	logf("Disconnecting...")

	err := conn.Close()

	if err != nil {
		logger.Printf("Close connection error: %s\n", err.Error())
	}
}
