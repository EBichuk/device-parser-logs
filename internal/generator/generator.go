package generator

import (
	"device-parser-logs/internal/models"
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf/v2"
)

type GeneratorPdf struct {
	dir string
}

func New(dir string) *GeneratorPdf {
	return &GeneratorPdf{
		dir: dir,
	}
}

func (g *GeneratorPdf) Generate(guid string, data []*models.DeviceLogs) error {
	pdf := gofpdf.New("L", "mm", "A4", "")

	pdf.AddPage()

	fileName := fmt.Sprintf("%s/%s_%v.pdf", g.dir, guid, time.Now().Format("2006-01-02_15:04:05"))

	pdf.AddUTF8Font("DejaVuBold", "", "./pkg/font/DejaVuSans-Bold.ttf")
	pdf.SetFont("DejaVuBold", "", 8)

	pdf.CellFormat(277, 6, fmt.Sprintf("Total records: %d", len(data)), "", 1, "L", false, 0, "")
	pdf.Ln(3)
	pdf.SetFillColor(240, 240, 240)

	rowDeviceLogs := GetRowDeviceLogs()

	for _, header := range rowDeviceLogs {
		pdf.CellFormat(header.w, 8, header.header, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	pdf.AddUTF8Font("DejaVu", "", "./pkg/font/DejaVuSans.ttf")
	pdf.SetFont("DejaVu", "", 8)

	for i, d := range data {
		for j, table := range rowDeviceLogs {
			if j == 0 {
				pdf.CellFormat(table.w, 10, fmt.Sprint(i+1), "1", 0, "L", false, 0, "")
			} else {
				pdf.CellFormat(table.w, 10, table.getData(*d), "1", 0, "L", false, 0, "")
			}
		}
		pdf.Ln(-1)
	}
	pdf.Ln(-1)

	err := pdf.OutputFileAndClose(fileName)

	return err
}

type tableCell struct {
	header  string
	w       float64
	h       float64
	getData func(models.DeviceLogs) string
}

func GetRowDeviceLogs() []tableCell {
	return []tableCell{
		{
			header: "N",
			w:      6,
			h:      8,
			getData: func(models.DeviceLogs) string {
				return ""
			},
		},
		{
			header: "MQTT",
			w:      13,
			h:      8,
			getData: func(d models.DeviceLogs) string {
				return fmt.Sprint(d.Mqtt)
			},
		},
		{
			header: "Invid",
			w:      16,
			h:      8,
			getData: func(d models.DeviceLogs) string {
				return d.Invid
			},
		},
		{
			header: "Msg",
			w:      33,
			h:      8,
			getData: func(d models.DeviceLogs) string {
				return d.MsgId
			},
		},
		{
			header: "Text",
			w:      40,
			h:      8,
			getData: func(d models.DeviceLogs) string {
				return d.Text
			},
		},
		{
			header: "Context",
			w:      18,
			h:      8,
			getData: func(d models.DeviceLogs) string {
				return d.Context
			},
		},
		{
			header: "Class",
			w:      20,
			h:      8,
			getData: func(d models.DeviceLogs) string {
				return d.ClassMsg
			},
		},
		{
			header: "Level",
			w:      15,
			h:      8,
			getData: func(d models.DeviceLogs) string {
				return fmt.Sprint(d.Level)
			},
		},
		{
			header: "Area",
			w:      20,
			h:      8,
			getData: func(d models.DeviceLogs) string {
				return d.Area
			},
		},
		{
			header: "Address",
			w:      45,
			h:      8,
			getData: func(d models.DeviceLogs) string {
				return d.Addr
			},
		},
		{
			header: "Block",
			w:      15,
			h:      8,
			getData: func(d models.DeviceLogs) string {
				return d.Block
			},
		},
		{
			header: "Type",
			w:      10,
			h:      8,
			getData: func(d models.DeviceLogs) string {
				return d.Type
			},
		},
		{
			header: "Bit",
			w:      12,
			h:      8,
			getData: func(d models.DeviceLogs) string {
				return d.Bit
			},
		},
		{
			header: "Invert Bit",
			w:      16,
			h:      8,
			getData: func(d models.DeviceLogs) string {
				return d.InvertBit
			},
		},
	}
}
