package httpReload

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	defaultTimeout = 20 * time.Second      // default timeout is 20s
	graceEnv       = "go_http_reload=true" // env flag for reload
)

// A Grace carries actions for graceful restart or shutdown.
type Grace interface {
	Run(*http.Server) error
	ListenAndServe(string, http.Handler) error
}

type grace struct {
	srv      *http.Server
	listener net.Listener
	timeout  time.Duration
	err      error
}

func (g *grace) reload() *grace {
	f, err := g.listener.(*net.TCPListener).File()
	if err != nil {
		g.err = err
		return g
	}
	defer f.Close()

	var args []string
	if len(os.Args) > 1 {
		args = append(args, os.Args[1:]...)
	}
	cmd := exec.Command(os.Args[0], args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), graceEnv)
	cmd.ExtraFiles = []*os.File{f}

	g.err = cmd.Start()
	return g
}

func (g *grace) stop() *grace {
	if g.err != nil {
		return g
	}

	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()

	if err := g.srv.Shutdown(ctx); err != nil {
		g.err = err
	}

	return g
}

func (g *grace) run() (err error) {
	if _, ok := syscall.Getenv(strings.Split(graceEnv, "=")[0]); ok {
		f := os.NewFile(3, "")
		if g.listener, err = net.FileListener(f); err != nil {
			return
		}
	} else {
		if g.listener, err = net.Listen("tcp", g.srv.Addr); err != nil {
			return
		}
	}

	terminate := make(chan error)
	go func() {
		if err := g.srv.Serve(g.listener); err != nil {
			terminate <- err
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit)

	for {
		select {
		case s := <-quit:
			switch s {
			case syscall.SIGINT, syscall.SIGTERM:
				signal.Stop(quit)
				return g.stop().err

			case syscall.SIGUSR2:
				return g.reload().stop().err
			}

		case err = <-terminate:
			return
		}
	}
}

// WithTimeout returns a custom timeout Grace.
func WithTimeout(timeout time.Duration) Grace {
	return &grace{timeout: timeout}
}

func (g *grace) Run(srv *http.Server) error {
	g.srv = srv
	return g.run()
}

func (g *grace) ListenAndServe(addr string, handler http.Handler) error {
	g.srv = &http.Server{Addr: addr, Handler: handler}
	return g.run()
}

var _ Grace = (*grace)(nil) // assert *grace implements Grace.

// Run accepts a custom http Server and provice signal magic.
func Run(srv *http.Server) error {
	return WithTimeout(defaultTimeout).Run(srv)
}

// ListenAndServe wraps http.ListenAndServe and provides signal magic.
func ListenAndServe(addr string, handler http.Handler) error {
	return WithTimeout(defaultTimeout).ListenAndServe(addr, handler)
}
