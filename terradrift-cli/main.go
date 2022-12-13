package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/rootsami/terradrift/pkg/config"
	"github.com/rootsami/terradrift/pkg/stacks"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"

	log "github.com/sirupsen/logrus"
)

var (
	app              = kingpin.New("terradrift", "A tool to detect drifts in terraform IaC")
	workspace        = app.Flag("workspace", "workspace of a project that contains all terraform directories").Default("./").String()
	configPath       = app.Flag("config", "Path for configuration file holding the stack information").String()
	extraBackendVars = app.Flag("extra-backend-vars", "Extra backend environment variables ex. GOOGLE_CREDENTIALS OR AWS_ACCESS_KEY").StringMap()
	debug            = app.Flag("debug", "Enable debug mode").Default("false").Bool()
	generateConfig   = app.Flag("generate-config-only", "Generate a config file with based on a provided worksapce").Default("false").Bool()
	output           = app.Flag("output", "Output format supported: json, yaml and table").Default("table").Enum("table", "json", "yaml")
)

type stackOutput struct {
	Name    string `json:"name" yaml:"name"`
	Path    string `json:"path" yaml:"path"`
	Drift   bool   `json:"drift" yaml:"drift"`
	Add     int    `json:"add" yaml:"add"`
	Change  int    `json:"change" yaml:"change"`
	Destroy int    `json:"destroy" yaml:"destroy"`
	TFver   string `json:"tfver" yaml:"tfver"`
}

func main() {

	kingpin.MustParse(app.Parse(os.Args[1:]))

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	// if worksapce path is not absolute, make it absolute and add a trailing slash
	if !path.IsAbs(*workspace) {
		absPath, err := filepath.Abs(*workspace)
		if err != nil {
			log.Fatalf("Error getting absolute path for workspace: %s", err)
		}

		*workspace = absPath + "/"

	} else if !strings.HasSuffix(*workspace, "/") {
		*workspace = *workspace + "/"
	}

	var table [][]string
	var stackList []config.Stack
	var cfg *config.Config
	var stackOutputs []stackOutput

	// if config file is provided, load it and assign it to stackList
	if *configPath != "" {

		log.WithFields(log.Fields{"config": *configPath}).Debug("Loading config file")
		cfg, err := config.ConfigLoader(*workspace, *configPath)
		if err != nil {
			log.Fatalf("Error loading config file: %s", err)
		}
		stackList = cfg.Stacks

	} else {

		var err error
		// if config file is not provided, find all directories that contain .tf files
		// and create a stack for each directory
		log.Debug("Config file not found, running stack init on each directory that contains .tf files")
		cfg, err = config.ConfigGenerator(*workspace)
		if err != nil {
			log.Warnf("Error generating config file: %s", err)
		}

		stackList = cfg.Stacks

		// if generate-config-only flag is provided, print the generated config file and exit
		if *generateConfig {
			yamlConfig, err := yaml.Marshal(&cfg)
			if err != nil {
				log.Error("Error while marshaling. %v", err)
			}

			fmt.Print(string(yamlConfig))
			os.Exit(0)
		}

	}

	for _, stack := range stackList {

		// if configPath is not provided, create a config file with default values
		// create a config file with default values
		cfg := config.Stack{
			Name:    stack.Name,
			Path:    stack.Path,
			TFvars:  stack.TFvars,
			Backend: stack.Backend,
		}

		response, tfver, err := stacks.StackInit(*workspace, cfg, *extraBackendVars)
		if err != nil {
			log.Error(err)
		}

		stackOutputs = append(stackOutputs, stackOutput{
			Name:    cfg.Name,
			Path:    cfg.Path,
			Drift:   response.Drift,
			Add:     response.Add,
			Change:  response.Change,
			Destroy: response.Destroy,
			TFver:   tfver,
		})

		tableRow := []string{cfg.Name,
			strconv.FormatBool(response.Drift),
			strconv.Itoa(response.Add),
			strconv.Itoa(response.Change),
			strconv.Itoa(response.Destroy),
			cfg.Path,
			tfver,
		}
		table = append(table, tableRow)

	}

	switch *output {
	case "json":
		o, err := json.Marshal(stackOutputs)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(string(o))
	case "yaml":
		o, err := yaml.Marshal(stackOutputs)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(string(o))
	case "table":
		tablePrinter(table, []string{"STACK-NAME", "DRIFT", "ADD", "CHANGE", "DESTROY", "PATH", "TF-VERSION"})

	}
}

func tablePrinter(data [][]string, columns []string) {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(columns)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(3)
	table.SetAlignment(3)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)
	table.AppendBulk(data) // Add Bulk Data
	table.Render()

}
