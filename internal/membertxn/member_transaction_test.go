package membertxn

import (
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
	"testing"
	"time"
)

func TestMemberTransaction_isCheckExpires(t1 *testing.T) {
	type fields struct {
		MemberTransaction ent.MemberTransaction
	}
	type args struct {
		expiresTimeMinutes int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "Test expire time",
			fields: fields{ent.MemberTransaction{CreatedAt: time.Now()}},
			args:   args{1},
			want:   false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &MemberTransaction{
				MemberTransaction: tt.fields.MemberTransaction,
			}
			if got := t.IsCheckExpires(tt.args.expiresTimeMinutes); got != tt.want {
				t1.Errorf("isCheckExpires() = %v, want %v", got, tt.want)
			}
		})
	}
}
