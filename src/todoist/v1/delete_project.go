package todoist

import (
	"fmt"
	"github.com/robomotionio/robomotion-go/message"
	"github.com/robomotionio/robomotion-go/runtime"
	"net/http"
)

type DeleteProject struct {
	runtime.Node `spec:"id=Robomotion.Todoist.DeleteProject,name=Delete Project,icon=mdiCactus,color=#E44331"`

	//Input
	InConnectionId runtime.InVariable[string] `spec:"title=Connection Id,type=string,scope=Message,name=connection_id,messageScope,customScope"`
	InProjectId    runtime.InVariable[string] `spec:"title=Project Id,type=string,scope=Message,name=project_id,messageScope,customScope"`

	//Output

	//Options

}

func (n *DeleteProject) OnCreate() error {

	return nil
}

func (n *DeleteProject) OnMessage(ctx message.Context) error {
	connection_id, err := n.InConnectionId.Get(ctx)
	if err != nil {
		return err
	}
	bearer := getToken(connection_id)
	project_id, err := n.InProjectId.Get(ctx)
	if err != nil {
		return err
	}
	endpoint := fmt.Sprintf("https://api.todoist.com/rest/v2/projects/%s", project_id)
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

func (n *DeleteProject) OnClose() error {

	return nil
}
