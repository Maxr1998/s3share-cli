package store

type FileMetadata struct {
	Name   string `json:"name"`
	NameIv string `json:"name_iv"`
	Iv     string `json:"iv"`
	Size   int64  `json:"size"`
}
