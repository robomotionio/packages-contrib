package datetime

import (
	"time"

	"github.com/robomotionio/robomotion-go/message"
	"github.com/robomotionio/robomotion-go/runtime"
)

type Add struct {
	runtime.Node `spec:"id=Robomotion.DateTime.Add,name=Add Time,icon=mdiTimelinePlus,color=#77AF38"`

	//Input
	InTime     runtime.InVariable[string] `spec:"title=Time,type=string,scope=Custom,messageScope,customScope,format=datetime"`
	InDuration runtime.InVariable[int64]  `spec:"title=Duration(ms),type=number,scope=Message,name=duration,messageScope,customScope"`

	//Output
	OutTimeResult runtime.OutVariable[string] `spec:"title=Time Result,type=string,scope=Message,name=result,messageScope,customScope"`

	//Options
	OptLayout string `spec:"title=Layout,value=RFC3339,enum=ANSIC|UnixDate|RubyDate|RFC822|RFC822Z|RFC850|RFC1123|RFC1123Z|RFC3339|RFC3339Nano,enumNames=ANSIC|UnixDate|RubyDate|RFC822|RFC822Z|RFC850|RFC1123|RFC1123Z|RFC3339|RFC3339Nano,option"`
}

func (n *Add) OnCreate() error {

	return nil
}

func (n *Add) OnMessage(ctx message.Context) error {

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

	duration, err := n.InDuration.Get(ctx)
	if err != nil {
		return err
	}

	result := t1.Add(time.Duration(duration) * time.Millisecond)
	return n.OutTimeResult.Set(ctx, result.String())
}

func (n *Add) OnClose() error {

	return nil
}
