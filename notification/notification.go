package notification

import "fmt"

type Notification interface {
	Send(msg string) error
}

func NewNotify(n map[string]any) (Notification, error) {
	name, ok := n["name"]
	if !ok {
		return nil, fmt.Errorf("name not found")
	}

	enable, ok := n["enable"]
	if !ok {
		return nil, fmt.Errorf("enable not found")
	}

	e, ok := enable.(bool)
	if !ok {
		return nil, fmt.Errorf("enable is not bool")
	}

	if !e {
		return nil, nil
	}

	webhook, ok := n["webhook"]
	if !ok {
		return nil, fmt.Errorf("webhook not found")
	}

	switch name {
	case "feishu":
		return &FeishuNotification{
			Name:    "feishu",
			Enable:  e,
			Webhook: webhook.(string),
		}, nil

	}

	return nil, nil
}
