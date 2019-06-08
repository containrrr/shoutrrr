package pushbullet

import "regexp"

type Plugin struct {}

var (
	serviceUrl = "https://api.pushbullet.com/v2/pushes"
)

func (plugin *Plugin) Send(url string, message string) error {
	config, err := CreateConfigFromURL(url)
	if err != nil {
		return err
	}

	for _, target := range config.Targets {
		if err := doSend(config.Token, target, message); err != nil {
			return err
		}
	}
	return nil
}

func doSend(token string, target string, message string) error {
	_, err := getTargetType(target)
	if err != nil {
		return err
	}

	/*
		payload format:
		{
			"type": "note",
	        "title": title,
			"body": message,
			"x": target // replace x with email, channel_tag or device_iden based on target type
		}
	 */

	return nil
}

func getTargetType(target string) (TargetType, error) {
	matchesEmail, err := regexp.MatchString(`.*@.*\..*`, target)

	if matchesEmail && err == nil {
		return EmailTarget, nil
	} else if string(target[0]) == "#" {
		return ChannelTarget, nil
	} else {
		return DeviceTarget, nil
	}
}


type TargetType int

const (
	EmailTarget TargetType = 1
	ChannelTarget TargetType = 2
	DeviceTarget TargetType = 3
)