package provision

import (
	"testing"
)

func TestProvisionSha256(t *testing.T) {
	_, err := downloadSha256("v0.0.1")
	if err != nil {
		t.Errorf("downloadFramework fail: %s\n", err)
	}
}

func TestProvisionDownload(t *testing.T) {
	path, err := downloadFramework("v0.0.1")
	if err != nil {
		t.Errorf("downloadFramework fail: %s\n", err)
	}
	t.Logf("downloaded: %s\n", path)
}

func TestProvisionFetchFramework(t *testing.T) {
	err := FetchFramework("v0.0.1")
	if err != nil {
		t.Errorf("GetFramework fail: %#v\n", err)
	}
}
