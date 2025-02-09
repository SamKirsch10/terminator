package internal

import (
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

func parseDiskSize(sizeStr string) (float64, error) {
	if sizeStr == "" {
		return 0, errors.New("no string size detected")
	}
	sizeStr = strings.TrimSpace(sizeStr)
	log.Debug("trying to parse disk size: ", sizeStr)

	var unit float64 = 1
	switch strings.ToUpper(string(sizeStr[len(sizeStr)-1])) {
	case "K":
		unit = 1024
	case "M":
		unit = 1024 * 1024
	case "G":
		unit = 1024 * 1024 * 1024
	case "T":
		unit = 1024 * 1024 * 1024 * 1024
	}

	sizeStr = sizeStr[:len(sizeStr)-1]

	size, err := strconv.ParseFloat(sizeStr, 64)
	if err != nil {
		return 0, err
	}

	return size * unit, nil
}

// locateDiskByUUID - pass in something like `/dev/disk/by-partuuid/abdc-123`
// and get back something like `/dev/sdg`
func locateDiskByUUID(diskPath string) (string, error) {
	log.Debugf("checking disk path: %s", diskPath)
	realDisk, err := os.Readlink(diskPath)
	if err != nil {
		return "", err
	}
	log.Debugf("got back disk: %s", realDisk)
	// remove partition if it's linked, we just want to whole disk
	if regexp.MustCompile(`.*[0-9]`).MatchString(realDisk) {
		realDisk = string(realDisk[len(realDisk)-1])
	}
	parts := strings.Split(realDisk, "/")
	return "/dev/" + parts[len(parts)-1], nil
}
