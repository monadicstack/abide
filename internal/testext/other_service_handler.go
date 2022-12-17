package testext

import (
	"context"
	"fmt"
	"strings"
)

type OtherServiceHandler struct {
	Sequence      *Sequence
	SampleService SampleService
}

func (svc OtherServiceHandler) SpaceOut(_ context.Context, req *OtherRequest) (*OtherResponse, error) {
	svc.Sequence.Append("SpaceOut:" + req.Text)
	runes := strings.Split(req.Text, "")
	return &OtherResponse{Text: strings.Join(runes, " ")}, nil
}

func (svc OtherServiceHandler) RPCExample(ctx context.Context, req *OtherRequest) (*OtherResponse, error) {
	svc.Sequence.Append("RPCExample:" + req.Text)
	res, err := svc.SampleService.TriggerUpperCase(ctx, &SampleRequest{Text: req.Text})
	if err != nil {
		return nil, fmt.Errorf("wtf: %w", err)
	}
	return &OtherResponse{Text: res.Text}, nil
}

func (svc OtherServiceHandler) ListenWell(_ context.Context, req *OtherRequest) (*OtherResponse, error) {
	svc.Sequence.Append("ListenWell:" + req.Text)
	return &OtherResponse{Text: "ListenWell:" + req.Text}, nil
}
