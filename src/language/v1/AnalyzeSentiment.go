package language

import (
	"context"
	"fmt"
	"log"

	languagepb "cloud.google.com/go/language/apiv1/languagepb"
	"github.com/robomotionio/robomotion-go/message"
	"github.com/robomotionio/robomotion-go/runtime"
)

type SentimentAnalysisResult struct {
	SentimentScore float32 `json:"sentiment_score"`
	SentimentLabel string  `json:"sentiment_label"`
}

type AnalyzeSentiment struct {
	runtime.Node `spec:"id=Robomotion.GoogleNaturalLanguage.AnalyzeSentiment,name=Analyze Sentiment,icon=mdiCardText,color=#9C27B0"`

	// Input
	InClientID runtime.InVariable[string] `spec:"title=Client ID,type=string,cope=Message,name=client_id,messageScope,customScope"`
	InText     runtime.InVariable[string] `spec:"title=Input Text,type=string, cope=Message,name=text,messageScope,customScope"`

	// Outpute
	OutSentiment runtime.OutVariable[interface{}] `spec:"title=Result,type=object,scope=Message,name=result,messageScope,customScope"`

	// Options

}

func (n *AnalyzeSentiment) OnCreate() error {
	return nil
}

func (n *AnalyzeSentiment) OnMessage(ctx message.Context) error {

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
	gctx := context.Background()
	sentiment, err := client.AnalyzeSentiment(gctx, &languagepb.AnalyzeSentimentRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: inputText,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
		EncodingType: languagepb.EncodingType_UTF8,
	})
	if err != nil {
		log.Fatalf("Failed to analyze text: %v", err)
	}
	var sentimentLabel string
	if sentiment.DocumentSentiment.Score >= 0 {
		sentimentLabel = "positive"
	} else {
		sentimentLabel = "negative"
	}

	result := SentimentAnalysisResult{
		SentimentScore: sentiment.DocumentSentiment.Score,
		SentimentLabel: sentimentLabel,
	}

	return n.OutSentiment.Set(ctx, result)
}

func (n *AnalyzeSentiment) OnClose() error {
	return nil
}
