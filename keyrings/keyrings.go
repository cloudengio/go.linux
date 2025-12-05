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
	opts options
}

// New creates a new keyrings.T. If WithKeyring is not specified then the
// session keyring is used.
func New(opts ...Option) (*T, error) {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}
	kr := &T{opts: o}
	if o.keyring == nil {
		skr, err := keyctl.SessionKeyring()
		if err != nil {
			return nil, err
		}
		kr.opts.keyring = skr
	}
	return kr, nil
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
	key, err := s.opts.keyring.Search(name)
	return key, mapError(err)
}

// ReadFileCtx reads a secret from the kernel keyring.
func (s *T) ReadFileCtx(_ context.Context, name string) ([]byte, error) {
	key, err := s.get(name)
	if err != nil {
		return nil, err
	}
	return key.Get()
}

func (s *T) ReadFile(name string) ([]byte, error) {
	return s.ReadFileCtx(context.Background(), name)
}

// WriteFileCtx writes a secret to the kernel keyring. The fs.FileMode parameter is ignored.
func (s *T) WriteFileCtx(_ context.Context, name string, data []byte, _ fs.FileMode) error {
	_, err := s.opts.keyring.Add(name, data)
	return err
}

func (s *T) WriteFile(name string, data []byte, mode fs.FileMode) error {
	return s.WriteFileCtx(context.Background(), name, data, mode)
}

// Delete removes a secret from the kernel keyring.
func (s *T) Delete(_ context.Context, name string) error {
	key, err := s.get(name)
	if err != nil {
		return err
	}
	return key.Unlink()
}
