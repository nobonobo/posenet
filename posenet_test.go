package posenet

import (
	"reflect"
	"syscall/js"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		config Config
	}
	tests := []struct {
		name string
		args args
		want *PoseNet
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPoseNet_Start(t *testing.T) {
	type fields struct {
		net    js.Value
		video  js.Value
		Config Config
	}
	type args struct {
		videoID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &PoseNet{
				net:    tt.fields.net,
				video:  tt.fields.video,
				Config: tt.fields.Config,
			}
			if err := n.Start(tt.args.videoID); (err != nil) != tt.wantErr {
				t.Errorf("PoseNet.Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPoseNet_Stop(t *testing.T) {
	type fields struct {
		net    js.Value
		video  js.Value
		Config Config
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &PoseNet{
				net:    tt.fields.net,
				video:  tt.fields.video,
				Config: tt.fields.Config,
			}
			n.Stop()
		})
	}
}
