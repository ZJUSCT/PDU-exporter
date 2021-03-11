package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	powerData = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "ZJUSCT",
		Name:      "pdu_processed_power",
		Help:      "node power consumption",
	},
		[]string{
		"node",
		"place",
		})
)

type conf struct {
	Url   string
	Nodes []struct {
		Name  string `yaml:"name"`
		Place []int  `yaml:",flow"`
	}
}

func parseYamlConfig(filename string) conf {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	c := conf{}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return c
}

func recordMetrics() {
	go func() {
		for {
			powerData.WithLabelValues("mu00", "1").Set(233.3)
			powerData.WithLabelValues("mu00", "2").Set(666.6)
			time.Sleep(2 * time.Second)
		}
	}()
}

func main() {
	app := &cli.App{
		Name:  "PDU Data Exporter",
		Usage: "PDU power exporter for prometheus",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config",
				Aliases:  []string{"c"},
				Usage:    "Load configuration from `FILE`",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			config := parseYamlConfig(c.String("config"))
			fmt.Printf("config:\n%v\n", config)

			recordMetrics()
			http.Handle("/metrics", promhttp.Handler())
			_ = http.ListenAndServe(":2112", nil)

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
