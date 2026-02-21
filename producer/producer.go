package producer

import (
	"encoding/csv"
	"fmt"
	"math/rand/v2"
	"os"

	"github.com/google/uuid"
)

type Producer struct {
	dir string
}

func New(dir string) *Producer {
	return &Producer{
		dir: dir,
	}
}

func (p *Producer) Produce() error {
	guid := uuid.New()

	rows := p.generateRows()

	file, err := os.Create(fmt.Sprintf("%s/%s.tsv", p.dir, guid))
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = '\t'
	defer writer.Flush()

	header := p.generateHeader()
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	if err := writer.WriteAll(rows); err != nil {
		return fmt.Errorf("failed to write rows: %w", err)
	}
	return nil
}

func (p *Producer) generateHeader() []string {
	return []string{
		"n", "mqtt", "invid", "unit_guid", "msg_id", "text",
		"context", "class", "level", "area", "addr", "block",
		"type", "bit", "invert_bit"}
}

func (p *Producer) generateRows() [][]string {
	numsRows := rand.IntN(10) + 1

	guids := p.generateGuid(numsRows)

	text := []string{"Разморозка", "Вентилятор", "Высокая температура", "Компрессор", "Охлаждение"}
	class := []string{"alarm", "warning", "info", "event", "comand"}

	var rows [][]string

	for i := 0; i < numsRows; i++ {
		row := []string{
			fmt.Sprint(i + 1),
			"",
			fmt.Sprintf("G-044%d", 300+rand.IntN(20)),
			guids[rand.IntN(len(guids))],
			"cold7_Test",
			text[rand.IntN(5)],
			"",
			class[rand.IntN(5)],
			fmt.Sprint(rand.IntN(101)),
			"LOCAL",
			"cold7_status.Test",
			"",
			"",
			"",
			"",
		}
		rows = append(rows, row)
	}
	return rows
}

func (p *Producer) generateGuid(nums int) []string {
	numsGuids := rand.IntN(nums) + 1
	guids := make([]string, numsGuids)

	for i := 0; i < numsGuids; i++ {
		guids[i] = uuid.NewString()
	}
	return guids
}
