package postgres

import (
	"go-api/pkg/model"
	"gorm.io/gorm"
)

type Report struct {
	gorm.Model
	Description        string `gorm:"not null"`
	Status             int    `gorm:"not null"`
	CreatedById        uint   `gorm:"not null"`
	ReportedDropId     uint
	ReportedCommentId  uint
	ReportedResponseId uint
	CreatedBy          User              `gorm:"foreignKey:CreatedById;references:ID"`
	ReportedDrop       Drop              `gorm:"foreignKey:ReportedDropId;references:ID"`
	ReportedComment    Comment           `gorm:"foreignKey:ReportedCommentId;references:ID"`
	ReportedResponse   []CommentResponse `gorm:"foreignKey:ReportedResponseId;references:ID"`
}

func (r *Report) GetID() uint {
	return r.ID
}

func (r *Report) GetDescription() string {
	return r.Description
}

func (r *Report) GetStatus() int {
	return r.Status
}

func (r *Report) GetCreatedBy() model.UserModel {
	return &r.CreatedBy
}

func (r *Report) GetCreatedAt() int {
	return int(r.CreatedAt.Unix())
}

func (r *Report) GetReportedDrop() model.DropModel {
	return &r.ReportedDrop
}

func (r *Report) GetReportedComment() model.CommentModel {
	return &r.ReportedComment
}

func (r *Report) GetReportedResponse() []model.CommentResponseModel {
	var result []model.CommentResponseModel
	for _, response := range r.ReportedResponse {
		result = append(result, &response)
	}
	return result

}

type ReportStatusActive struct{}

func (r *ReportStatusActive) ToInt() int {
	return 1
}

type ReportStatusResolved struct{}

func (r *ReportStatusResolved) ToInt() int {
	return 2
}

type ReportStatusDeleted struct{}

func (r *ReportStatusDeleted) ToInt() int {
	return -1
}

type repoReportPrivate struct {
	db *gorm.DB
}

func NewReportRepo(db *gorm.DB) model.ReportRepository {
	return &repoReportPrivate{db: db}
}

func (r *repoReportPrivate) CreateReport(
	description string,
	createdById uint,
	reportedDropId uint,
	reportedCommentId uint,
	reportedResponseId uint,
) (model.ReportModel, error) {
	var createdByUser User
	if err := r.db.First(&createdByUser, createdById).Error; err != nil {
		return nil, err
	}

	var reportedDrop Drop
	if reportedDropId != 0 {
		if err := r.db.First(&reportedDrop, reportedDropId).Error; err != nil {
			return nil, err
		}
	}

	var reportedComment Comment
	if reportedCommentId != 0 {
		if err := r.db.First(&reportedComment, reportedCommentId).Error; err != nil {
			return nil, err
		}
	}

	var reportedResponse []CommentResponse
	if reportedResponseId != 0 {
		if err := r.db.First(&reportedResponse, reportedResponseId).Error; err != nil {
			return nil, err
		}
	}

	status := new(ReportStatusActive).ToInt()

	report := Report{
		Description:      description,
		Status:           status,
		CreatedBy:        createdByUser,
		ReportedDrop:     reportedDrop,
		ReportedComment:  reportedComment,
		ReportedResponse: reportedResponse,
	}
	if err := r.db.Create(&report).Error; err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *repoReportPrivate) DeleteReport(reportId uint) error {
	return r.db.Delete(&Report{}, reportId).Error
}

func (r *repoReportPrivate) GetReportsByDropId(dropId uint) ([]model.ReportModel, error) {
	var reports []Report
	if err := r.db.Preload("CreatedBy").Preload("ReportedDrop").Preload("ReportedComment").Preload("ReportedResponse").Where("reported_drop_id = ?", dropId).Find(&reports).Error; err != nil {
		return nil, err
	}
	var result []model.ReportModel
	for _, report := range reports {
		result = append(result, &report)
	}
	return result, nil
}

func (r *repoReportPrivate) GetReportsByCommentId(commentId uint) ([]model.ReportModel, error) {
	var reports []Report
	if err := r.db.Preload("CreatedBy").Preload("ReportedDrop").Preload("ReportedComment").Preload("ReportedResponse").Where("reported_comment_id = ?", commentId).Find(&reports).Error; err != nil {
		return nil, err
	}
	var result []model.ReportModel
	for _, report := range reports {
		result = append(result, &report)
	}
	return result, nil
}

func (r *repoReportPrivate) GetReportsByUserId(userId uint) ([]model.ReportModel, error) {
	var reports []Report
	if err := r.db.Preload("CreatedBy").Preload("ReportedDrop").Preload("ReportedComment").Preload("ReportedResponse").Where("created_by_id = ?", userId).Find(&reports).Error; err != nil {
		return nil, err
	}
	var result []model.ReportModel
	for _, report := range reports {
		result = append(result, &report)
	}
	return result, nil
}

func (r *repoReportPrivate) GetReportById(reportId uint) (model.ReportModel, error) {
	var report Report
	if err := r.db.Preload("CreatedBy").Preload("ReportedDrop").Preload("ReportedComment").Preload("ReportedResponse").First(&report, reportId).Error; err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *repoReportPrivate) GetAllReports() ([]model.ReportModel, error) {
	var reports []Report
	if err := r.db.Preload("CreatedBy").Preload("ReportedDrop").Preload("ReportedComment").Preload("ReportedResponse").Find(&reports).Error; err != nil {
		return nil, err
	}
	var result []model.ReportModel
	for _, report := range reports {
		result = append(result, &report)
	}
	return result, nil
}

func (r *repoReportPrivate) UpdateReportStatus(reportId uint, status uint) error {
	return r.db.Model(&Report{}).Where("id = ?", reportId).Update("status", status).Error
}

func (r *repoReportPrivate) GetActiveReportByDropAndUser(dropId uint, userId uint) (model.ReportModel, error) {
	var report Report
	if err := r.db.Preload("CreatedBy").Preload("ReportedDrop").Preload("ReportedComment").Preload("ReportedResponse").Where("reported_drop_id = ? AND created_by_id = ? AND status = 1", dropId, userId).First(&report).Error; err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *repoReportPrivate) GetActiveReportByCommentAndUser(commentId uint, userId uint) (model.ReportModel, error) {
	var report Report
	if err := r.db.Preload("CreatedBy").Preload("ReportedDrop").Preload("ReportedComment").Preload("ReportedResponse").Where("reported_comment_id = ? AND created_by_id = ? AND status = 1", commentId, userId).First(&report).Error; err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *repoReportPrivate) GetActiveReportByResponseAndUser(responseId uint, userId uint) (model.ReportModel, error) {
	var report Report
	if err := r.db.Preload("CreatedBy").Preload("ReportedDrop").Preload("ReportedComment").Preload("ReportedResponse").Where("reported_response_id = ? AND created_by_id = ? AND status = 1", responseId, userId).First(&report).Error; err != nil {
		return nil, err
	}
	return &report, nil
}
