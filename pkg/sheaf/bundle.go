/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

import (
	"io"

	"github.com/bryanl/sheaf/pkg/images"
)

//go:generate mockgen -destination=../mocks/mock_bundle.go -package mocks github.com/bryanl/sheaf/pkg/sheaf Bundle
//go:generate mockgen -destination=../mocks/mock_manifest_service.go -package mocks github.com/bryanl/sheaf/pkg/sheaf ManifestService

// ManifestGenerator generates manifests.
type ManifestGenerator interface {
	Show(w io.Writer) error
}

// BundleFactory is a factory for creating bundles given a URI.
type BundleFactory func(uri string) (Bundle, error)

// Bundle manages bundles.
type Bundle interface {
	Codec() Codec
	Path() string
	Config() BundleConfig
	Artifacts() ArtifactsService
	Manifests() (ManifestService, error)
	Images() (images.Set, error)
}

// ManifestService is a service for interacting with manifests.
type ManifestService interface {
	List() ([]BundleManifest, error)
	Add(manifestURIs ...string) error
}

// BundleManifest describes a manifest in a fs.
type BundleManifest struct {
	ID   string
	Data []byte
}
