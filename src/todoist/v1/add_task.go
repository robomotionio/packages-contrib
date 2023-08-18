package todoist

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/robomotionio/robomotion-go/message"
	"github.com/robomotionio/robomotion-go/runtime"
	"net/http"
)

type AddTask struct {
	runtime.Node `spec:"id=Robomotion.Todoist.AddTask,name=Add Task,icon=mdiCactus,color=#E44331"`

	//Input
	InConnectionId runtime.InVariable[string] `spec:"title=Connection Id,type=string,scope=Message,name=connection_id,messageScope,customScope"`
	InProjectId    runtime.InVariable[string] `spec:"title=Project Id,type=string,scope=Message,name=project_id,messageScope,customScope"`
	InTaskName     runtime.InVariable[string] `spec:"title=Task Name,type=string,scope=Custom,messageScope,customScope"`
	InDueString    runtime.InVariable[string] `spec:"title=Due String,type=string,scope=Custom,messageScope,customScope"`

	//Output
	OutTask runtime.OutVariable[interface{}] `spec:"title=Result,type=object,scope=Message,name=task,messageScope,customScope"`

	//Options

}

func (n *AddTask) OnCreate() error {

	return nil
}

func (n *AddTask) OnMessage(ctx message.Context) error {
	connection_id, err := n.InConnectionId.Get(ctx)
	if err != nil {
		return err
	}
	bearer := getToken(connection_id)
	content, err := n.InTaskName.Get(ctx)
	if err != nil {
		return err
	}
	project_id, err := n.InProjectId.Get(ctx)
	if err != nil {
		return err
	}
	due_string, err := n.InDueString.Get(ctx)
	if err != nil {
		return err
	}
	type Data struct {
		Content   string `json:"content"`
		ProjectId string `json:"project_id"`
		DueString string `json:"due_string"`
	}
	const endpoint = "https://api.todoist.com/rest/v2/tasks"
	data := Data{
		Content:   content,
		ProjectId: project_id,
		DueString: due_string,
	}
	marshalled, _ := json.Marshal(data)

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(marshalled))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", bearer)
	req.Header.Set("X-Request-ID", uuid.New().String())
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	var responseBody interface{}
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return err
	}

	return n.OutTask.Set(ctx, responseBody)

}

func (n *AddTask) OnClose() error {

	return nil
}
