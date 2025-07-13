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

/*
aethergate is a tool that generates static files to support vanity Go remote
import paths.

It allows you to create vanity remote import paths like `import example.com/pkg`
while hosting your code on platforms such as GitHub, GitLab, or SourceHut. Each
repository defined in the `aethergate.toml` configuration file generates an
`index.html` file containing `go-import` meta tags, which redirect the go tool
to the actual repository.

Usage:

	aethergate [command]

Commands:

	init
		Initialize a new config and layout file

	build <dir>
		Build static site

Customization:

	aethergate.toml
		The configuration file. Use it to define your vanity imports.

	layout.html
		The HTML template for generating pages. Available template actions:

		{{.MetaTags}}: Outputs the required meta tags for Go imports.
		{{.Domain}}: The domain of the import path.
		{{.Path}}: The path component of the import.
		{{.VCS}}: The version control system (e.g., git, hg).
		{{.Repo}}: The repository URL.
*/
package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	config = "aethergate.toml"
	layout = "layout.html"
)

func main() {
	initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	initCmd.Usage = func() {
		fmt.Fprintln(initCmd.Output(), "Usage: aethergate init")
	}
	buildCmd := flag.NewFlagSet("build", flag.ExitOnError)
	buildCmd.Usage = func() {
		fmt.Fprintln(buildCmd.Output(), "Usage: aethergate build <dir>")
	}

	if len(os.Args) <= 1 {
		fmt.Println("error: no command specified.")
		fmt.Println("Try 'aethergate -h' for more information.")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		initCmd.Parse(os.Args[2:])
		cfg := DefaultConfig()

		if _, err := os.Stat(config); os.IsNotExist(err) {
			WriteDefault(cfg, config)
			fmt.Println("Created config file")
		} else if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		} else {
			fmt.Println("Config file already exists")
		}

		if _, err := os.Stat(layout); os.IsNotExist(err) {
			err := WriteDefaultLayout(layout)
			if err != nil {
				fmt.Println("error:", err)
				os.Exit(1)
			}
			fmt.Println("Created layout file")
		} else if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		} else {
			fmt.Println("Layout file already exists")
		}
	case "build":
		buildCmd.Parse(os.Args[2:])
		buildDir := "dist"
		if buildCmd.Arg(0) != "" {
			buildDir = buildCmd.Arg(0)
		}
		cfg, err := Load(config)
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}

		err = Build(cfg, layout, buildDir)
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}

		fmt.Println("Site built successfully")
	case "-h", "--h", "-help", "--help":
		fmt.Println(`Usage: aethergate [command]
Commands:
  init
        Initialize a new config and layout file
  build <dir>
        Build static site`)
	default:
		fmt.Println("error: no command specified.")
		fmt.Println("Try 'aethergate -h' for more information.")
		os.Exit(1)
	}
}
