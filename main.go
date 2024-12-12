package main

import (
	"encoding/json"
	"fmt"
	"github.com/fzf-labs/openapi-jmeter/config"
	"github.com/fzf-labs/openapi-jmeter/openapi"
	"github.com/spf13/cobra"
	"log"
)

var CmdJmeter = &cobra.Command{
	Use:   "jmeter",
	Short: "Generate JMeter test scripts",
	Long:  "Generate JMeter test scripts from Swagger documentation. Example: jmeter",
	Run:   run,
}

var (
	conf string
)

func init() {
	CmdJmeter.Flags().StringVarP(&conf, "config", "c", "./config/config.yaml", "config file")
}

func run(_ *cobra.Command, _ []string) {
	c, err := config.NewConfig(conf)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	apis, err := openapi.NewOpenAPI(c).Run()
	if err != nil {
		log.Fatalf("Failed to load OpenAPI: %v", err)
		return
	}
	b, _ := json.Marshal(apis)
	fmt.Println(string(b))
}

func main() {
	err := CmdJmeter.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
