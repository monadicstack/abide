package testext

import (
	"bytes"
	"context"
	"io"
	"strings"
	"time"

	"github.com/monadicstack/abide/fail"
	"github.com/monadicstack/abide/metadata"
)

type SampleServiceHandler struct {
	Sequence *Sequence
}

func (s SampleServiceHandler) Defaults(_ context.Context, req *SampleRequest) (*SampleResponse, error) {
	s.Sequence.Append("Defaults:" + req.Text)
	return &SampleResponse{Text: "Defaults:" + req.Text}, nil
}

func (s SampleServiceHandler) ComplexValues(_ context.Context, req *SampleComplexRequest) (*SampleComplexResponse, error) {
	s.Sequence.Append("ComplexValues:" + req.InUser.ID)
	res := SampleComplexResponse{
		OutFlag:    req.InFlag,
		OutFloat:   req.InFloat,
		OutTime:    req.InTime,
		OutTimePtr: req.InTimePtr,
		OutUser: &SampleUser{
			ID:              req.InUser.ID,
			Name:            req.InUser.Name,
			Age:             req.InUser.Age,
			Attention:       req.InUser.Attention,
			AttentionString: req.InUser.AttentionString,
			PhoneNumber:     req.InUser.PhoneNumber,
			MarshalToString: req.InUser.MarshalToString,
			MarshalToObject: req.InUser.MarshalToObject,
		},
	}
	return &res, nil
}

func (s SampleServiceHandler) ComplexValuesPath(_ context.Context, req *SampleComplexRequest) (*SampleComplexResponse, error) {
	s.Sequence.Append("ComplexValuesPath:" + req.InUser.ID)
	res := SampleComplexResponse{
		OutFlag:    req.InFlag,
		OutFloat:   req.InFloat,
		OutTime:    req.InTime,
		OutTimePtr: req.InTimePtr,
		OutUser: &SampleUser{
			ID:              req.InUser.ID,
			Name:            req.InUser.Name,
			Age:             req.InUser.Age,
			Attention:       req.InUser.Attention,
			AttentionString: req.InUser.AttentionString,
			PhoneNumber:     req.InUser.PhoneNumber,
			MarshalToString: req.InUser.MarshalToString,
			MarshalToObject: req.InUser.MarshalToObject,
		},
	}
	return &res, nil
}

func (s SampleServiceHandler) Fail4XX(_ context.Context, req *SampleRequest) (*SampleResponse, error) {
	s.Sequence.Append("Fail4XX:" + req.Text)
	return nil, fail.AlreadyExists("always a conflict")
}

func (s SampleServiceHandler) Fail5XX(_ context.Context, req *SampleRequest) (*SampleResponse, error) {
	s.Sequence.Append("Fail5XX:" + req.Text)
	return nil, fail.BadGateway("always a bad gateway")
}

func (s SampleServiceHandler) CustomRoute(_ context.Context, req *SampleRequest) (*SampleResponse, error) {
	s.Sequence.Append("CustomRoute:" + req.Text)
	return &SampleResponse{ID: req.ID, Text: "Route:" + req.Text}, nil
}

func (s SampleServiceHandler) CustomRouteQuery(_ context.Context, req *SampleRequest) (*SampleResponse, error) {
	s.Sequence.Append("CustomRouteQuery:" + req.Text)
	return &SampleResponse{ID: req.ID, Text: "Route:" + req.Text}, nil
}

func (s SampleServiceHandler) CustomRouteBody(_ context.Context, req *SampleRequest) (*SampleResponse, error) {
	s.Sequence.Append("CustomRouteBody:" + req.Text)
	return &SampleResponse{ID: req.ID, Text: "Route:" + req.Text}, nil
}

func (s SampleServiceHandler) OmitMe(_ context.Context, req *SampleRequest) (*SampleResponse, error) {
	s.Sequence.Append("OmitMe:" + req.Text)
	return &SampleResponse{Text: "Doesn't matter..."}, nil
}

func (s SampleServiceHandler) Download(_ context.Context, req *SampleDownloadRequest) (*SampleDownloadResponse, error) {
	s.Sequence.Append("Download:" + req.Format)
	res := SampleDownloadResponse{}
	switch req.Format {
	case "text/csv":
		content := "ID,Name,Enabled\n1,Dude,true\n2,Walter,false"
		res.SetContent(io.NopCloser(bytes.NewBufferString(content)))
		res.SetContentType("text/csv")
		res.SetContentLength(len([]byte(content)))
		res.SetContentFileName("dude.csv")
	default:
		content := "Donny, you're out of your element!"
		res.SetContent(io.NopCloser(bytes.NewBufferString(content)))
		res.SetContentType("text/plain")
		res.SetContentLength(len([]byte(content)))
		res.SetContentFileName("dude.txt")
	}
	return &res, nil
}

func (s SampleServiceHandler) DownloadResumable(_ context.Context, req *SampleDownloadRequest) (*SampleDownloadResponse, error) {
	s.Sequence.Append("DownloadResumable:" + req.Format)
	content := "<h1>The Dude Abides</h1>"
	res := SampleDownloadResponse{}
	res.SetContentType("text/html")
	res.SetContentRange(50, 50+len(content), 1024)
	res.SetContent(io.NopCloser(bytes.NewBufferString(content)))
	res.SetContentFileName("dude.html")
	return &res, nil
}

func (s SampleServiceHandler) Redirect(_ context.Context, _ *SampleRedirectRequest) (*SampleRedirectResponse, error) {
	s.Sequence.Append("Redirect")
	return &SampleRedirectResponse{URI: "/v2/download?Format=text/csv"}, nil
}

func (s SampleServiceHandler) Authorization(ctx context.Context, req *SampleRequest) (*SampleResponse, error) {
	s.Sequence.Append("Authorization:" + req.Text)
	return &SampleResponse{Text: metadata.Authorization(ctx)}, nil
}

func (s SampleServiceHandler) Sleep(_ context.Context, req *SampleRequest) (*SampleResponse, error) {
	s.Sequence.Append("Sleep:" + req.Text)
	time.Sleep(5 * time.Second)
	return &SampleResponse{Text: "The Dude Abides"}, nil
}

func (s SampleServiceHandler) TriggerUpperCase(_ context.Context, req *SampleRequest) (*SampleResponse, error) {
	s.Sequence.Append("TriggerUpperCase:" + req.Text)
	return &SampleResponse{Text: strings.ToUpper(req.Text)}, nil
}

func (s SampleServiceHandler) TriggerLowerCase(_ context.Context, req *SampleRequest) (*SampleResponse, error) {
	s.Sequence.Append("TriggerLowerCase:" + req.Text)
	return &SampleResponse{Text: strings.ToLower(req.Text)}, nil
}

func (s SampleServiceHandler) TriggerFailure(_ context.Context, req *SampleRequest) (*SampleResponse, error) {
	s.Sequence.Append("TriggerFailure:" + req.Text)
	return nil, fail.Unexpected("nope...")
}

func (s SampleServiceHandler) ListenerA(_ context.Context, req *SampleRequest) (*SampleResponse, error) {
	s.Sequence.Append("ListenerA:" + req.Text)
	return &SampleResponse{Text: "ListenerA:" + req.Text}, nil
}

func (s SampleServiceHandler) ListenerB(_ context.Context, req *SampleRequest) (*SampleResponse, error) {
	s.Sequence.Append("ListenerB:" + req.Text)
	return &SampleResponse{Text: "ListenerB:" + req.Text}, nil
}

func (s SampleServiceHandler) SecureWithRoles(ctx context.Context, req *SampleSecurityRequest) (*SampleSecurityResponse, error) {
	s.Sequence.Append("SecureWithRoles:" + req.ID)
	return &SampleSecurityResponse{Roles: metadata.Route(ctx).Roles}, nil
}

func (s SampleServiceHandler) Panic(_ context.Context, req *SampleRequest) (*SampleResponse, error) {
	panic("don't")
}
