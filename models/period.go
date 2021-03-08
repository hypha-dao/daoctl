package models

import (
	"context"
	"fmt"
	"time"

	eos "github.com/eoscanada/eos-go"
	"github.com/hypha-dao/document-graph/docgraph"
	"go.uber.org/zap"
)

// Period represents a period of time aligning to a payroll period, typically a week
type Period struct {
	Label          string
	StartTimePoint eos.TimePoint
	StartTime      time.Time
	Next           *Period
	Document       docgraph.Document
}

func NewPeriod(ctx context.Context, api *eos.API, contract eos.AccountName, doc docgraph.Document) (Period, error) {
	p := Period{}
	p.Document = doc

	startTime, err := doc.GetContentFromGroup("details", "start_time")
	if err != nil {
		return Period{}, fmt.Errorf("get content failed: %v", err)
	}
	p.StartTimePoint, err = startTime.TimePoint()
	if err != nil {
		return Period{}, fmt.Errorf("get content failed: %v", err)
	}
	p.StartTime = time.Unix(int64(p.StartTimePoint)/1000000, 0).UTC()
	zap.S().Debugf("Loading a period with a start date: %v", p.StartTime.Format("2006 Jan 02 15:04:05"))

	label, err := doc.GetContentFromGroup("details", "label")
	if err != nil {
		return Period{}, fmt.Errorf("get content failed: %v", err)
	}
	p.Label = label.String()

	nextEdges, err := docgraph.GetEdgesFromDocumentWithEdge(ctx, api, contract, doc, eos.Name("next"))
	if err != nil {
		return Period{}, fmt.Errorf("error while retrieving next edge: %v", err)
	}
	if len(nextEdges) == 0 {
		zap.S().Debugf("There is no edge period, returning: %v", p.Document.Hash.String())
		p.Next = nil
		return p, nil
	} else {
		zap.S().Debugf("Loading the next period as: %v", nextEdges[0].ToNode.String())

		nextDocument, err := docgraph.LoadDocument(ctx, api, contract, nextEdges[0].ToNode.String())
		if err != nil {
			return Period{}, fmt.Errorf("unable to load next edge: %v", err)
		}
		nextPeriod, err := NewPeriod(ctx, api, contract, nextDocument)
		if err != nil {
			return Period{}, fmt.Errorf("unable to create next Period: %v", err)
		}
		p.Next = &nextPeriod
	}
	return p, nil
}
