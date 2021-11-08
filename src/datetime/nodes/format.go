package nodes

import (
	"time"

	"github.com/robomotionio/robomotion-go/message"
	"github.com/robomotionio/robomotion-go/runtime"
)

type Format struct {
	runtime.Node `spec:"id=Robomotion.DateTime.Format,name=Format Time,icon=mdiTimelineText,color=#77AF38"`

	//Input
	InTime runtime.InVariable `spec:"title=Time,type=string,scope=Custom,messageScope,customScope,format=datetime"`

	//Output
	OutFormattedTime runtime.OutVariable `spec:"title=Formatted Time,type=string,scope=Message,name=time,messageScope,customScope"`

	//Options
	OptInLayout  string `spec:"title=In Layout,value=RFC3339,enum=ANSIC|UnixDate|RubyDate|RFC822|RFC822Z|RFC850|RFC1123|RFC1123Z|RFC3339|RFC3339Nano,enumNames=ANSIC|UnixDate|RubyDate|RFC822|RFC822Z|RFC850|RFC1123|RFC1123Z|RFC3339|RFC3339Nano,option"`
	OptOutLayout string `spec:"title=Out Layout,value=RFC3339,enum=ANSIC|UnixDate|RubyDate|RFC822|RFC822Z|RFC850|RFC1123|RFC1123Z|RFC3339|RFC3339Nano,enumNames=ANSIC|UnixDate|RubyDate|RFC822|RFC822Z|RFC850|RFC1123|RFC1123Z|RFC3339|RFC3339Nano,option"`
}

func (n *Format) OnCreate() error {

	return nil
}

func (n *Format) OnMessage(ctx message.Context) error {

	inLayout, ok := layouts[n.OptInLayout]
	if !ok {
		inLayout = time.RFC3339
	}

	outLayout, ok := layouts[n.OptOutLayout]
	if !ok {
		outLayout = time.RFC3339
	}

	inTime, err := n.InTime.GetString(ctx)
	if err != nil {
		return err
	}

	t1, err := time.Parse(inLayout, inTime)
	if err != nil {
		return err
	}

	outTime := t1.Format(outLayout)
	return n.OutFormattedTime.Set(ctx, outTime)
}

func (n *Format) OnClose() error {

	return nil
}
