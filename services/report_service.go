package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type ReportService struct {
	repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) GetReportToday() (*models.Report, error) {
	return s.repo.GetReportToday()
}

func (s *ReportService) GetReportByDate(startDate string, endDate string) (*models.Report, error) {
	return s.repo.GetReportByDate(startDate, endDate)
}
