package mqtt

import (
	"fmt"
	"log"
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/util"
	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
)

const (
	maxLength = 268435455
)

// Service sends notifications to mqtt topic
type Service struct {
	standard.Standard
	config *Config
	pkr    format.PropKeyResolver
}

// Send notification to mqtt
func (service *Service) Send(message string, params *types.Params) error {

	message, omitted := MessageLimit(message)

	if omitted > 0 {
		service.Logf("omitted %v character(s) from the message", omitted)
	}

	config := *service.config
	if err := service.pkr.UpdateConfigFromParams(&config, params); err != nil {
		return err
	}

	if err := service.PublishMessageToTopic(message, &config); err != nil {
		return fmt.Errorf("an error occurred while sending notification to the MQTT topic: %s", err.Error())
	}

	return nil
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service
func (service *Service) Initialize(configURL *url.URL, logger *log.Logger) error {
	service.Logger.SetLogger(logger)
	service.config = &Config{
		DisableTLS: false,
		Port:       8883,
	}
	service.pkr = format.NewPropKeyResolver(service.config)
	err := service.config.setURL(&service.pkr, configURL)

	return err
}

// MessageLimit returns a string with the maximum size and the amount of omitted characters
func MessageLimit(message string) (string, int) {
	size := util.Min(maxLength, len(message))
	omitted := len(message) - size

	return message[:size], omitted
}

// GetConfig returns the Config for the service
func (service *Service) GetConfig() *Config {
	return service.config
}

// Publish to topic
func (service *Service) Publish(client mqtt.Client, topic string, message string) {
	token := client.Publish(topic, 0, false, message)
	token.Wait()
}

// PublishMessageToTopic initializes the client and publishes the message
func (service *Service) PublishMessageToTopic(message string, config *Config) error {
	postURL := config.MqttURL()
	opts := config.GetClientConfig(postURL)
	client := mqtt.NewClient(opts)
	token := client.Connect()

	if token.Error() != nil {
		return token.Error()
	}

	token.Wait()

	service.Publish(client, config.Topic, message)

	client.Disconnect(250)

	return nil
}
