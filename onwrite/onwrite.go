package onwrite

import (
	"context"
	"fmt"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/googleapis/google-cloudevents-go/cloud/firestoredata"
	"github.com/nohe427/token-pipeline/formatter"
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

func countTokens(input string) (int, error) {

	return 0, nil
}

func docUpdated(ctx context.Context, e event.Event) error {
	uf := UserFields{InputText: "inputText", OutputText: "outputText", InputTextTokenCount: "inputTextTokenCount", OutputTextTokenCount: "outputTextTokenCount"}
	configureUserFields(&uf)

	fe, err := protoEventToFirestoreEvent(e)
	if err != nil {
		return err
	}
	inputStr := fe.GetValue().GetFields()[uf.InputText].GetStringValue()
	fmtInputStr, err := formatter.FormatInput(inputStr)
	if err != nil {
		return err
	}

	outputStr := fe.GetValue().GetFields()[uf.OutputText].GetStringValue()
	fmtOutputStr, err := formatter.FormatInput(outputStr)
	if err != nil {
		return err
	}

	fmt.Println(fmtInputStr, fmtOutputStr)

	return nil
}
