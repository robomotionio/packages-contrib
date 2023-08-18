package todoist

import (
	"fmt"
	"github.com/robomotionio/robomotion-go/message"
	"github.com/robomotionio/robomotion-go/runtime"
	"net/http"
)

type DeleteTask struct {
	runtime.Node `spec:"id=Robomotion.Todoist.DeleteTask,name=Delete Task,icon=mdiCactus,color=#E44331"`

	//Input
	InConnectionId runtime.InVariable[string] `spec:"title=Connection Id,type=string,scope=Message,name=connection_id,messageScope,customScope"`
	InTaskId       runtime.InVariable[string] `spec:"title=Task Id,type=string,scope=Message,name=task_id,messageScope,customScope"`

	//Output

	//Options

}

func (n *DeleteTask) OnCreate() error {

	return nil
}

func (n *DeleteTask) OnMessage(ctx message.Context) error {
	connection_id, err := n.InConnectionId.Get(ctx)
	if err != nil {
		return err
	}
	bearer := getToken(connection_id)
	task_id, err := n.InTaskId.Get(ctx)
	if err != nil {
		return err
	}
	endpoint := fmt.Sprintf("https://api.todoist.com/rest/v2/tasks/%s", task_id)
	req, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", bearer)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return nil
}

func (n *DeleteTask) OnClose() error {

	return nil
}
