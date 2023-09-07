package fpsclient

import (
	"reflect"
	"testing"
	"time"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fpsclient"
)

func Test_toClientDTO(t *testing.T) {
	timeNow := time.Now()
	client := &fpsclient.Client{
		FpsClient: ent.FpsClient{
			ID:            1,
			ClientID:      -1,
			Name:          "abc",
			Description:   "desc",
			SampleFileURL: "https://abc.com/sample_file.xlsx",
			CreatedAt:     timeNow,
			CreatedBy:     "created by asdfs",
			UpdatedAt:     timeNow,
		},
	}

	clientExpected := ClientDTO{
		ID:            client.ID,
		Name:          client.Name,
		Description:   client.Description,
		SampleFileURL: client.SampleFileURL,
		CreatedAt:     client.CreatedAt.UnixMilli(),
		CreatedBy:     client.CreatedBy,
	}

	type args struct {
		client *fpsclient.Client
	}
	tests := []struct {
		name string
		args args
		want ClientDTO
	}{
		{"test toClientDTO with full data", args{client: client}, clientExpected},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toClientDTO(tt.args.client); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toClientDTO() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
