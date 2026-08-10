package todoist

import (
	"github.com/robomotionio/robomotion-go/message"
	"github.com/robomotionio/robomotion-go/runtime"
)

type Connect struct {
	runtime.Node `spec:"id=Robomotion.Todoist.Connect,name=Connect,icon=mdiCactus,color=#E44331"`

	//Input

	//Output
	OutConnectionId runtime.OutVariable[string] `spec:"title=Connection Id,type=string,scope=Message,name=connection_id,messageScope,customScope"`

	//Options
	OptToken runtime.Credential `spec:"title=Todoist Token,scope=Custom,category=4,messageScope,customScope"`
}

func (n *Connect) OnCreate() error {

	return nil
}

func (n *Connect) OnMessage(ctx message.Context) error {
	item, err := n.OptToken.Get(ctx)
	if err != nil {
		return err
	}

	token, ok := item["value"].(string)
	if !ok {
		return runtime.NewError("ErrInvalidArg", "No Token Value")
	}

	clientID := addToken(token)
	return n.OutConnectionId.Set(ctx, clientID)

}

func (n *Connect) OnClose() error {

	return nil
}
