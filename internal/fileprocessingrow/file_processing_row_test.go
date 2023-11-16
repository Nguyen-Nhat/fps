package fileprocessingrow

import "testing"

func TestCustomStatisticModel_IsSuccessAll(t *testing.T) {
	type fields struct {
		Statuses string
		Count    int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// no success
		{"test IsSuccessAll case no statuses", fields{"", 0}, false},
		{"test IsSuccessAll case 1 status but it is 1", fields{"1", 1}, false},
		{"test IsSuccessAll case 1 status but it is 2", fields{"2", 1}, false},
		{"test IsSuccessAll case 1 status but it is 3", fields{"3", 1}, false},
		{"test IsSuccessAll case 1 status but it is 5", fields{"5", 1}, false},
		{"test IsSuccessAll case 2 statuses, has one is 1", fields{"1,4", 2}, false},
		{"test IsSuccessAll case 2 statuses, has one is 2", fields{"2,4", 2}, false},
		{"test IsSuccessAll case 2 statuses, has one is 3", fields{"3,4", 2}, false},
		{"test IsSuccessAll case 2 statuses, has one is 5", fields{"5,4", 2}, false},
		{"test IsSuccessAll case 3 statuses, has one is not success", fields{"4,1,4", 3}, false},
		{"test IsSuccessAll case 4 statuses, has one is not success", fields{"4,1,4,4", 4}, false},
		{"test IsSuccessAll case 5 statuses, has one is not success", fields{"4,1,4,4,4", 5}, false},
		// success
		{"test IsSuccessAll case 1 status, all is 4", fields{"4", 1}, true},
		{"test IsSuccessAll case 2 statuses, all is 4", fields{"4,4", 2}, true},
		{"test IsSuccessAll case 3 statuses, all is 4", fields{"4,4,4", 3}, true},
		{"test IsSuccessAll case 4 statuses, all is 4", fields{"4,4,4,4", 4}, true},
		{"test IsSuccessAll case 5 statuses, all is 4", fields{"4,4,4,4,4", 5}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &CustomStatisticModel{
				Statuses: tt.fields.Statuses,
				Count:    tt.fields.Count,
			}
			if got := s.IsSuccessAll(); got != tt.want {
				t.Errorf("IsSuccessAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCustomStatisticModel_IsContainsFailed(t *testing.T) {
	tests := []struct {
		name     string
		Statuses string
		want     bool
	}{
		{"test IsContainsFailed case 1 status, it is 1", "1", false},
		{"test IsContainsFailed case 1 status, it is 2", "2", false},
		{"test IsContainsFailed case 1 status, it is 3", "3", true},
		{"test IsContainsFailed case 1 status, it is 4", "4", false},
		{"test IsContainsFailed case 1 status, it is 5", "5", false},
		{"test IsContainsFailed case 2 statuses, has one is 3", "1,3", true},
		{"test IsContainsFailed case 2 statuses, no one is 3", "1,2", false},
		{"test IsContainsFailed case 2 statuses, no one is 3", "2,4", false},
		{"test IsContainsFailed case 3 statuses, has one is 3", "4,4,3", true},
		{"test IsContainsFailed case 3 statuses, has one is 3", "4,3,5", true},
		{"test IsContainsFailed case 3 statuses, no one is 3", "4,4,4", false},
		{"test IsContainsFailed case 4 statuses, has one is 3", "4,3,4,4", true},
		{"test IsContainsFailed case 4 statuses, no one is 3", "4,4,4,4", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &CustomStatisticModel{Statuses: tt.Statuses}
			if got := s.IsContainsFailed(); got != tt.want {
				t.Errorf("IsContainsFailed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCustomStatisticModel_IsProcessed(t *testing.T) {
	tests := []struct {
		name     string
		statuses string
		want     bool
	}{
		{"test IsProcessed case has both 3 and 4", "4,3,3,4", true},
		{"test IsProcessed case not has at least 3 or 4", "1,2,1,5", false},
		{"test IsProcessed case has only 3", "1,2,3,5", true},
		{"test IsProcessed case has only 4", "1,2,1,4", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &CustomStatisticModel{Statuses: tt.statuses}
			if got := s.IsProcessed(); got != tt.want {
				t.Errorf("IsProcessed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCustomStatisticModel_IsWaiting(t *testing.T) {
	tests := []struct {
		name     string
		statuses string
		want     bool
	}{
		{"test IsProcessed case has 1", "4,3,3,4", false},
		{"test IsProcessed case has 2", "4,2", false},
		{"test IsProcessed case has 3", "4,3,3,4", false},
		{"test IsProcessed case has 4", "4,3,3,4", false},
		{"test IsProcessed case has 5", "4,5,3,4", true},
		{"test IsProcessed case has 6", "4,3,3,6", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &CustomStatisticModel{Statuses: tt.statuses}
			if got := s.IsWaiting(); got != tt.want {
				t.Errorf("IsWaiting() = %v, want %v", got, tt.want)
			}
		})
	}
}
