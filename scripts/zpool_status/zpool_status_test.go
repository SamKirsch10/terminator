package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"testing"
)

func parseTestFile(testFile string, t *testing.T) ZpoolOutput {
	cmdBytes, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal("failed to read in test file")
	}

	var data ZpoolOutput
	err = json.Unmarshal(cmdBytes, &data)
	if err != nil {
		t.Fatal("func parseStatusOutput returned error", err)
	}

	return data
}

func TestStatusScrub(t *testing.T) {
	data := parseTestFile("tests/scrub.txt", t).Pools["tank"]
	state := getPoolState(data)

	if state != ZpoolState("ONLINE/SCRUBBING") {
		t.Errorf("Pool state check failed, got: %s, want %s", data.State, "ONLINE/SCRUBBING")
	}

	isEqual := func(a, b float64) bool {
		return math.Abs(a-b) <= 1e-9
	}

	want := 51.597052
	perc, err := parseScrubPercent(data)
	if err != nil {
		t.Errorf("%v", err)
	}
	if isEqual(perc, want) {
		fmt.Printf("%T : %T\n", perc, want)
		t.Errorf("Pool scan percent done check failed, got: %f, want %f", perc, want)
	}
}

func TestStatus(t *testing.T) {
	data := parseTestFile("tests/normal.txt", t).Pools["tank"]
	if data.State != "ONLINE" {
		t.Errorf("Pool state check failed, got: %s, want %s", data.State, "ONLINE")
	}

	var diskCount int
	for name, vdev := range data.Vdevs {
		if name == "tank" {
			continue
		}
		diskCount = diskCount + len(vdev.DiskMakeup)
	}

	if diskCount != 5 {
		t.Errorf("Pool disk count failed, got: %d, want %d", diskCount, 7)
	}

}
