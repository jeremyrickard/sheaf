/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package fs

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bryanl/sheaf/pkg/reporter"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

func TestManifestService_List(t *testing.T) {
	withBundleDir(t, func(dir string) {
		stageFile(t, sheaf.BundleConfigFilename, filepath.Join(dir, sheaf.BundleConfigFilename))

		manifestDir := filepath.Join("testdata", "manifests")
		m, err := NewManifestService(manifestDir, ManifestServiceReporter(reporter.Nop{}))
		require.NoError(t, err)

		actual, err := m.List()
		require.NoError(t, err)

		wanted := []sheaf.BundleManifest{
			{
				ID:   filepath.Join(manifestDir, "deploy.yaml"),
				Data: slurpData(t, filepath.Join(manifestDir, "deploy.yaml")),
			},
			{
				ID:   filepath.Join(manifestDir, "service.yaml"),
				Data: slurpData(t, filepath.Join(manifestDir, "service.yaml")),
			},
		}
		require.Equal(t, wanted, actual)
	})
}

func TestManifestService_Add_from_http_url(t *testing.T) {
	withBundleDir(t, func(bundleDir string) {
		stageFile(t, sheaf.BundleConfigFilename, filepath.Join(bundleDir, sheaf.BundleConfigFilename))

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "data")
		}))

		defer ts.Close()

		manifestDir := filepath.Join(bundleDir, "app", "manifests")
		m, err := NewManifestService(manifestDir, ManifestServiceReporter(reporter.Nop{}))
		require.NoError(t, err)

		err = m.Add(ts.URL + "/deploy.yaml")
		require.NoError(t, err)

		wantedPaths := []string{"deploy.yaml"}
		for _, p := range wantedPaths {
			_, err = os.Stat(filepath.Join(manifestDir, p))
			require.NoError(t, err)
		}
	})
}

func TestManifestService_Add_from_non_http_url(t *testing.T) {
	withBundleDir(t, func(bundleDir string) {
		stageFile(t, sheaf.BundleConfigFilename, filepath.Join(bundleDir, sheaf.BundleConfigFilename))

		manifestDir := filepath.Join(bundleDir, "app", "manifests")
		m, err := NewManifestService(manifestDir, ManifestServiceReporter(reporter.Nop{}))
		require.NoError(t, err)

		err = m.Add("ws://example.com/deploy.yaml")
		require.Error(t, err)
	})
}

func TestManifestService_Add_from_fs(t *testing.T) {
	cases := []struct {
		name        string
		manifestURI string
		setup       func(t *testing.T, bundleDir string)
		wantedPaths []string
		wantErr     bool
	}{
		{
			name:        "add file",
			manifestURI: filepath.Join("testdata", "manifests", "deploy.yaml"),
			wantedPaths: []string{"deploy.yaml"},
		},
		{
			name:        "add file (already exists)",
			manifestURI: filepath.Join("testdata", "manifests", "deploy.yaml"),
			setup: func(t *testing.T, bundleDir string) {
				dir := genManifestDir(bundleDir)
				require.NoError(t, os.MkdirAll(dir, 0700))

				f, err := os.Create(filepath.Join(dir, "deploy.yaml"))
				require.NoError(t, err)
				defer require.NoError(t, f.Close())
			},
			wantErr: true,
		},
		{
			name:        "add from directory",
			manifestURI: filepath.Join("testdata", "manifests"),
			wantedPaths: []string{"deploy.yaml", "service.yaml"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			withBundleDir(t, func(bundleDir string) {
				stageFile(t, sheaf.BundleConfigFilename, filepath.Join(bundleDir, sheaf.BundleConfigFilename))

				if tc.setup != nil {
					tc.setup(t, bundleDir)
				}

				manifestDir := filepath.Join(bundleDir, "app", "manifests")
				m, err := NewManifestService(manifestDir, ManifestServiceReporter(reporter.Nop{}))
				require.NoError(t, err)

				err = m.Add(tc.manifestURI)
				if tc.wantErr {
					require.Error(t, err)
					return
				}
				require.NoError(t, err)

				for _, p := range tc.wantedPaths {
					_, err = os.Stat(filepath.Join(manifestDir, p))
					require.NoError(t, err)
				}
			})
		})
	}
}