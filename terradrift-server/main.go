package main

import (
	"io/ioutil"
	"os"

	"github.com/rootsami/terradrift/pkg/server"
	log "github.com/sirupsen/logrus"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	app              = kingpin.New("terradrift-server", "A tool to detect drifts in terraform IaC, As a server mode it will expose a rest api to query the drifts and also as prometheus metrics on /metrics endpoint")
	hostname         = app.Flag("hostname", "hostname that apil will be exposed.").Default("localhost").String()
	port             = app.Flag("port", "port of the service api is listening on").Default("8080").String()
	scheme           = app.Flag("scheme", "The scheme of exposed endpoint http/https").Default("http").String()
	repository       = app.Flag("repository", "The git repository which include all terraform stacks ").Required().String()
	gitToken         = app.Flag("git-token", "Personal access token to access git repositories").Required().String()
	gitTimeout       = app.Flag("git-timeout", "Wait timeout for git repoistory to clone or pull updates").Default("120").Int()
	interval         = app.Flag("interval", "The interval for scan scheduler").Default("60").Int()
	configPath       = app.Flag("config", "Path for configuration file holding the stack information").Default("config.yaml").String()
	extraBackendVars = app.Flag("extra-backend-vars", "Extra backend environment variables ex. GOOGLE_CREDENTIALS, AWS_ACCESS_KEY_ID or AWS_CONFIG_FILE ..etc").StringMap()
	debug            = app.Flag("debug", "Enable debug mode").Default("false").Bool()
	workdir          string
	version          string
)

func init() {
	dir, err := ioutil.TempDir("", "terradrift")
	if err != nil {
		log.Fatal(err)
	}
	workdir = dir + "/"

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006/01/02 - 15:04:05",
	})

}

func main() {

	app.Version(version)
	kingpin.MustParse(app.Parse(os.Args[1:]))

	args := server.Server{
		Workdir:          workdir,
		GitToken:         *gitToken,
		GitTimeout:       *gitTimeout,
		ConfigPath:       *configPath,
		ExtraBackendVars: *extraBackendVars,
		Interval:         *interval,
		Repository:       *repository,
		Scheme:           *scheme,
		Hostname:         *hostname,
		Port:             *port,
		Debug:            *debug,
	}

	err := server.Server.Start(args)
	if err != nil {
		log.Fatal(err)
	}

}
