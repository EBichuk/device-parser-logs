package watcher

import (
	"context"
	"device-parser-logs/internal/models"
	"fmt"
	"log/slog"
	"os"

	"time"
)

type service interface {
	SaveFileTSV(context.Context, string) error
	GetProssedFile(context.Context, string) (*models.ProcessedFile, error)
	SaveProssedFile(context.Context, *models.ProcessedFile)
}

type Pool struct {
	ctx context.Context

	workersNum   int
	directoryTsv string
	interval     string

	taskQueueCh chan string

	stopChan chan struct{}
	service  service
	logger   *slog.Logger
}

func NewPool(ctx context.Context, workersNum int, dirTsv string, interval string, logger *slog.Logger, service service) *Pool {
	return &Pool{
		ctx:          ctx,
		workersNum:   workersNum,
		directoryTsv: dirTsv,
		interval:     interval,
		taskQueueCh:  make(chan string),
		stopChan:     make(chan struct{}),
		logger:       logger,
		service:      service,
	}
}

func (p *Pool) RunPool() {
	go p.Dirmon()

	for workerId := range p.workersNum {
		go p.WorkerRun(workerId)
	}
}

func (p *Pool) WorkerRun(id int) {
	p.logger.Info("worker running", "id", id)
	for {
		select {
		case <-p.ctx.Done():
			return
		case task, ok := <-p.taskQueueCh:
			if ok {
				err := p.service.SaveFileTSV(p.ctx, task)

				file := models.ProcessedFile{
					Name:        task,
					ProcessedAt: time.Now(),
					Status:      "success",
					MsgError:    "",
				}

				if err != nil {
					file.Status = "error"
					file.MsgError = err.Error()
				}

				p.service.SaveProssedFile(p.ctx, &file)
			} else {
				p.logger.Info("worker stopped", "id", id)
				return
			}
		}
	}
}

func (p *Pool) Dirmon() error {
	dur, err := time.ParseDuration(p.interval)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(dur)
	defer ticker.Stop()

	err = p.ScanDirectory()
	if err != nil {
		p.logger.Error("failed to scan directory", "error", err.Error())
	}

	for {
		select {
		case <-p.ctx.Done():
			return nil
		case <-p.stopChan:
			p.logger.Info("directory monitoring stopped")
			return nil
		case <-ticker.C:
			if err := p.ScanDirectory(); err != nil {
				p.logger.Error("failed to scan directory", "error", err.Error())
			}
		}
	}
}

func (p *Pool) ScanDirectory() error {
	files, err := os.ReadDir(p.directoryTsv)
	if err != nil {
		return fmt.Errorf("failed to read directory %w", err)
	}

	for _, file := range files {
		select {
		case <-p.ctx.Done():
			return nil
		case <-p.stopChan:
			return nil
		default:
		}
		fileExistence, err := p.service.GetProssedFile(p.ctx, file.Name())

		if err == nil && fileExistence == nil {
			p.taskQueueCh <- file.Name()
		}
		if err != nil {
			p.logger.Error("error to find prossed file in db", "error", err.Error())
		}
	}
	return nil
}

func (p *Pool) Stop() {
	p.stopChan <- struct{}{}
	close(p.stopChan)
	close(p.taskQueueCh)
}
