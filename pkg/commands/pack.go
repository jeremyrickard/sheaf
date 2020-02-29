/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/bryanl/sheaf/pkg/bundle"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

// NewPackCommand creates a pack command.
func NewPackCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pack",
		Short: "pack a bundle",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("requires path to bundle directory")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			config := sheaf.PackConfig{
				BundleURI:     args[0],
				BundleFactory: bundle.DefaultBundleFactory,
				Packer:        bundle.NewPacker(os.Stdout),
			}
			return sheaf.Pack(config)
		},
	}

	return cmd
}
