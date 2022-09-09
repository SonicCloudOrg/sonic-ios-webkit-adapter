package protocols

import (
	"github.com/tidwall/gjson"
	adapters "sonic-ios-webkit-adapter/adapter"
)

type IOS12 struct {
	adapter *adapters.Adapter
}

func (i *IOS12) targetCreated(message []byte) []byte {
	i.adapter.SetTargetID(gjson.Get(string(message), "params.targetInfo.targetId").String())
	return message
}
