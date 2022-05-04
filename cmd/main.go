package main

import (
	"flag"
	"fmt"
	"github.com/soerenschneider/gobot-lux/internal"
	"github.com/soerenschneider/gobot-lux/internal/config"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/platforms/firmata"
	"gobot.io/x/gobot/platforms/mqtt"
	"log"
	"os"
	"time"
)

const (
	cliConfFile = "config"
	cliVersion  = "version"
)

func main() {
	conf := getConfig()
	log.Printf("Started %s, version %s, commit %s", config.BotName, internal.BuildVersion, internal.CommitHash)
	conf.FormatTopic()
	conf.Print()
	err := conf.Validate()
	if err != nil {
		log.Fatalf("Invalid config: %v", err)
	}

	if conf.MetricConfig != "" {
		go internal.StartMetricsServer(conf.MetricConfig)
	}

	adaptor := firmata.NewAdaptor(conf.FirmAtaPort)
	driver := aio.NewAnalogSensorDriver(adaptor, conf.AioPin, time.Millisecond*time.Duration(conf.AioPollingIntervalMs))
	clientId := fmt.Sprintf("%s_%s", config.BotName, conf.Placement)
	mqttAdaptor := mqtt.NewAdaptor(conf.MqttConfig.Host, clientId)
	mqttAdaptor.SetAutoReconnect(true)
	mqttAdaptor.SetQoS(1)

	if conf.MqttConfig.UsesSslCerts() {
		log.Println("Setting TLS client cert and key...")
		mqttAdaptor.SetClientCert(conf.MqttConfig.ClientCertFile)
		mqttAdaptor.SetClientKey(conf.MqttConfig.ClientKeyFile)

		if len(conf.MqttConfig.ServerCaFile) > 0 {
			log.Println("Setting server CA...")
			mqttAdaptor.SetServerCert(conf.MqttConfig.ServerCaFile)
		}
	}

	adaptors := &internal.BrightnessBot{
		Driver:      driver,
		Adaptor:     adaptor,
		MqttAdaptor: mqttAdaptor,
		Config:      conf,
	}

	bot := internal.AssembleBot(adaptors)
	err = bot.Start()
	if err != nil {
		log.Fatalf("could not start bot: %v", err)
	}
}

func getConfig() config.Config {
	var configFile string
	flag.StringVar(&configFile, cliConfFile, "", "File to read configuration from")
	version := flag.Bool(cliVersion, false, "Print version and exit")
	flag.Parse()

	if *version {
		fmt.Printf("%s (revision %s)", internal.BuildVersion, internal.CommitHash)
		os.Exit(0)
	}

	if configFile == "" {
		log.Println("Building config from env vars")
		return config.ConfigFromEnv()
	}

	log.Printf("Reading config from file %s", configFile)
	conf, err := config.ReadJsonConfig(configFile)
	if err != nil {
		log.Fatalf("Could not read config from %s: %v", configFile, err)
	}
	if nil == conf {
		log.Fatalf("Received empty config, should not happen")
	}

	return *conf
}
