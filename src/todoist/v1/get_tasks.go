package todoist

import (
	"encoding/json"
	"github.com/robomotionio/robomotion-go/message"
	"github.com/robomotionio/robomotion-go/runtime"
	"log"
	"net/http"
)

type GetTasks struct {
	runtime.Node `spec:"id=Robomotion.Todoist.GetTasks,name=Get Tasks,icon=mdiCactus,color=#E44331"`

	//Input
	InConnectionId runtime.InVariable[string] `spec:"title=Connection Id,type=string,scope=Message,name=connection_id,messageScope,customScope"`
	//Output
	OutTasks runtime.OutVariable[interface{}] `spec:"title=Result,type=object,scope=Message,name=task_list,messageScope,customScope"`

	//Options

}

func (n *GetTasks) OnCreate() error {

	return nil
}

func (n *GetTasks) OnMessage(ctx message.Context) error {
	connection_id, err := n.InConnectionId.Get(ctx)
	if err != nil {
		return err
	}
	bearer := getToken(connection_id)
	const endpoint = "https://api.todoist.com/rest/v2/tasks"
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", bearer)
	client := &http.Client{}
	response, er := client.Do(req)
	if err != nil {
		log.Fatalf("Impossible to send request: %s", er)
		return err
	}
	defer response.Body.Close()

	var responseBody interface{}
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return err
	}

	return n.OutTasks.Set(ctx, responseBody)

}

func (n *GetTasks) OnClose() error {

	return nil
}
