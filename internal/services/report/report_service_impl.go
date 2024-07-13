package report

import (
	"go-api/internal/repositories"
	"go-api/pkg/errors2"
	"go-api/pkg/model"
)

type ReportService struct {
	Repo *repositories.Repositories
}

func (s *ReportService) ReportDrop(userId uint, dropId uint, description string) (model.ReportModel, error) {
	if _, err := s.Repo.DropRepository.GetDropById(dropId); err != nil {
		return nil, errors2.NotFoundError{Entity: "Drop"}
	}

	if _, err := s.Repo.ReportRepository.GetActiveReportByDropAndUser(dropId, userId); err == nil {
		return nil, errors2.AlreadyReportedError{Entity: "Drop"}
	}

	report, err := s.Repo.ReportRepository.CreateReport(description, userId, dropId, 0, 0)
	if err != nil {
		return nil, err
	}

	return report, nil
}

func (s *ReportService) ReportComment(userId uint, commentId uint, description string) (model.ReportModel, error) {
	if _, err := s.Repo.CommentRepository.GetById(commentId); err != nil {
		return nil, errors2.NotFoundError{Entity: "Comment"}
	}

	if _, err := s.Repo.ReportRepository.GetActiveReportByCommentAndUser(commentId, userId); err == nil {
		return nil, errors2.AlreadyReportedError{Entity: "Comment"}
	}

	report, err := s.Repo.ReportRepository.CreateReport(description, userId, 0, commentId, 0)
	if err != nil {
		return nil, err
	}

	return report, nil
}

func (s *ReportService) ReportResponse(userId uint, responseId uint, description string) (model.ReportModel, error) {
	if _, err := s.Repo.CommentResponseRepository.GetById(responseId); err != nil {
		return nil, errors2.NotFoundError{Entity: "Response"}
	}

	if _, err := s.Repo.ReportRepository.GetActiveReportByResponseAndUser(responseId, userId); err == nil {
		return nil, errors2.AlreadyReportedError{Entity: "Response"}
	}

	report, err := s.Repo.ReportRepository.CreateReport(description, userId, 0, 0, responseId)
	if err != nil {
		return nil, err
	}

	return report, nil
}

func (s *ReportService) DeleteReport(reportId uint) error {
	if err := s.Repo.ReportRepository.DeleteReport(reportId); err != nil {
		return err
	}

	return nil
}
