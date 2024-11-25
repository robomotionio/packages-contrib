package language

import (
	"context"
	"log"

	language "cloud.google.com/go/language/apiv1"
	"github.com/robomotionio/robomotion-go/message"
	"github.com/robomotionio/robomotion-go/runtime"
	"google.golang.org/api/option"
)

type Connect struct {
	runtime.Node `spec:"id=Robomotion.GoogleNaturalLanguage.Connect,name=Connect,icon=mdiLanConnect,color=#9C27B0"`

	// Output
	OutClientID runtime.OutVariable[string] `spec:"title=Client ID,type=string,scope=Message,name=client_id,messageScope,customScope"`

	// Option
	OptCredential runtime.Credential `spec:"title=Credentials,scope=Custom,category=6,messageScope,customScope"`
}

func (n *Connect) OnCreate() error {
	return nil
}

func (n *Connect) OnMessage(ctx message.Context) error {

	credential, err := n.OptCredential.Get(ctx)
	if err != nil {
		return err
	}

	credDocument, ok := credential["content"].(string)
	if !ok {
		return runtime.NewError("ErrInvalidArg", "No Credential Content")
	}

	gctx := context.Background()

	client, err := language.NewClient(gctx, option.WithCredentialsJSON([]byte(credDocument)))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	clientID := addClient(client)

	return n.OutClientID.Set(ctx, clientID)

}

func (n *Connect) OnClose() error {

	return nil
}
