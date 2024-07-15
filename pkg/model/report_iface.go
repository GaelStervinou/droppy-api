package model

type ReportModel interface {
	GetID() uint
	GetDescription() string
	GetStatus() int
	GetCreatedBy() UserModel
	GetCreatedAt() int
	GetReportedDrop() DropModel
	GetReportedComment() CommentModel
	GetReportedResponse() []CommentResponseModel
}

type ReportRepository interface {
	CreateReport(description string, createdBy uint, reportedDrop uint, reportedComment uint, reportedResponse uint) (ReportModel, error)
	DeleteReport(reportId uint) error
	GetReportsByDropId(dropId uint) ([]ReportModel, error)
	GetReportsByCommentId(commentId uint) ([]ReportModel, error)
	GetReportsByUserId(userId uint) ([]ReportModel, error)
	GetReportById(reportId uint) (ReportModel, error)
	GetAllReports() ([]ReportModel, error)
	UpdateReportStatus(reportId uint, status uint) error
	GetActiveReportByDropAndUser(dropId uint, userId uint) (ReportModel, error)
	GetActiveReportByCommentAndUser(commentId uint, userId uint) (ReportModel, error)
	GetActiveReportByResponseAndUser(responseId uint, userId uint) (ReportModel, error)
}

type ReportService interface {
	ReportDrop(userId uint, dropId uint, args ReportCreationParam) (ReportModel, error)
	ReportComment(userId uint, commentId uint, args ReportCreationParam) (ReportModel, error)
	ReportResponse(userId uint, responseId uint, args ReportCreationParam) (ReportModel, error)
	DeleteReport(reportId uint) error
}

type ReportCreationParam struct {
	Description string `json:"description"`
	DropId      uint   `json:"dropId"`
	CommentId   uint   `json:"commentId"`
	ResponseId  uint   `json:"responseId"`
}
