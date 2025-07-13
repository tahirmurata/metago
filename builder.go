// Copyright 2025 endmin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package builder handles the static site generation.
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func Build(cfg *Config, layoutPath string, outDir string) error {
	// Clean output directory first
	err := os.RemoveAll(outDir)
	if err != nil {
		return fmt.Errorf("remove output dir: %w", err)
	}

	// Create output directory
	err = os.Mkdir(outDir, 0755)
	if err != nil {
		return fmt.Errorf("make output dir: %w", err)
	}

	// Generate index.html in folder for each repository
	for _, repo := range cfg.Repository {
		templateData := Data{
			Domain: cfg.Domain,
			Path:   repo.Path,
			VCS:    repo.VCS,
			Repo:   repo.Repo,
		}

		outputPath := filepath.Join(outDir, fmt.Sprintf("%s.html", repo.Path))
		err = RenderAndWrite(templateData, layoutPath, outputPath)
		if err != nil {
			return fmt.Errorf("render and write static site: %w", err)
		}
	}

	return nil
}
