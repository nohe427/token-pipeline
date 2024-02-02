package onwrite

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/compute/metadata"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/googleapis/google-cloudevents-go/cloud/firestoredata"
	"github.com/nohe427/token-pipeline/adc"
	"github.com/nohe427/token-pipeline/formatter"
	"github.com/nohe427/token-pipeline/vertexhelp"
	"google.golang.org/protobuf/proto"
)

type UserFields struct {
	InputText            string
	OutputText           string
	InputTextTokenCount  string
	OutputTextTokenCount string
}

type CountTokenOpts struct {
	Location  string
	ProjectID string
	Model     string
	Token     string
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

func countTokens(input string, opts CountTokenOpts) (int, error) {
	ctp := vertexhelp.NewCountTokenParams(opts.Location, opts.ProjectID, opts.Model)
	ctReq := vertexhelp.CountTokenRequest{Instances: []vertexhelp.Prompt{{Prompt: input}}}
	return vertexhelp.RequestTokenCount(ctp, &ctReq, opts.Token)
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

	adcToken, err := adc.GetADCToken()
	if err != nil {
		return err
	}

	projId := metadata.ProjectID()

	opts := CountTokenOpts{Location: vertexhelp.DEFAULT_LOCATION, ProjectID: projId, Model: vertexhelp.DEFAULT_MODEL, Token: adcToken}
	countTokens(fmtInputStr, opts)

	return nil
}
