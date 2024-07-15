package response_models

import (
	"go-api/pkg/model"
	"time"
)

type GetReportResponse struct {
	ID          uint
	Description string
	CreatedAt   *time.Time
	CreatedBy   GetUserResponseInterface
	Status      int
	Drop        GetDropResponse              `json:",omitempty"`
	Comment     GetCommentResponse           `json:",omitempty"`
	Response    []GetCommentResponseResponse `json:",omitempty"`
}

func FormatGetReportResponse(report model.ReportModel) GetReportResponse {
	if nil == report {
		return GetReportResponse{}
	}

	createdAt := time.Unix(int64(report.GetCreatedAt()), 0)

	return GetReportResponse{
		ID:          report.GetID(),
		Description: report.GetDescription(),
		CreatedAt:   &createdAt,
		CreatedBy:   FormatGetUserResponse(report.GetCreatedBy()),
		Status:      report.GetStatus(),
		Drop:        FormatGetDropResponse(report.GetReportedDrop(), false),
		Comment:     FormatGetCommentResponse(report.GetReportedComment()),
		Response:    FormatGetCommentResponsesResponse(report.GetReportedResponse()),
	}
}
