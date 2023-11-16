package config

import "testing"

func TestSchedulerConfig_GetNumDigesters(t *testing.T) {
	type fields struct {
		NumDigesters       int
		NumDigestersCustom string
	}
	tests := []struct {
		name     string
		fields   fields
		clientID int
		want     int
	}{
		{"case empty", fields{123, ``}, 1, 123},
		{"case empty array", fields{123, `[]`}, 1, 123},
		{"case wrong json", fields{123, `[}asd`}, 1, 123},
		{"case json not match", fields{123, `{"client_id":1,"value":11}`}, 1, 123},

		{"case normal client 1", fields{123, `[{"client_id":1,"value":11},{"client_id":2,"value":22}]`}, 1, 11},
		{"case normal client 2", fields{123, `[{"client_id":1,"value":11},{"client_id":2,"value":22}]`}, 2, 22},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := SchedulerConfig{
				NumDigesters:       tt.fields.NumDigesters,
				NumDigestersCustom: tt.fields.NumDigestersCustom,
			}
			if got := sc.GetNumDigesters(tt.clientID); got != tt.want {
				t.Errorf("GetNumDigesters() = %v, want %v", got, tt.want)
			}
		})
	}
}
