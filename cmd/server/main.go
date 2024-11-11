package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/mwantia/metrics-merger/pkg/common"
	pkg_http "github.com/mwantia/metrics-merger/pkg/http"
)

func init() {
	flag.Set("logtostderr", "true")
}

var (
	Config = flag.String("config", "", "Defines the configuration path used by this application")
)

func main() {
	flag.Parse()

	if strings.TrimSpace(*Config) == "" {
		log.Fatal(fmt.Errorf("configuration path has not been defined"))
	}

	cfg, err := common.LoadServerConfig(*Config)
	if err != nil {
		log.Fatal(err)
	}

	if err := pkg_http.CreateAndServe(cfg); err != nil {
		log.Fatal(err)
	}
}
