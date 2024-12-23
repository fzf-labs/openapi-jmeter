package main

import (
	"log"

	"github.com/fzf-labs/openapi-jmeter/config"
	"github.com/fzf-labs/openapi-jmeter/jmeter"
	"github.com/fzf-labs/openapi-jmeter/openapi"
	"github.com/spf13/cobra"
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
	log.Println("Start generate jmeter script")
	// 读取配置文件
	c, err := config.NewConfig(conf)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Println("Load config success")
	// 获取OpenAPI文件
	apis, err := openapi.NewOpenAPI(c).Run()
	if err != nil {
		log.Fatalf("Failed to load OpenAPI: %v", err)
		return
	}
	log.Println("Get OpenAPI file success")
	// 生成JMeter脚本
	err = jmeter.NewJMeter(c, apis).GenerateJMX()
	if err != nil {
		log.Fatalf("Failed to generate JMX: %v", err)
	}
	log.Println("Generate JMX success")
}

func main() {
	err := CmdJmeter.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
