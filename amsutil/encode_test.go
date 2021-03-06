package amsutil

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestEncode(t *testing.T) {
	ctx := context.TODO()
	cnf := testConfigFromFile(t, "config.json")
	AMS, err := cnf.Client(ctx)
	if err != nil {
		t.Fatalf("client construct failed: %v", err)
	}

	f, err := os.Open(filepath.Join(cnf.BaseDir, "testdata", "small.mp4"))
	if err != nil {
		t.Fatalf("video file open failed: %v", err)
	}
	defer f.Close()

	asset, err := UploadFile(ctx, AMS, f, 4*1024*1024, 5)
	if err != nil {
		t.Fatalf("file uploading failed: %v", err)
	}

	mediaProcessors, err := AMS.GetMediaProcessors(ctx)
	if err != nil {
		t.Fatalf("get media processors failed: %v", err)
	}

	var MES string
	for _, mediaProcessor := range mediaProcessors {
		if mediaProcessor.Name == "Media Encoder Standard" {
			MES = mediaProcessor.ID
			break
		}
	}

	if len(MES) == 0 {
		t.Fatal("'Media Encoder Standard' not found")
	}

	encodedAssets, job, err := Encode(ctx, AMS, asset.ID, MES, "Adaptive Streaming")
	if err != nil {
		t.Fatalf("encode rejected: %v", err)
	}

	if err := WaitJob(ctx, AMS, job.ID, 3*time.Second); err != nil {
		t.Fatalf("encode failed: %v", err)
	}

	if err := AMS.DeleteAsset(ctx, asset.ID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	for _, encodedAsset := range encodedAssets {
		if err := AMS.DeleteAsset(ctx, encodedAsset.ID); err != nil {
			t.Fatalf("delete failed [asset#%v]: %v", encodedAsset.ID, err)
		}
	}
}
