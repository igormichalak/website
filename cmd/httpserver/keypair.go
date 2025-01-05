package main

import (
	"crypto/tls"
	"fmt"
	"sync"

	"github.com/fsnotify/fsnotify"
)

type KeyPairReloader struct {
	CertFile    string
	KeyFile     string
	Certificate *tls.Certificate
	mu          sync.RWMutex
}

func (app *application) StartKeyPairReloader(certFile, keyFile string) (*KeyPairReloader, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	kpr := &KeyPairReloader{
		CertFile:    certFile,
		KeyFile:     keyFile,
		Certificate: &cert,
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create key pair watcher: %w", err)
	}
	go func() {
		defer watcher.Close()
		for {
			select {
			case ev, ok := <-watcher.Events:
				if !ok {
					return
				}
				if ev.Has(fsnotify.Write) {
					kpr.Reload()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				app.Logger.Error("key pair watcher error", "err", err)
			}
		}
	}()
	if err := watcher.Add(certFile); err != nil {
		return nil, err
	}
	if err := watcher.Add(keyFile); err != nil {
		return nil, err
	}
	return kpr, nil
}

func (kpr *KeyPairReloader) Reload() error {
	cert, err := tls.LoadX509KeyPair(kpr.CertFile, kpr.KeyFile)
	if err != nil {
		return err
	}
	kpr.mu.Lock()
	kpr.Certificate = &cert
	kpr.mu.Unlock()
	return nil
}

func (kpr *KeyPairReloader) GetCertificate(_ *tls.ClientHelloInfo) (*tls.Certificate, error) {
	kpr.mu.RLock()
	defer kpr.mu.RUnlock()
	return kpr.Certificate, nil
}
