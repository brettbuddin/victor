package httpserver

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	*http.ServeMux
	mutex       sync.Mutex
	listener    net.Listener
	tls         bool
	tlsKeyFile  string
	tlsCertFile string
}

func New() *Server {
	return &Server{
		ServeMux: http.NewServeMux(),
	}
}

func (s *Server) SetTLS(keyFile, certFile string) {
	s.tls = true
	s.tlsKeyFile = keyFile
	s.tlsCertFile = certFile
}

func (s *Server) ListenURL() string {
	scheme := "http"
	if s.tls {
		scheme = "https"
	}
	if s.listener != nil {
		if addr, ok := s.listener.Addr().(*net.TCPAddr); ok {
			if addr.IP.IsUnspecified() {
				return fmt.Sprintf("%s://localhost:%d", scheme, addr.Port)
			}
			return fmt.Sprintf("%s://%s", scheme, s.listener.Addr())
		}
	}
	return ""
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Started %s \"%s\"", r.Method, r.RequestURI)

	w = &ResponseWriteTracker{ResponseWriter: w}
	s.ServeMux.ServeHTTP(w, r)

	wt := w.(*ResponseWriteTracker)
	log.Printf("Completed %s %d (%d bytes written)", r.Method, wt.code, wt.size)
}

func (s *Server) Listen(addr string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.listener != nil {
		return nil
	}

	var err error
	if s.listener, err = net.Listen("tcp", addr); err != nil {
		return fmt.Errorf("could not listen on %s: %s", addr, err)
	}

	if s.tls {
		config := &tls.Config{
			Rand:       rand.Reader,
			Time:       time.Now,
			NextProtos: []string{"HTTP/1.1"},
		}
		config.Certificates = make([]tls.Certificate, 1)
		config.Certificates[0], err = tls.LoadX509KeyPair(s.tlsCertFile, s.tlsKeyFile)
		if err != nil {
			return fmt.Errorf("could not load TLS cert: %s", err)
		}
		s.listener = tls.NewListener(s.listener, config)
	}

	return nil
}

func (s *Server) Serve() error {
	if err := http.Serve(s.listener, s); err != nil {
		return fmt.Errorf("could not start server: %s", err)
	}
	return nil
}

func (s *Server) ListenAndServe(addr string) error {
	if err := s.Listen(addr); err != nil {
		return err
	}

	return s.Serve()
}

func (s *Server) Stop() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.listener.Close()
}
