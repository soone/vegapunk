package notification

import (
	"time"

	"github.com/guonaihong/gout"
)

type FeishuNotification struct {
	Name    string `json:"name" mapstructure:"name"`
	Enable  bool   `json:"enable" mapstructure:"enable"`
	Webhook string `json:"webhook" mapstructure:"webhook"`
}

func (f *FeishuNotification) Send(msg string) error {
	return gout.POST(f.Webhook).SetJSON(gout.H{"msg_type": "text", "content": gout.H{"text": msg}}).SetTimeout(time.Second * 10).Do()
}
