package nextdate

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNextDate(t *testing.T) {
	now := time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		dstart   string
		repeat   string
		expected string
		wantErr  bool
	}{
		{
			name:     "повторение через семь дней",
			dstart:   "20240113",
			repeat:   "d 7",
			expected: "20240127",
		},
		{
			name:     "ежедневное повторение",
			dstart:   "20240125",
			repeat:   "d 1",
			expected: "20240127",
		},
		{
			name:     "ежегодное повторение",
			dstart:   "20240126",
			repeat:   "y",
			expected: "20250126",
		},
		{
			name:     "високосный год",
			dstart:   "20240229",
			repeat:   "y",
			expected: "20250301",
		},
		{
			name:    "пустое правило",
			dstart:  "20240113",
			repeat:  "",
			wantErr: true,
		},
		{
			name:    "некорректная дата",
			dstart:  "20240230",
			repeat:  "d 1",
			wantErr: true,
		},
		{
			name:    "слишком большой интервал",
			dstart:  "20240113",
			repeat:  "d 405",
			wantErr: true,
		},
		{
			name:    "не указан интервал",
			dstart:  "20240113",
			repeat:  "d",
			wantErr: true,
		},
		{
			name:    "лишний параметр для ежегодного повторения",
			dstart:  "20240113",
			repeat:  "y 2",
			wantErr: true,
		},
		{
			name:    "неизвестное правило",
			dstart:  "20240113",
			repeat:  "l 34",
			wantErr: true,
		},
		{
			name:     "повторение по воскресеньям",
			dstart:   "20240125",
			repeat:   "w 7",
			expected: "20240128",
		},
		{
			name:     "повторение по понедельникам, четвергам и пятницам",
			dstart:   "20240126",
			repeat:   "w 1,4,5",
			expected: "20240129",
		},
		{
			name:    "недопустимый день недели",
			dstart:  "20240125",
			repeat:  "w 8",
			wantErr: true,
		},
		{
			name:    "не указаны дни недели",
			dstart:  "20240125",
			repeat:  "w",
			wantErr: true,
		},

		{
			name:     "повторение по 16 и 5 числам месяца",
			dstart:   "20240116",
			repeat:   "m 16,5",
			expected: "20240205",
		},
		{
			name:     "повторение по последнему и 18 числу",
			dstart:   "20240201",
			repeat:   "m -1,18",
			expected: "20240218",
		},
		{
			name:     "повторение по третьему числу января марта и июня",
			dstart:   "20240103",
			repeat:   "m 3 1,3,6",
			expected: "20240303",
		},
		{
			name:     "повторение по предпоследнему числу февраля",
			dstart:   "20240130",
			repeat:   "m -2",
			expected: "20240228",
		},
		{
			name:    "недопустимый день месяца",
			dstart:  "20240113",
			repeat:  "m 4,40",
			wantErr: true,
		},
		{
			name:    "недопустимый месяц",
			dstart:  "20240113",
			repeat:  "m 5 1,13",
			wantErr: true,
		},
		{
			name:    "несуществующее число месяца",
			dstart:  "20240113",
			repeat:  "m 31 2",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NextDate(now, tt.dstart, tt.repeat)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
