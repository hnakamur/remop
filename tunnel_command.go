package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"slices"
	"time"
)

type OpReadResponse struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

func runTunnelServer(config *Config, targetHost string) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	remoteSocketPath := config.Remote.SocketPath
	sshPath := config.Local.SshPath
	rmCmd := exec.Command(sshPath, targetHost, "rm -f "+remoteSocketPath)
	if err := rmCmd.Run(); err != nil {
		return fmt.Errorf("remove remote socket, err: %s", err)
	}

	localSocketPath := config.Local.SocketPath
	cmd := exec.CommandContext(ctx, sshPath, "-N", "-R",
		remoteSocketPath+":"+localSocketPath, targetHost)
	err := cmd.Start()
	if err != nil {
		return err
	}

	ln, err := listenUnix(localSocketPath)
	if err != nil {
		return err
	}
	defer ln.Close()

	slog.Info("tunnel start")
	defer slog.Info("tunnel exit")

outer:
	for {
		acceptDeadline := time.Now().Add(time.Second)
		if err := ln.SetDeadline(acceptDeadline); err != nil {
			return err
		}
		conn, err := ln.Accept()
		if err != nil {
			if !errors.Is(err, os.ErrDeadlineExceeded) {
				return err
			}
			select {
			case <-ctx.Done():
				break outer
			default:
				continue outer
			}
		}

		go func(c net.Conn) {
			defer c.Close()

			reader := bufio.NewReader(conn)
			line, err := reader.ReadBytes('\n')
			if err != nil {
				slog.Error("failed to read input item URL", "err", err)
				return
			}

			itemURL := string(bytes.TrimRight(line, "\n"))
			if !slices.Contains(config.Local.AllowedItemURLs, itemURL) {
				slog.Warn("denied request which is out of allowed list", "itemURL", itemURL)
			}

			stdout, stderr, err := runOpRead(config.Local.OpPath, itemURL)
			if err != nil && stderr == "" {
				slog.Error("remote: op read failed", "itemURL", itemURL, "err", err)
			}
			res := OpReadResponse{Stdout: stdout, Stderr: stderr}
			resBytes, err := json.Marshal(res)
			if err != nil {
				slog.Error("remote: marshal response", "err", err)
				return
			}
			if _, err := c.Write(resBytes); err != nil {
				slog.Error("remote: write response", "err", err)
				return
			}
			if stderr != "" {
				slog.Warn("sent error response", "itemURL", itemURL, "stderr", stderr)
			} else {
				slog.Info("sent success response", "itemURL", itemURL)
			}
		}(conn)
	}

	if err := cmd.Wait(); err != nil {
		if ee, ok := errors.AsType[*exec.ExitError](err); ok {
			if ee.ExitCode() != 255 {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}
