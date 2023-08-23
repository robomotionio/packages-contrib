package language

import (
	"context"
	"fmt"
	"log"

	languagepb "cloud.google.com/go/language/apiv1/languagepb"
	"github.com/robomotionio/robomotion-go/message"
	"github.com/robomotionio/robomotion-go/runtime"
)

type AnalyzeEntities struct {
	runtime.Node `spec:"id=Robomotion.GoogleNaturalLanguage.AnalyzeEntities,name=Analyze Entities,icon=mdiCardText,color=#9C27B0"`

	// Input
	InClientID runtime.InVariable[string] `spec:"title=Client ID,type=string,cope=Message,name=client_id,messageScope,customScope"`
	InText     runtime.InVariable[string] `spec:"title=Input Text,type=string,cope=Message,name=text,messageScope,customScope"`

	// Output
	OutEntites runtime.OutVariable[interface{}] `spec:"title=Result,type=object,scope=Message,name=result,messageScope,customScope"`

	// Options

}

func (n *AnalyzeEntities) OnCreate() error {
	return nil
}

func (n *AnalyzeEntities) OnMessage(ctx message.Context) error {

	inputText, err := n.InText.Get(ctx)
	if err != nil {
		return err
	}

	clientID, err := n.InClientID.Get(ctx)
	if err != nil {
		return err
	}

	client := getClient(clientID)
	if client == nil {
		return fmt.Errorf("Invalid ClientId")
	}

	document := &languagepb.Document{
		Source: &languagepb.Document_Content{
			Content: inputText,
		},
		Type: languagepb.Document_PLAIN_TEXT,
	}

	req := &languagepb.AnalyzeEntitiesRequest{
		Document: document,
	}
	gctx := context.Background()

	resp, err := client.AnalyzeEntities(gctx, req)
	if err != nil {
		log.Fatalf("Failed to analyze entities: %v", err)
	}

	for _, entity := range resp.Entities {
		fmt.Printf("Entity: %v\n", entity.Name)
		fmt.Printf("Type: %v\n", entity.Type)
		fmt.Printf("Salience: %v\n", entity.Salience)
		fmt.Printf("Mentions:\n")
		for _, mention := range entity.Mentions {
			fmt.Printf("\tText: %v\n", mention.Text.Content)
			fmt.Printf("\tType: %v\n", mention.Type)
		}
		fmt.Println("--------------------")
	}

	return n.OutEntites.Set(ctx, resp.Entities)
}

func (n *AnalyzeEntities) OnClose() error {
	return nil
}
