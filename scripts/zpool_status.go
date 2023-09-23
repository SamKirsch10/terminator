///usr/bin/env true; exec /usr/bin/env go run "$0" "$@"

package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

const cmdOutput = `  pool: tank
state: DEGRADED
status: One or more devices has been taken offline by the administrator.
   Sufficient replicas exist for the pool to continue functioning in a
   degraded state.
action: Online the device using 'zpool online' or replace the device with
   'zpool replace'.
config:

   NAME                                                  STATE     READ WRITE CKSUM
   tank                                                  DEGRADED     0     0     0
	 raidz2-0                                            DEGRADED     0     0     0
	   ata-TOSHIBA_HDWG440_23Q0A056FZ1G-part1            ONLINE       0     0     0
	   ata-WDC_WD4000F9YZ-09N20L1_WD-WMC5D0D25KFW-part1  ONLINE       0     0     0
	   ata-WDC_WD4003FRYZ-01F0DB0_V1J1PJZG-part1         ONLINE       0     0     0
	   ata-WDC_WD4003FRYZ-01F0DB0_V1KW5XVG-part1         ONLINE       0     0     0
	   /tmp/fake.img                                     OFFLINE      0     0     0

errors: No known data errors`

type ZDisk struct {
	Disk        string
	State       string
	ReadErrors  string
	WriteErrors string
	CKSumErrors string
}

type ZPool struct {
	Name  string
	Type  string
	State string
	Disks []ZDisk
}

func delete_empty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, strings.Replace(str, " ", "", -1))
		}
	}
	return r
}

func main() {
	fmt.Println("Getting Info")

	out := new(ZPool)

	out.Name = regexp.MustCompile("pool:\\s(.*)\n").FindAllStringSubmatch(cmdOutput, 1)[0][1]
	out.State = regexp.MustCompile("state:\\s(.*)\n").FindAllStringSubmatch(cmdOutput, 1)[0][1]

	tmp := cmdOutput[strings.Index(cmdOutput, "config:")+7 : strings.Index(cmdOutput, "errors:")]

	var dataFound bool
	var poolFound bool
	for _, line := range strings.Split(tmp, "\n") {
		d := new(ZDisk)
		if line == "" {
			continue
		}
		elements := delete_empty(strings.Split(strings.Replace(line, "\t", "", -1), "  "))
		if elements[0] == out.Name {
			poolFound = true
		} else if elements[0][0:4] == "raid" {
			dataFound = true
			continue
		}

		if dataFound && poolFound {
			d.Disk = elements[0]
			d.State = elements[1]
			d.ReadErrors = elements[2]
			d.WriteErrors = elements[3]
			d.CKSumErrors = elements[4]
			out.Disks = append(out.Disks, *d)
		}
	}

	j_son, _ := json.Marshal(out)

	fmt.Println(string(j_son))

}
