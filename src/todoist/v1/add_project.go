package todoist

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/robomotionio/robomotion-go/message"
	"github.com/robomotionio/robomotion-go/runtime"
	"net/http"
)

type AddProject struct {
	runtime.Node `spec:"id=Robomotion.Todoist.AddProject,name=Add Project,icon=mdiCactus,color=#E44331"`

	//Input
	InConnectionId runtime.InVariable[string] `spec:"title=Connection Id,type=string,scope=Message,name=connection_id,messageScope,customScope"`
	InName         runtime.InVariable[string] `spec:"title=Name,type=string,scope=Custom,messageScope,customScope"`
	//Output
	OutProject runtime.OutVariable[interface{}] `spec:"title=Result,type=object,scope=Message,name=project,messageScope,customScope"`

	//Options

}

func (n *AddProject) OnCreate() error {

	return nil
}

func (n *AddProject) OnMessage(ctx message.Context) error {
	connection_id, err := n.InConnectionId.Get(ctx)
	if err != nil {
		return err
	}
	bearer := getToken(connection_id)
	name, err := n.InName.Get(ctx)
	if err != nil {
		return err
	}
	type Data struct {
		Name string `json:"name"`
	}
	const endpoint = "https://api.todoist.com/rest/v2/projects"
	data := Data{
		Name: name,
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

	return n.OutProject.Set(ctx, responseBody)
}

func (n *AddProject) OnClose() error {

	return nil
}
