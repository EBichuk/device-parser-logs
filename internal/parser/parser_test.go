package parser

import (
	"device-parser-logs/internal/models"
	"reflect"
	"testing"
)

func TestParsertoLogDevice(t *testing.T) {
	tests := []struct {
		name    string
		log     []string
		want    *models.DeviceLogs
		wantErr bool
	}{
		{
			name: "success",
			log: []string{"1", "", "G-044322", "01749246-95f6-57db-b7c3-2ae0e8be671f", "cold7_Temp_Al_HH",
				"Высокая температура", "", "alarm", "100", "LOCAL", "cold7_status.Temp_Al_HH", "", "", "", ""},
			want: &models.DeviceLogs{
				Mqtt:      "",
				Invid:     "G-044322",
				Guid:      "01749246-95f6-57db-b7c3-2ae0e8be671f",
				MsgId:     "cold7_Temp_Al_HH",
				Text:      "Высокая температура",
				Context:   "",
				ClassMsg:  "alarm",
				Level:     100,
				Area:      "LOCAL",
				Addr:      "cold7_status.Temp_Al_HH",
				Block:     "",
				Type:      "",
				Bit:       "",
				InvertBit: "",
			},
			wantErr: false,
		},
		{
			name: "error",
			log: []string{"2", "", "G-044322", "01749246-95f6-57db-b7c3-2ae0e8be671f", "cold7_Temp_Al_HH",
				"Высокая температура", "", "alarm", "", "LOCAL", "cold7_status.Temp_Al_HH", "", "", "", ""},
			want: &models.DeviceLogs{
				Mqtt:      "",
				Invid:     "G-044322",
				Guid:      "01749246-95f6-57db-b7c3-2ae0e8be671f",
				MsgId:     "cold7_Temp_Al_HH",
				Text:      "Высокая температура",
				Context:   "",
				ClassMsg:  "alarm",
				Level:     -1,
				Area:      "LOCAL",
				Addr:      "cold7_status.Temp_Al_HH",
				Block:     "",
				Type:      "",
				Bit:       "",
				InvertBit: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToDeviceLogs(tt.log)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParsertoLogDevice() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParsertoLogDevice() = %v, want %v", got, tt.want)
			}
		})
	}
}
