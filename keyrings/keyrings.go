// Copyright 2025 cloudeng llc. All rights reserved.
// Use of this source code is governed by the Apache-2.0
// license that can be found in the LICENSE file.

//go:build linux

// Package keyrings provides a file based API for working with
// the linux kernel keyrings facility. It is primarily intended
// for injecting and accssing secrets inside a linux docker
// container without having to use docker swarm or other
// container orchestration tools.
// A typical use case would be to pipe in cloud service
// credentials to a container when it is started and for those
// credentials to be used to access other services, such as
// a cloud secrets manager etc. It is not intended for
// user centric, interactive use, as per the gnome or macos
// keychain/keyring facilities.
package keyrings

import (
	"context"
	"io/fs"
	"sync"
	"syscall"

	"github.com/cloudengio/keyctl"
)

// Option represents an option for configuring a keyrings.T
type Option func(o *options)

// WithKeyring specifies the keyring to use for storing and retrieving secrets.
func WithKeyring(keyring keyctl.Keyring) Option {
	return func(o *options) {
		o.keyring = keyring
	}
}

type options struct {
	keyring keyctl.Keyring
}

// T provides access to linux kernel keyrings.
type T struct {
	mu   sync.Mutex
	opts options
	err  error
}

// New creates a new keyrings.T. If WithKeyring is not specified then the
// session keyring is used.
func New(opts ...Option) *T {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}
	kr := &T{opts: o}
	if o.keyring == nil {
		kr.opts.keyring, kr.err = keyctl.SessionKeyring()
	}
	return kr
}

func (s *T) error() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.err
}

func mapError(err error) error {
	if err == nil {
		return nil
	}
	errno, ok := err.(syscall.Errno)
	if !ok {
		return err
	}
	if errno == syscall.ENOKEY {
		return fs.ErrNotExist
	}
	return err
}

func (s *T) get(name string) (*keyctl.Key, error) {
	if err := s.error(); err != nil {
		return nil, err
	}
	key, err := s.opts.keyring.Search(name)
	return key, mapError(err)
}

func (s *T) ReadFileCtx(ctx context.Context, name string) ([]byte, error) {
	if err := s.error(); err != nil {
		return nil, err
	}
	key, err := s.get(name)
	if err != nil {
		return nil, err
	}
	return key.Get()
}

func (s *T) WriteFileCtx(ctx context.Context, name string, data []byte) error {
	if err := s.error(); err != nil {
		return err
	}
	_, err := s.opts.keyring.Add(name, data)
	return err
}

func (s *T) Delete(ctx context.Context, name string) error {
	if err := s.error(); err != nil {
		return err
	}
	key, err := s.get(name)
	if err != nil {
		return err
	}
	return key.Unlink()
}
