package nodes

import (
	"time"

	"github.com/robomotionio/robomotion-go/message"
	"github.com/robomotionio/robomotion-go/runtime"
)

type Split struct {
	runtime.Node `spec:"id=Robomotion.DateTime.Split,name=Split Date,icon=mdiCalendarBlankMultiple,color=#77AF38"`

	//Input
	InTime runtime.InVariable `spec:"title=Time,type=string,scope=Message,name=time,messageScope,customScope"`

	//Output
	OutParts runtime.OutVariable `spec:"title=Parts,type=object,scope=Message,name=parts,messageScope,customScope"`

	//Options
	OptLayout string `spec:"title=Layout,value=RFC3339,enum=ANSIC|UnixDate|RubyDate|RFC822|RFC822Z|RFC850|RFC1123|RFC1123Z|RFC3339|RFC3339Nano,enumNames=ANSIC|UnixDate|RubyDate|RFC822|RFC822Z|RFC850|RFC1123|RFC1123Z|RFC3339|RFC3339Nano,option"`
}

func (n *Split) OnCreate() error {

	return nil
}

func (n *Split) OnMessage(ctx message.Context) error {

	layout, ok := layouts[n.OptLayout]
	if !ok {
		layout = time.RFC3339
	}

	inTime, err := n.InTime.GetString(ctx)
	if err != nil {
		return err
	}

	t1, err := time.Parse(layout, inTime)
	if err != nil {
		return err
	}

	parts := map[string]interface{}{
		"year":         t1.Year(),
		"month":        int(t1.Month()),
		"month_name":   t1.Month().String(),
		"day":          t1.Day(),
		"year_day":     t1.YearDay(),
		"weekday":      int(t1.Weekday()),
		"weekday_name": t1.Weekday().String(),
		"minute":       t1.Minute(),
		"second":       t1.Second(),
	}

	return n.OutParts.Set(ctx, parts)
}

func (n *Split) OnClose() error {

	return nil
}
