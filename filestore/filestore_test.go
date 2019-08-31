package filestore

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestFileManager_CreateBucket(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		f       *FileManager
		args    args
		wantErr bool
	}{
		{
			name:    "Create Bucket",
			f:       New(),
			args:    args{"preview_bucket"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.f.CreateBucket(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("FileManager.CreateBucket() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileManager_Upload(t *testing.T) {
	readFile := func(filename string) []byte {
		file, err := os.Open(filename)
		if err != nil {
			t.Error(err)
		}

		b, err := ioutil.ReadAll(file)
		if err != nil {
			t.Error(err)
		}
		return b
	}
	type args struct {
		bucket   string
		filename string
		reader   io.Reader
	}
	tests := []struct {
		name    string
		f       *FileManager
		args    args
		wantErr bool
	}{
		{
			name: "Upload to S3",
			f:    New(),
			args: args{"preview_bucket", "testfile", bytes.NewReader(readFile("../samples/invoice_preview.pdf"))},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.f.Upload(tt.args.bucket, tt.args.filename, tt.args.reader); (err != nil) != tt.wantErr {
				t.Errorf("FileManager.Upload() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileManager_Download(t *testing.T) {
	type args struct {
		bucket   string
		filename string
	}
	tests := []struct {
		name string
		f    *FileManager
		args args
		// want    io.Reader
		wantErr bool
	}{
		{
			name: "Test Download",
			f:    New(),
			args: args{"preview_bucket", "testfile"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.Download(tt.args.bucket, tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileManager Download() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			buf := make([]byte, 100000)
			n, err := got.Read(buf)
			f, err := os.Create("test.pdf")
			f.Write(buf[:n])
		})
	}
}
