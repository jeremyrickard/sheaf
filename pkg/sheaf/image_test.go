/*
 * Copyright 2020 Sheaf Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package sheaf

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadFromIndex(t *testing.T) {
	tests := []struct {
		name     string
		index    string
		wantErr  bool
		expected []Image
	}{
		{
			name:  "valid archive",
			index: "index.json",
			expected: []Image{
				{
					MediaType: "application/vnd.docker.distribution.manifest.list.v2+json",
					Size:      1412,
					Digest:    "sha256:ad5552c786f128e389a0263104ae39f3d3c7895579d45ae716f528185b36bc6f",
					Annotations: map[string]string{
						"org.opencontainers.image.ref.name": "docker.io/library/nginx:1.17",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index := filepath.Join("testdata", tt.index)

			got, err := LoadFromIndex(index)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expected, got)
		})
	}
}
