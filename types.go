package gvm

type Version struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
	Files   []File `json:"files"`
}

type File struct {
	Filename string `json:"filename"`
	Os       string `json:"os"`
	Arch     string `json:"arch"`
	Size     int64  `json:"size"`
	Sha256   string `json:"sha256"`
	Version  string `json:"version"`
	Kind     string `json:"kind"`
}
