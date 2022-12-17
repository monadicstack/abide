package fieldtypes

import (
	"context"
	"fmt"
	"time"

	"github.com/dimfeld/httptreemux/v5"
	"github.com/monadicstack/abide/parser/testdata"
)

type HappyLittleService interface {
	PaintTree(context.Context, *Request) (*Response, error)
}

type Request struct {
	EmbeddedFields
	EmbeddedString
	embeddedBool

	notExported string

	Basic        string
	BasicPointer *string

	ExportedStruct        ExportedStruct
	ExportedStructPointer *ExportedStruct

	NotExportedStruct        notExportedStruct
	NotExportedStructPointer *notExportedStruct

	Time        time.Time
	TimePointer *time.Time

	Duration        time.Duration
	DurationPointer *time.Duration

	Interface any
	Stringer  fmt.Stringer

	BasicSlice []string
	BasicMap   map[string]string

	AliasBasic        AliasBasic
	AliasBasicPointer *AliasBasic

	AliasStruct        AliasStruct
	AliasStructPointer *AliasStruct

	AliasSlice        AliasSlice
	AliasSlicePointer *AliasSlice

	ThirdParty httptreemux.ContextRouteData
	SharedType testdata.SharedType
}

type Response struct {
}

type ExportedStruct struct {
	Name string
}

type notExportedStruct struct {
	name string
}

type EmbeddedFields struct {
	EmbeddedA string
	EmbeddedB bool
	EmbeddedC ExportedStruct
	embeddedD int
}

type EmbeddedString string
type embeddedBool bool

type AliasBasic uint64
type AliasStruct ExportedStruct
type AliasSlice []string
