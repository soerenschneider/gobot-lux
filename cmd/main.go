package main

import (
	"flag"
	"fmt"
	"gobot-lux/internal"
	"gobot-lux/internal/config"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/platforms/firmata"
	"gobot.io/x/gobot/platforms/mqtt"
	"log"
	"time"
)

func main() {
	log.Printf("Started %s, version %s, commit %s, built at %s", config.BotName, internal.BuildVersion, internal.CommitHash, internal.BuildTime)
	conf := getConfig()
	if conf.MetricConfig != "" {
		go internal.StartMetricsServer(conf.MetricConfig)
	}

	adaptor := firmata.NewAdaptor(conf.FirmAtaPort)
	driver := aio.NewAnalogSensorDriver(adaptor, conf.AioPin, time.Millisecond*time.Duration(conf.AioPollingIntervalMs))
	clientId := fmt.Sprintf("%s_%s", config.BotName, conf.Location)
	mqttAdaptor := mqtt.NewAdaptor(conf.MqttConfig.Host, clientId)
	adaptors := &internal.BrightnessBot{
		Driver:      driver,
		Adaptor:     adaptor,
		MqttAdaptor: mqttAdaptor,
		Config:      conf,
	}

	bot := internal.AssembleBot(adaptors)
	err := bot.Start()
	if err != nil {
		log.Fatalf("could not start bot: %v", err)
	}
}

func getConfig() config.Config {
	var configFile string
	flag.StringVar(&configFile, "config", "", "File to read configuration from")
	flag.Parse()
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

	conf.Print()
	err = conf.Validate()
	if err != nil {
		log.Fatalf("Invalid config: %v", err)
	}
	conf.FormatTopic()

	return *conf
}
