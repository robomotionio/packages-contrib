package language

import (
	"github.com/robomotionio/robomotion-go/message"
	"github.com/robomotionio/robomotion-go/runtime"
)

type Disconnect struct {
	runtime.Node `spec:"id=Robomotion.GoogleNaturalLanguage.Disconnect,name=Disconnect,icon=mdiLanDisconnect,color=#9C27B0"`

	//Input
	InClientID runtime.InVariable[string] `spec:"title=Client ID,type=string,scope=Message,name=client_id,messageScope,customScope"`

	//Output

	//Options
}

func (n *Disconnect) OnCreate() error {

	return nil
}

func (n *Disconnect) OnMessage(ctx message.Context) error {
	connection_id, err := n.InClientID.Get(ctx)
	if err != nil {
		return err
	}
	removeClient(connection_id)

	return nil

}

func (n *Disconnect) OnClose() error {

	return nil
}
