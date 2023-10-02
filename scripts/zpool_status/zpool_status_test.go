package main

import (
	"os"
	"testing"
)

func parseTestFile(testFile string, t *testing.T) *ZPool {
	cmdBytes, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal("failed to read in test file")
	}

	data, err := parseStatusOutput(string(cmdBytes))
	if err != nil {
		t.Fatal("func parseStatusOutput returned error", err)
	}

	return data
}

func TestStatusScrub(t *testing.T) {
	data := parseTestFile("tests/scrub.txt", t)

	if data.State != "ONLINE/SCRUBBING" {
		t.Errorf("Pool state check failed, got: %s, want %s", data.State, "ONLINE/SCRUBBING")
	}

	if data.Scan.PercentDone != 38.89 {
		t.Errorf("Pool scan percent done check failed, got: %f, want %f", data.Scan.PercentDone, 38.89)
	}
}

func TestStatus(t *testing.T) {
	data := parseTestFile("tests/normal.txt", t)
	if data.State != "ONLINE" {
		t.Errorf("Pool state check failed, got: %s, want %s", data.State, "ONLINE")
	}

	if len(data.Disks) != 5 {
		t.Errorf("Pool disk count failed, got: %d, want %d", len(data.Disks), 5)
	}

}
