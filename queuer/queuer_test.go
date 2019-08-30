package queuer

import (
	"testing"

	"os/exec"
)

var initSQSMocker = func() {
	exec.Command("docker run -p 9324:9324 softwaremill/elasticmq")
}

func TestQueuer_Init(t *testing.T) {
	tests := []struct {
		name    string
		q       *Queuer
		wantNil bool
	}{
		{
			name:    "Success",
			q:       New(),
			wantNil: false,
		},
	}
	for _, tt := range tests {
		initSQSMocker()
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.Init(); (got != nil) == tt.wantNil {
				t.Errorf("Queuer.Init() = %v, want %v", got, tt.wantNil)
			}
		})
	}
}

func TestQueuer_CreateQueue(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		q       *Queuer
		args    args
		wantErr bool
	}{
		{
			name:    "Success Create Queue",
			q:       New(),
			args:    args{"QUEUE_NAME"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		initSQSMocker()
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.q.CreateQueue(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("Queuer.CreateQueue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQueuer_GetQueueURL(t *testing.T) {
	type args struct {
		queueName string
	}
	tests := []struct {
		name    string
		q       *Queuer
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "Success Get Queue URL",
			q:       New(),
			args:    args{"QUEUE_NAME"},
			want:    "http://localhost:9324/queue/QUEUE_NAME",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		initSQSMocker()
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.q.GetQueueURL(tt.args.queueName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Queuer.GetQueueURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Queuer.GetQueueURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueuer_WriteToQueue(t *testing.T) {
	type args struct {
		queueURL string
		body     []byte
	}
	tests := []struct {
		name    string
		q       *Queuer
		args    args
		wantErr bool
	}{
		{
			name:    "Success Write To Queue",
			q:       New(),
			args:    args{"http://localhost:9324/queue/QUEUE_NAME", []byte("Hello")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		initSQSMocker()
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.q.WriteToQueue(tt.args.queueURL, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("Queuer.WriteToQueue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQueuer_ReadFromQueue(t *testing.T) {
	type args struct {
		queueURL string
	}
	tests := []struct {
		name    string
		q       *Queuer
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "Success Read from QUEUE",
			q:       New(),
			args:    args{"http://localhost:9324/queue/QUEUE_NAME"},
			want:    "Hello",
			wantErr: false,
		},
		{
			name:    "Non Existent Queue",
			q:       New(),
			args:    args{"http://localhost:9324/queue/UNEXISTING"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.q.WriteToQueue(tt.args.queueURL, []byte("Hello"))
			got, err := tt.q.ReadFromQueue(tt.args.queueURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("Queuer.ReadFromQueue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Queuer.ReadFromQueue() = %v, want %v", got, tt.want)
			}
		})
	}
}
