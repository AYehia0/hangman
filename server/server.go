package server

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/logging"
)

func publicKeyHandler(_ctx ssh.Context, _key ssh.PublicKey) bool {
	return true
}

// start the server!
func Start(host string, port int) {
	fmt.Println("Starting the server!!!")
	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", host, port)),
		wish.WithPublicKeyAuth(publicKeyHandler),
		wish.WithMiddleware(func(h ssh.Handler) ssh.Handler {
			return func(s ssh.Session) {
				wish.Println(s, "Hello, from the ssh server!")
				h(s)
			}
		},
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
