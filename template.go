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

// Package template handles HTML template rendering.
package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
)

type Data struct {
	Domain   string
	Path     string
	VCS      string
	Repo     string
	MetaTags template.HTML
}

const defaultMetaTemplate = `<title>{{.Path}}</title>
		<meta name="go-import" content="{{.Domain}}/{{.Path}} {{.VCS}} {{.Repo}}" />
		<meta name="generator" content="aethergate v1.0.0" />`

const defaultLayoutTemplate = `<!doctype html>
<html>
    <head>
        {{.MetaTags}}
        <meta http-equiv="refresh" content="0;url={{.Repo}}" />
        <meta name="robots" content="noindex,noarchive" />
        <style>
            html {
                background-color: oklch(98.5% 0 0);
                color: oklch(14.1% 0.005 285.823);
                transition: background-color 0.3s ease;
            }

            body {
            	font-family: ui-sans-serif, system-ui, sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol', 'Noto Color Emoji';
            }

            @media (prefers-color-scheme: dark) {
                html {
                    background-color: oklch(21% 0.006 285.885);
                    color: oklch(98.5% 0 0);
                }
            }

            .centered-text {
                text-align: center;
                margin-top: calc(.24rem * 6);
            }

            .repo-link {
                color: inherit;
                text-decoration: underline;
            }
        </style>
    </head>
    <body>
        <p class="centered-text">
            Redirecting to <a href="{{.Repo}}" class="repo-link">repository</a>...
        </p>
    </body>
</html>`

func RenderAndWrite(data Data, layoutPath string, outputPath string) error {
	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("text/css", css.Minify)

	metaTmpl, err := template.New("goMeta").Parse(defaultMetaTemplate)
	if err != nil {
		return fmt.Errorf("parse meta template: %w", err)
	}

	layoutTmpl, err := template.ParseFiles(layoutPath)
	if err != nil {
		return fmt.Errorf("parse layout template: %w", err)
	}

	var metaBuffer bytes.Buffer
	err = metaTmpl.Execute(&metaBuffer, data)
	if err != nil {
		return fmt.Errorf("execute meta template: %w", err)
	}

	// Add the meta HTML to the data
	data.MetaTags = template.HTML(metaBuffer.String())

	// Create the output directory if it doesn't exist
	err = os.MkdirAll(filepath.Dir(outputPath), 0755)
	if err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	// Execute the layout template into new buffer
	var finalBuffer bytes.Buffer
	err = layoutTmpl.Execute(&finalBuffer, data)
	if err != nil {
		return fmt.Errorf("execute layout template: %w", err)
	}

	// Create the output file
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}
	defer f.Close()

	// Minify the HTML output before writing to file
	err = m.Minify("text/html", f, &finalBuffer)
	if err != nil {
		return fmt.Errorf("minify HTML output: %w", err)
	}

	return f.Sync()
}

func WriteDefaultLayout(path string) error {
	err := os.WriteFile(path, []byte(defaultLayoutTemplate), 0644)
	if err != nil {
		return fmt.Errorf("write default layout: %w", err)
	}
	return nil
}
