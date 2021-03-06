/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package archive

import (
	"fmt"

	"github.com/spf13/cobra"

	archive2 "github.com/bryanl/sheaf/pkg/archive"
	"github.com/bryanl/sheaf/pkg/fs"
)

// NewStageCommand creates a stage command.
func NewStageCommand() *cobra.Command {
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "relocate",
		Short: "Relocate images in archive to new registry",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("requires fs location and registry prefix")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			relocator := fs.NewImageRelocator(
				fs.ImageRelocatorDryRun(dryRun))

			stager := archive2.NewStager(
				archive2.StagerOptionImageRelocator(relocator))

			return stager.Stage(args[0], args[1])
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "dry run")

	return cmd
}
