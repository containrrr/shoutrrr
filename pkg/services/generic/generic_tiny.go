//go:build tinygo
// +build tinygo

package generic

import (
	"errors"

	"github.com/containrrr/shoutrrr/pkg/types"
)

func (service *Service) doSend(config *Config, message string, params *types.Params) error {
	return errors.New("not supported in tinygo")
}
