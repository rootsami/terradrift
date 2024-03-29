package config

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
)

// ConfigGenerator creates a stack for each directory that contains .tf files
func ConfigGenerator(workdir string) (*Config, error) {

	var stacks []Stack
	var stack Stack
	cfgPaths := findStack(workdir, ".tf")

	for _, path := range cfgPaths {

		// if the directory contains .tfvars files, create a stack for each .tfvars file
		if len(findStack(workdir+path, ".tfvars")) > 0 {

			for _, subStack := range findStack(workdir+path, ".tfvars") {
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

// findStack discover all directories in the workdir that contains files with .tf extension and returns a list of stacks
// if the extension is .tfvars, it returns a list of .tfvars files to determine how many stacks are in the same directory
func findStack(workdir, ext string) []string {

	var list []string
	filepath.WalkDir(workdir, func(s string, f fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(f.Name()) == ext && !strings.Contains(s, ".terraform") {

			switch ext {
			case ".tf":
				f := filepath.Dir(s)
				d := f[len(workdir):]

				// Append the directory path to the list and eliminate duplicates
				if !slices.Contains(list, d) {
					list = append(list, d)
				}
			case ".tfvars":
				d := s[len(workdir):]
				// set max depth of 1 in case of tfvars are in a subdirectory
				// which create an issue by duplicating stacks for tfvars of a subdirectory that doesn't belong to it
				if strings.Count(d, string(os.PathSeparator)) > 2 {
					return fs.SkipDir
				}
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
