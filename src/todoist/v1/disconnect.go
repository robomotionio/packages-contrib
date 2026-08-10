package todoist

import (
	"github.com/robomotionio/robomotion-go/message"
	"github.com/robomotionio/robomotion-go/runtime"
)

type Disconnect struct {
	runtime.Node `spec:"id=Robomotion.Todoist.Disconnect,name=Disconnect,icon=mdiCactus,color=#E44331"`

	//Input
	InConnectionId runtime.InVariable[string] `spec:"title=Connection Id,type=string,scope=Message,name=connection_id,messageScope,customScope"`

	//Output

	//Options
}

func (n *Disconnect) OnCreate() error {

	return nil
}

func (n *Disconnect) OnMessage(ctx message.Context) error {
	connection_id, err := n.InConnectionId.Get(ctx)
	if err != nil {
		return err
	}
	removeToken(connection_id)
	return nil

}

func (n *Disconnect) OnClose() error {

	return nil
}
