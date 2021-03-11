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

func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			totalPower.Set(23.3)
			time.Sleep(2 * time.Second)
		}
	}()
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})
	totalPower = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "pdu_processed_power_total",
		Help: "Total power consumption",
	})
)

type conf struct {
	A string
	B struct {
		RenamedC int   `yaml:"c"`
		D        []int `yaml:",flow"`
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
	fmt.Printf("--- t:\n%v\n\n", c)

	return c
}

func main() {
	app := &cli.App{
			Name: "PDU Data Exporter",
			Usage: "PDU power exporter for prometheus",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "config",
					Aliases: []string{"c"},
					Usage:   "Load configuration from `FILE`",
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