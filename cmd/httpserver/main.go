package main

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const DefaultPort = "8080"
const BaseURL = "https://www.igormichalak.com"

type application struct {
	Logger        *slog.Logger
	Debug         bool
	TemplateCache map[string]*template.Template
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	certFile := flag.String("cert-file", "./tls/cert.pem", "certificate file path")
	keyFile := flag.String("key-file", "./tls/key.pem", "key file path")
	debug := flag.Bool("debug", false, "debug mode")
	flag.Parse()

	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = DefaultPort
	}
	redirectPort := os.Getenv("REDIRECT_PORT")

	var logger *slog.Logger
	if *debug {
		handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
		logger = slog.New(handler)
	} else {
		handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
		logger = slog.New(handler)
	}

	tmplCache, err := newTemplateCache()
	if err != nil {
		logger.Error("Failed to instantiate template cache", "err", err)
		return fmt.Errorf("failed to instantiate template cache: %w", err)
	}

	app := &application{
		Logger:        logger,
		Debug:         *debug,
		TemplateCache: tmplCache,
	}

	kpr, err := app.StartKeyPairReloader(*certFile, *keyFile)
	if err != nil {
		logger.Error("Failed to start key pair reloader", "err", err)
		return fmt.Errorf("failed to start key pair reloader: %w", err)
	}

	tlsConfig := &tls.Config{
		GetCertificate: kpr.GetCertificate,
		MinVersion:     tls.VersionTLS12,
		MaxVersion:     tls.VersionTLS13,
		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
		},
		CurvePreferences:   []tls.CurveID{tls.X25519, tls.CurveP256},
		ClientSessionCache: tls.NewLRUClientSessionCache(128),
	}

	errorLog := slog.NewLogLogger(logger.Handler(), slog.LevelError)

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           app.routes(),
		ReadTimeout:       6 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      12 * time.Second,
		IdleTimeout:       time.Minute,
		MaxHeaderBytes:    8_192,
		ErrorLog:          errorLog,
		TLSConfig:         tlsConfig,
	}

	redirectSrv := &http.Server{
		Addr:              ":" + redirectPort,
		Handler:           http.HandlerFunc(app.redirectToTLS),
		ReadTimeout:       6 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      6 * time.Second,
		IdleTimeout:       30 * time.Second,
		MaxHeaderBytes:    8_192,
		ErrorLog:          errorLog,
	}

	stopC := make(chan os.Signal, 1)
	signal.Notify(stopC, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	errorC := make(chan error, 1)

	go func(ec chan error) {
		logger.Info(fmt.Sprintf("Starting server on %s...", srv.Addr))

		err := srv.ListenAndServeTLS("", "")
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Server failed", "err", err)
			ec <- fmt.Errorf("server failed: %w", err)
		} else {
			ec <- nil
		}
	}(errorC)

	go func(ec chan error) {
		if redirectPort == "" {
			return
		}
		logger.Info(fmt.Sprintf("Starting redirect server on %s...", redirectSrv.Addr))

		err := redirectSrv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Redirect server failed", "err", err)
			ec <- fmt.Errorf("redirect server failed: %w", err)
		} else {
			ec <- nil
		}
	}(errorC)

	select {
	case err := <-errorC:
		return err
	case <-stopC:
	}

	logger.Info("Shutting down server(s)...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if redirectPort != "" {
		if err := redirectSrv.Shutdown(ctx); err != nil {
			logger.Error("Redirect server shutdown failed", "err", err)
			return fmt.Errorf("redirect server shutdown failed: %w", err)
		}
	}
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown failed", "err", err)
		return fmt.Errorf("server shutdown failed: %w", err)
	}
	logger.Info("Server(s) gracefully stopped.")

	return nil
}
