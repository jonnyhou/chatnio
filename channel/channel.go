package channel

import (
	"chat/utils"
	"errors"
	"fmt"
	"strings"
)

var defaultMaxRetries = 1

func (c *Channel) GetId() int {
	return c.Id
}

func (c *Channel) GetName() string {
	return c.Name
}

func (c *Channel) GetType() string {
	return c.Type
}

func (c *Channel) GetPriority() int {
	return c.Priority
}

func (c *Channel) GetWeight() int {
	if c.Weight <= 0 {
		return 1
	}
	return c.Weight
}

func (c *Channel) GetModels() []string {
	return c.Models
}

func (c *Channel) GetRetry() int {
	if c.Retry <= 0 {
		return defaultMaxRetries
	}
	return c.Retry
}

func (c *Channel) GetSecret() string {
	return c.Secret
}

// GetRandomSecret returns a random secret from the secret list
func (c *Channel) GetRandomSecret() string {
	arr := strings.Split(c.GetSecret(), "\n")
	idx := utils.Intn(len(arr))
	return arr[idx]
}

func (c *Channel) SplitRandomSecret(num int) []string {
	secret := c.GetRandomSecret()
	arr := strings.Split(secret, "|")
	if len(arr) == num {
		return arr
	} else if len(arr) > num {
		return arr[:num]
	}

	for i := len(arr); i < num; i++ {
		arr = append(arr, "")
	}

	return arr
}

func (c *Channel) GetEndpoint() string {
	return c.Endpoint
}

func (c *Channel) GetMapper() string {
	return c.Mapper
}

func (c *Channel) GetReflect() map[string]string {
	if c.Reflect == nil {
		var reflect map[string]string
		arr := strings.Split(c.GetMapper(), "\n")
		for _, item := range arr {
			pair := strings.Split(item, ">")
			if len(pair) == 2 {
				reflect[pair[0]] = pair[1]
			}
		}

		c.Reflect = &reflect
	}

	return *c.Reflect
}

// GetModelReflect returns the reflection model name if it exists, otherwise returns the original model name
func (c *Channel) GetModelReflect(model string) string {
	ref := c.GetReflect()
	if reflect, ok := ref[model]; ok && len(reflect) > 0 {
		return reflect
	}

	return model
}

func (c *Channel) GetHitModels() []string {
	if c.HitModels == nil {
		var res []string

		models := c.GetModels()
		ref := c.GetReflect()

		for _, model := range models {
			if !utils.Contains(model, res) {
				res = append(res, model)
			}
		}

		for model := range ref {
			if !utils.Contains(model, res) {
				res = append(res, model)
			}
		}

		c.HitModels = &res
	}

	return *c.HitModels
}

func (c *Channel) GetState() bool {
	return c.State
}

func (c *Channel) IsHit(model string) bool {
	return utils.Contains(model, c.GetHitModels())
}

func (c *Channel) ProcessError(err error) error {
	if err == nil {
		return nil
	}
	content := err.Error()
	if strings.Contains(content, c.GetEndpoint()) {
		// hide the endpoint
		replacer := fmt.Sprintf("channel://%d", c.GetId())
		content = strings.Replace(content, c.GetEndpoint(), replacer, -1)
	}

	return errors.New(content)
}
