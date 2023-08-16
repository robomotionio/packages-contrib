package datetime

import (
	"time"

	"github.com/robomotionio/robomotion-go/message"
	"github.com/robomotionio/robomotion-go/runtime"
)

type Leap struct {
	runtime.Node `spec:"id=Robomotion.DateTime.Leap,name=Is Leap Year,icon=mdiCalendarStar,color=#77AF38"`

	//Input
	InTime runtime.InVariable[string] `spec:"title=Date,type=string,scope=Custom,customScope,format=date"`

	//Output
	OutIsLeap runtime.OutVariable[bool] `spec:"title=Is Leap,type=boolean,scope=Message,name=leap,messageScope,customScope"`

	//Options
	OptLayout string `spec:"title=Layout,value=RFC3339,enum=ANSIC|UnixDate|RubyDate|RFC822|RFC822Z|RFC850|RFC1123|RFC1123Z|RFC3339|RFC3339Nano,enumNames=ANSIC|UnixDate|RubyDate|RFC822|RFC822Z|RFC850|RFC1123|RFC1123Z|RFC3339|RFC3339Nano,option"`
}

func (n *Leap) OnCreate() error {

	return nil
}

func (n *Leap) OnMessage(ctx message.Context) error {

	layout, ok := layouts[n.OptLayout]
	if !ok {
		layout = time.RFC3339
	}

	inTime, err := n.InTime.Get(ctx)
	if err != nil {
		return err
	}

	t1, err := time.Parse(layout, inTime)
	if err != nil {
		return err
	}

	return n.OutIsLeap.Set(ctx, isLeap(t1.Year()))
}

func (n *Leap) OnClose() error {

	return nil
}

func isLeap(year int) bool {
	return year%400 == 0 || year%4 == 0 && year%100 != 0
}
