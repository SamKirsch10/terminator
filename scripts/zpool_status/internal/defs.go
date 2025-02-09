package internal

type Disk struct {
	Name           string   `json:"name"`
	VdevType       VDevType `json:"vdev_type"`
	GUID           string   `json:"guid"`
	Class          string   `json:"class"`
	State          string   `json:"state"`
	AllocSpace     string   `json:"alloc_space"`
	TotalSpace     string   `json:"total_space"`
	DefSpace       string   `json:"def_space"`
	ReadErrors     string   `json:"read_errors"`
	WriteErrors    string   `json:"write_errors"`
	ChecksumErrors string   `json:"checksum_errors"`
}

type VDev struct {
	Disk
	DiskMakeup map[string]Disk `json:"vdevs"`
}

type ZPool struct {
	Name       string `json:"name"`
	State      string `json:"state"`
	PoolGUID   string `json:"pool_guid"`
	Txg        string `json:"txg"`
	SpaVersion string `json:"spa_version"`
	ZplVersion string `json:"zpl_version"`
	ErrorCount string `json:"error_count"`
	ScanStats  struct {
		Function           string `json:"function"`
		State              string `json:"state"`
		StartTime          string `json:"start_time"`
		EndTime            string `json:"end_time"`
		ToExamine          string `json:"to_examine"`
		Examined           string `json:"examined"`
		Skipped            string `json:"skipped"`
		Processed          string `json:"processed"`
		Errors             string `json:"errors"`
		BytesPerScan       string `json:"bytes_per_scan"`
		PassStart          string `json:"pass_start"`
		ScrubPause         string `json:"scrub_pause"`
		ScrubSpentPaused   string `json:"scrub_spent_paused"`
		IssuedBytesPerScan string `json:"issued_bytes_per_scan"`
		Issued             string `json:"issued"`
	} `json:"scan_stats"`
	Vdevs map[string]VDev `json:"vdevs"`
}

type ZpoolStatusOutput struct {
	OutputVersion struct {
		Command   string `json:"command"`
		VersMajor int    `json:"vers_major"`
		VersMinor int    `json:"vers_minor"`
	} `json:"output_version"`
	Pools map[string]ZPool `json:"pools"`
}

type PoolListing struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	State      string `json:"state"`
	PoolGUID   string `json:"pool_guid"`
	Txg        string `json:"txg"`
	SpaVersion string `json:"spa_version"`
	ZplVersion string `json:"zpl_version"`
	Properties struct {
		Size struct {
			Value  string `json:"value"`
			Source struct {
				Type string `json:"type"`
				Data string `json:"data"`
			} `json:"source"`
		} `json:"size"`
		Allocated struct {
			Value  string `json:"value"`
			Source struct {
				Type string `json:"type"`
				Data string `json:"data"`
			} `json:"source"`
		} `json:"allocated"`
		Free struct {
			Value  string `json:"value"`
			Source struct {
				Type string `json:"type"`
				Data string `json:"data"`
			} `json:"source"`
		} `json:"free"`
		Checkpoint struct {
			Value  string `json:"value"`
			Source struct {
				Type string `json:"type"`
				Data string `json:"data"`
			} `json:"source"`
		} `json:"checkpoint"`
		Expandsize struct {
			Value  string `json:"value"`
			Source struct {
				Type string `json:"type"`
				Data string `json:"data"`
			} `json:"source"`
		} `json:"expandsize"`
		Fragmentation struct {
			Value  string `json:"value"`
			Source struct {
				Type string `json:"type"`
				Data string `json:"data"`
			} `json:"source"`
		} `json:"fragmentation"`
		Capacity struct {
			Value  string `json:"value"`
			Source struct {
				Type string `json:"type"`
				Data string `json:"data"`
			} `json:"source"`
		} `json:"capacity"`
		Dedupratio struct {
			Value  string `json:"value"`
			Source struct {
				Type string `json:"type"`
				Data string `json:"data"`
			} `json:"source"`
		} `json:"dedupratio"`
		Health struct {
			Value  string `json:"value"`
			Source struct {
				Type string `json:"type"`
				Data string `json:"data"`
			} `json:"source"`
		} `json:"health"`
		Altroot struct {
			Value  string `json:"value"`
			Source struct {
				Type string `json:"type"`
				Data string `json:"data"`
			} `json:"source"`
		} `json:"altroot"`
	} `json:"properties"`
}

type ZpoolListOutput struct {
	OutputVersion struct {
		Command   string `json:"command"`
		VersMajor int    `json:"vers_major"`
		VersMinor int    `json:"vers_minor"`
	} `json:"output_version"`
	Pools map[string]PoolListing `json:"pools"`
}
