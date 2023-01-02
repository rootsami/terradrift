package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/olekukonko/tablewriter"
	"github.com/rootsami/terradrift/pkg/config"
	"github.com/rootsami/terradrift/pkg/tfstack"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"

	log "github.com/sirupsen/logrus"
)

var (
	app              = kingpin.New("terradrift-cli", "A command-line tool to detect drifts in terraform IaC")
	workdir          = app.Flag("workdir", "workdir of a project that contains all terraform directories").Default("./").String()
	configPath       = app.Flag("config", "Path for configuration file holding the stack information").String()
	extraBackendVars = app.Flag("extra-backend-vars", "Extra backend environment variables ex. GOOGLE_CREDENTIALS, AWS_ACCESS_KEY_ID or AWS_CONFIG_FILE ..etc").StringMap()
	debug            = app.Flag("debug", "Enable debug mode").Default("false").Bool()
	generateConfig   = app.Flag("generate-config-only", "Generate a config file based on a provided worksapce").Default("false").Bool()
	output           = app.Flag("output", "Output format supported: json, yaml and table").Default("table").Enum("table", "json", "yaml")
	version          string
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

func init() {

	app.Version(version)
	kingpin.MustParse(app.Parse(os.Args[1:]))

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	// if worksapce path is not absolute, make it absolute and add a trailing slash
	if !path.IsAbs(*workdir) {
		absPath, err := filepath.Abs(*workdir)
		if err != nil {
			log.Fatalf("Error getting absolute path for workdir: %s", err)
		}

		*workdir = absPath + "/"

	} else if !strings.HasSuffix(*workdir, "/") {
		*workdir = *workdir + "/"
	}
}

func main() {

	var cfg *config.Config
	var stackOutputs []stackOutput
	var err error

	switch {
	// if config file is provided, load it and assign it to cfg
	case *configPath != "":
		log.WithFields(log.Fields{"config": *configPath}).Debug("loading config file")
		cfg, err = config.ConfigLoader(*workdir, *configPath)
		if err != nil {
			log.Fatalf("Error loading config file: %s", err)
		}

	// if --generate-config-only flag is provided, generate config file to stdout and exit
	case *generateConfig:
		cfg, err = config.ConfigGenerator(*workdir)
		if err != nil {
			log.Fatalf("Error generating config file: %s", err)
		}

		outputWriter(cfg, "yaml")
		os.Exit(0)

	// if config file is not provided, generate it and assign it to cfg
	case *configPath == "":

		log.Debug("config file not found, running stack init on each directory that contains .tf files")
		cfg, err = config.ConfigGenerator(*workdir)
		if err != nil {
			log.Fatalf("error generating config file: %s", err)
		}

	}

	var wg sync.WaitGroup
	for _, stack := range cfg.Stacks {

		wg.Add(1)
		go func(s config.Stack) {
			defer wg.Done()

			// catch panic and log it as error to continue to the next stack
			defer func() {
				if r := recover(); r != nil {
					log.WithFields(log.Fields{"stack": s.Name}).Error(r)
				}
			}()

			response, tfver, err := tfstack.StackInit(*workdir, s, *extraBackendVars)
			if err != nil {
				log.WithFields(log.Fields{"stack": s.Name}).Error(err)
			}

			stackOutputs = append(stackOutputs, stackOutput{
				Name:    s.Name,
				Path:    s.Path,
				Drift:   response.Drift,
				Add:     response.Add,
				Change:  response.Change,
				Destroy: response.Destroy,
				TFver:   tfver,
			})
		}(stack)

	}
	wg.Wait()

	// output the results based on the output flag
	switch *output {
	case "json":
		outputWriter(stackOutputs, "json")
	case "yaml":
		outputWriter(stackOutputs, "yaml")
	case "table":
		tableWriter(stackOutputs)

	}
}

func tableWriter(stackOutputs []stackOutput) {

	columns := []string{"STACK-NAME", "DRIFT", "ADD", "CHANGE", "DESTROY", "PATH", "TF-VERSION"}
	var data [][]string

	for _, stackOutput := range stackOutputs {
		row := []string{stackOutput.Name,
			strconv.FormatBool(stackOutput.Drift),
			strconv.Itoa(stackOutput.Add),
			strconv.Itoa(stackOutput.Change),
			strconv.Itoa(stackOutput.Destroy),
			stackOutput.Path,
			stackOutput.TFver,
		}
		data = append(data, row)
	}

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

// outputWriter takes a data interface and a format string and outputs the data in the specified format
func outputWriter(data interface{}, format string) {

	switch format {
	case "json":
		o, err := json.Marshal(data)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(string(o))
	case "yaml":
		o, err := yaml.Marshal(data)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(string(o))
	}

}
