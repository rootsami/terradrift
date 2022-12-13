package config

import (
	"io/fs"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
)

// ConfigGenerator creates a stack for each directory that contains .tf files
func ConfigGenerator(workspace string) (*Config, error) {

	var stacks []Stack
	var stack Stack
	cfgPaths := findStack(workspace, ".tf")

	for _, path := range cfgPaths {

		// if the directory contains .tfvars files, create a stack for each .tfvars file
		if len(findStack(workspace+path, ".tfvars")) > 0 {

			for _, subStack := range findStack(workspace+path, ".tfvars") {
				// remove the .tfvars extension and replace / with - to create a stack name
				name := strings.TrimSuffix(strings.ReplaceAll(path+"-"+subStack, "/", "-"), ".tfvars")
				stack = Stack{
					Name:    name,
					Path:    path,
					TFvars:  subStack,
					Backend: strings.ReplaceAll(subStack, "tfvars", "hcl"),
				}
				stacks = append(stacks, stack)
			}

		} else {

			stack = Stack{
				Name: strings.ReplaceAll(path, "/", "-"),
				Path: path,
			}

			stacks = append(stacks, stack)
		}

	}

	cfg := Config{
		Stacks: stacks,
	}

	return &cfg, nil
}

// findStack discover all directories in the workspace that contains files with .tf extension and returns a list of stacks
// if the extension is .tfvars, it returns a list of .tfvars files to determine how many stacks are in the same directory
func findStack(workspace, ext string) []string {

	var list []string
	filepath.WalkDir(workspace, func(s string, f fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(f.Name()) == ext && !strings.Contains(s, ".terraform") {

			switch ext {
			case ".tf":
				f := filepath.Dir(s)
				d := f[len(workspace):]

				// Append the directory path to the list and eliminate duplicates
				if !slices.Contains(list, d) {
					list = append(list, d)
				}
			case ".tfvars":
				d := s[len(workspace):]
				d = strings.Trim(d, "/")

				if !slices.Contains(list, d) {
					list = append(list, d)
				}
			}
		}

		return nil
	})

	return list
}
