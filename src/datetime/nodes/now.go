package nodes

import (
	"time"

	"github.com/robomotionio/robomotion-go/message"
	"github.com/robomotionio/robomotion-go/runtime"
)

type Now struct {
	runtime.Node `spec:"id=Robomotion.DateTime.Now,name=Now,icon=mdiCalendarClock,color=#77AF38"`

	//Input

	//Output
	OutNow runtime.OutVariable `spec:"title=Time,type=string,scope=Message,name=now,messageScope,customScope"`

	//Options
	OptLayout string `spec:"title=Layout,value=RFC3339,enum=ANSIC|UnixDate|RubyDate|RFC822|RFC822Z|RFC850|RFC1123|RFC1123Z|RFC3339|RFC3339Nano,enumNames=ANSIC|UnixDate|RubyDate|RFC822|RFC822Z|RFC850|RFC1123|RFC1123Z|RFC3339|RFC3339Nano,option"`
}

func (n *Now) OnCreate() error {

	return nil
}

func (n *Now) OnMessage(ctx message.Context) error {

	layout, ok := layouts[n.OptLayout]
	if !ok {
		layout = time.RFC3339
	}

	now := time.Now().Format(layout)
	return n.OutNow.Set(ctx, now)
}

func (n *Now) OnClose() error {

	return nil
}
