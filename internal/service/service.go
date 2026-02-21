package service

import (
	"context"
	"device-parser-logs/internal/models"
	"device-parser-logs/pkg/errorx"
	"errors"
	"log/slog"
)

type parser interface {
	ParseTSV(dir string) ([]*models.DeviceLogs, error)
}

type repository interface {
	SaveDeviceLogs(context.Context, []*models.DeviceLogs) error
	GetInfoByGuid(context.Context, string, int, int) (*models.PaginationResult, error)
	GetProssecedFile(context.Context, string) (*models.ProcessedFile, error)
	SaveProssecedFile(context.Context, *models.ProcessedFile) error
}

type generator interface {
	Generate(string, []*models.DeviceLogs) error
}

type Service struct {
	parser     parser
	repository repository
	generator  generator
	logger     *slog.Logger
}

func New(p parser, r repository, g generator, logger *slog.Logger) *Service {
	return &Service{
		parser:     p,
		repository: r,
		generator:  g,
		logger:     logger,
	}
}

func (s *Service) SaveFileTSV(ctx context.Context, filePath string) error {
	res, err := s.parser.ParseTSV(filePath)
	if err != nil {
		if errors.Is(err, errorx.ErrParseNotInt) {
			s.logger.Warn("parser: validation error", "filename", filePath, "err", err)
		} else {
			s.logger.Error("failed to parse tsv file", "file", filePath, "error", err.Error())
			return err
		}
	}

	err = s.repository.SaveDeviceLogs(ctx, res)
	if err != nil {
		s.logger.Error("failed to save tsv file", "file", filePath, "error", err.Error())
		return err
	}

	guids := s.CheckGuidsInfile(res)

	for guid, deviceLog := range guids {
		err = s.generator.Generate(guid, deviceLog)
		if err != nil {
			s.logger.Error("failed to generate pdf file", "guid", guid, "error", err.Error())
		}
	}

	return err
}

func (s *Service) CheckGuidsInfile(res []*models.DeviceLogs) map[string][]*models.DeviceLogs {
	mapGuid := make(map[string][]*models.DeviceLogs)
	for _, deviceLog := range res {
		mapGuid[deviceLog.Guid] = append(mapGuid[deviceLog.Guid], deviceLog)
	}

	return mapGuid
}

func (s *Service) GetDeviceLogs(ctx context.Context, guid string, page, limit int) (*models.PaginationResult, error) {
	res, err := s.repository.GetInfoByGuid(ctx, guid, page, limit)
	if err != nil {
		s.logger.Error("failed to get devicelogs", "guid", guid, "error", err.Error())
		return nil, err
	}

	return res, nil
}

func (s *Service) GetProssedFile(ctx context.Context, fileName string) (*models.ProcessedFile, error) {
	file, err := s.repository.GetProssecedFile(ctx, fileName)
	if err != nil {
		s.logger.Error("failed to get prossedfile", "file", fileName, "error", err.Error())
		return nil, err
	}
	return file, err
}

func (s *Service) SaveProssedFile(ctx context.Context, file *models.ProcessedFile) {
	err := s.repository.SaveProssecedFile(ctx, file)
	if err != nil {
		s.logger.Error("failed to save prossedfile", "file", file.Name, "error", err.Error())
	}
}
