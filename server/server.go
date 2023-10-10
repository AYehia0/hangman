package server

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AYehia0/hangman/game"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	wishtea "github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

func publicKeyHandler(_ctx ssh.Context, _key ssh.PublicKey) bool {
	return true
}

// start the server!
func Start(host string, port int) {
	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", host, port)),
		wish.WithPublicKeyAuth(publicKeyHandler),
		wish.WithMiddleware(
			wishtea.Middleware(gameHandler()),
			logging.Middleware(),
		))

	if err != nil {
		log.Error("Something went wrong : %s", "error", err)
	}

	// unix signal channel
	done := make(chan os.Signal, 1)

	// notify if one of these signal types occur
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)

	// listening for requests to the server
	go func() {
		if err := s.ListenAndServe(); err != nil && errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping the SSH Server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() {
		cancel()
	}()

	// closing the server
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("could not stop server", "error", err)
	}
}

func gameHandler() wishtea.Handler {
	return func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
		opts := []tea.ProgramOption{}
		pty, _, active := s.Pty()
		if !active {
			log.Print("No active terminal PTY, Skipping", "info")
			return nil, nil
		}

		rand.Seed(time.Now().UnixNano())

		// TODO: choose to use the api or the static json file

		g := game.Play(pty.Window.Width, pty.Window.Height, s)

		// full screen mode added to the options
		opts = append(opts, tea.WithAltScreen())

		return g, opts
	}
}
