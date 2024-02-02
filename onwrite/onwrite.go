package onwrite

import (
	"context"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/googleapis/google-cloudevents-go/cloud/firestoredata"
	"google.golang.org/protobuf/proto"
)

type UserFields struct {
	InputText            string
	OutputText           string
	InputTextTokenCount  string
	OutputTextTokenCount string
}

func init() {
	functions.CloudEvent("DocumentUpdates", docUpdated)
}

func protoEventToFirestoreEvent(e event.Event) (*firestoredata.DocumentEventData, error) {
	ed := firestoredata.DocumentEventData{}

	if err := proto.Unmarshal(e.Data(), &ed); err != nil {
		return &ed, err
	}

	return &ed, nil
}

func configureUserFields(uf *UserFields) {
	if userIf := os.Getenv("inputText"); userIf != "" {
		uf.InputText = userIf
	}
	if userOf := os.Getenv("outputText"); userOf != "" {
		uf.OutputText = userOf
	}
	if userOfTc := os.Getenv("outputTextTokenCount"); userOfTc != "" {
		uf.OutputTextTokenCount = userOfTc
	}
	if userIfTc := os.Getenv("inputTextTokenCount"); userIfTc != "" {
		uf.InputTextTokenCount = userIfTc
	}
}

func countTokens(input string) error {

	return nil
}

func docUpdated(ctx context.Context, e event.Event) error {
	uf := UserFields{InputText: "inputText", OutputText: "outputText", InputTextTokenCount: "inputTextTokenCount", OutputTextTokenCount: "outputTextTokenCount"}
	configureUserFields(&uf)

	return nil
}
