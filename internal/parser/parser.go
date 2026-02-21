package parser

import (
	"device-parser-logs/internal/models"
	"device-parser-logs/pkg/errorx"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const constValidationErrorValue = -1

type Parser struct {
	dir string
}

func New(dir string) *Parser {
	return &Parser{
		dir: dir,
	}
}

func (p *Parser) ParseTSV(dir string) ([]*models.DeviceLogs, error) {
	var logsInfo []*models.DeviceLogs

	fileName := fmt.Sprintf("%s/%s", p.dir, dir)

	file, err := os.Open(fileName)
	if err != nil {
		return logsInfo, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	r := csv.NewReader(file)
	r.Comma = '\t'
	r.FieldsPerRecord = -1

	_, err = r.Read()
	if errors.Is(err, io.EOF) {
		return logsInfo, fmt.Errorf("empty file %w", err)
	}

	for {
		record, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return logsInfo, fmt.Errorf("failed to read row %w", err)
		}

		dto, err := ToDeviceLogs(record)
		if err != nil {
			continue
		}

		logsInfo = append(logsInfo, dto)
	}

	return logsInfo, err
}

func ToDeviceLogs(log []string) (*models.DeviceLogs, error) {
	for i, s := range log {
		log[i] = strings.TrimSpace(s)
	}

	var err error

	level, err := strconv.Atoi(log[8])
	if err != nil {
		level = constValidationErrorValue
		err = errorx.ErrParseNotInt
	}

	return &models.DeviceLogs{
		ID:        primitive.NewObjectID(),
		Mqtt:      log[1],
		Invid:     log[2],
		Guid:      log[3],
		MsgId:     log[4],
		Text:      log[5],
		Context:   log[6],
		ClassMsg:  log[7],
		Level:     level,
		Area:      log[9],
		Addr:      log[10],
		Block:     log[11],
		Type:      log[12],
		Bit:       log[13],
		InvertBit: log[14],
		CreateAt:  time.Now(),
	}, err
}
