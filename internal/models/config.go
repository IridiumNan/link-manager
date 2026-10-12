package models

type LinkConfig struct {
	// DirMap store the (src_dir -> dst_dir) key-value pair
	DirMap map[string]string `mapstructure:"dir_map" json:"dir_map"`

	// FileMap store the (src_name -> dst_name) key-value pair
	FileMap map[string]string `mapstructure:"file_map" json:"file_map"`
}
