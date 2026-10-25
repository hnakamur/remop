package main

import (
	"errors"
	"fmt"
	"io/fs"
	"net"
	"os"
)

// listenUnix removes a leftover socket file at path (if nothing is
// listening on it) and listens on a new unix socket there.
func listenUnix(path string) (*net.UnixListener, error) {
	if isActiveSocket(path) {
		return nil, fmt.Errorf("unix socket %s is already in use", path)
	}
	if err := os.Remove(path); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("remove %s: %w", path, err)
	}
	ln, err := net.Listen("unix", path)
	if err != nil {
		return nil, fmt.Errorf("listen on %s: %w", path, err)
	}
	unixLn, ok := ln.(*net.UnixListener)
	if !ok {
		return nil, errors.New("not UnixListener error")
	}
	unixLn.SetUnlinkOnClose(true)
	return unixLn, nil
}

// isActiveSocket reports whether something is listening on the unix socket
// at path.
func isActiveSocket(path string) bool {
	conn, err := net.Dial("unix", path)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

// removeSocketFile removes the socket file at path, ignoring the case where
// the file does not exist.
func removeSocketFile(path string) {
	if err := os.Remove(path); err != nil && !errors.Is(err, fs.ErrNotExist) {
		fmt.Fprintf(os.Stderr, "op-remote: remove %s: %v\n", path, err)
	}
}
