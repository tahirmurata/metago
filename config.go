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

// Package config handles the aethergate.toml configuration file.
package main

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Repository struct {
	Path string `toml:"path" comment:"Path of vanity remote import"`
	VCS  string `toml:"vcs" comment:"Version control system"`
	Repo string `toml:"repo" comment:"Repository URL'"`
}

type Config struct {
	Version int    `toml:"version" comment:"Version of aethergate"`
	Domain  string `toml:"domain" comment:"Domain of vanity remote import'"`

	Repository []Repository `toml:"repository" comment:"List of repositories to serve"`
}

func DefaultConfig() *Config {
	return &Config{
		Version: 1,
		Domain:  "go.endfieldind.com",
		Repository: []Repository{
			{
				Path: "aethergate",
				VCS:  "git",
				Repo: "https://git.sr.ht/~endmin/aethergate",
			},
		},
	}
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var cfg Config
	err = toml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("unmarshal data: %w", err)
	}

	return &cfg, nil
}

func WriteDefault(cfg *Config, path string) error {
	data, err := toml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}
