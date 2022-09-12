package protocols

import (
	"github.com/tidwall/gjson"
	adapters "sonic-ios-webkit-adapter/adapter"
)

type iOS12 struct {
	adapter *adapters.Adapter
}

func initIOS12(protocol *ProtocolAdapter) {
	protocol.adapter.SetTargetBased(true)
	result := &iOS12{
		adapter: protocol.adapter,
	}
	protocol.init()
	protocol.adapter.AddMessageFilter("Target.targetCreated", result.targetCreated)
}

func (i *iOS12) targetCreated(message []byte) []byte {
	i.adapter.SetTargetID(gjson.Get(string(message), "params.targetInfo.targetId").String())
	return message
}
