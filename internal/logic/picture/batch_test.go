package picture

import (
	v1 "cloud/api/user/v1"
	"context"
	"reflect"
	"testing"
)

func Test_sPicture_UploadByBatch(t *testing.T) {
	type args struct {
		ctx context.Context
		req *v1.PictureUploadByBatchReq
	}
	tests := []struct {
		name    string
		args    args
		wantRes *v1.PictureUploadByBatchRes
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				ctx: context.Background(),
				req: &v1.PictureUploadByBatchReq{
					NamePrefix: "",
					SearchText: "苹果",
					Count:      3,
				},
			},
			wantRes: nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sPicture{}
			gotRes, err := s.UploadByBatch(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadByBatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("UploadByBatch() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
