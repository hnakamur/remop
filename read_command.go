package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
)

func doRunReadSubcommand(socketPath, itemURL string) error {
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return fmt.Errorf("connect to %s: %w", socketPath, err)
	}
	defer conn.Close()

	slog.Debug("read: connected")
	_, err = fmt.Fprintf(conn, "%s\n", itemURL)
	if err != nil {
		return nil
	}
	slog.Debug("read: written message")
	unixConn, ok := conn.(*net.UnixConn)
	if !ok {
		return errors.New("not UnixConn error")
	}
	if err := unixConn.CloseWrite(); err != nil {
		return err
	}
	slog.Debug("read: closed write")

	dec := json.NewDecoder(conn)
	var res OpReadResponse
	if err := dec.Decode(&res); err != nil {
		return err
	}
	if res.Stderr != "" {
		if _, err := os.Stderr.WriteString(res.Stderr); err != nil {
			return err
		}
	}
	if _, err := os.Stdout.WriteString(res.Stdout); err != nil {
		return err
	}

	return nil
}
