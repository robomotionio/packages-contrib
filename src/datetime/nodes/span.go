package nodes

import (
	"time"

	"github.com/robomotionio/robomotion-go/message"
	"github.com/robomotionio/robomotion-go/runtime"
)

type Span struct {
	runtime.Node `spec:"id=Robomotion.DateTime.Span,name=Time Span,icon=mdiTimelapse,color=#77AF38"`

	//Input
	InStartDate runtime.InVariable `spec:"title=Start Date,type=string,scope=Message,name=start,messageScope,customScope"`
	InEndDate   runtime.InVariable `spec:"title=End Date,type=string,scope=Message,name=end,messageScope,customScope"`

	//Output
	OutSpan runtime.OutVariable `spec:"title=Time Span(ms),type=string,scope=Message,name=span,messageScope,customScope"`

	//Options
	OptLayout string `spec:"title=Layout,value=RFC3339,enum=ANSIC|UnixDate|RubyDate|RFC822|RFC822Z|RFC850|RFC1123|RFC1123Z|RFC3339|RFC3339Nano,enumNames=ANSIC|UnixDate|RubyDate|RFC822|RFC822Z|RFC850|RFC1123|RFC1123Z|RFC3339|RFC3339Nano,option"`
}

func (n *Span) OnCreate() error {

	return nil
}

func (n *Span) OnMessage(ctx message.Context) error {

	layout, ok := layouts[n.OptLayout]
	if !ok {
		layout = time.RFC3339
	}

	startDate, err := n.InStartDate.GetString(ctx)
	if err != nil {
		return err
	}

	endDate, err := n.InEndDate.GetString(ctx)
	if err != nil {
		return err
	}

	t1, err := time.Parse(layout, startDate)
	if err != nil {
		return err
	}

	t2, err := time.Parse(layout, endDate)
	if err != nil {
		return err
	}

	span := t2.Sub(t1)
	return n.OutSpan.Set(ctx, span.Milliseconds())
}

func (n *Span) OnClose() error {

	return nil
}
