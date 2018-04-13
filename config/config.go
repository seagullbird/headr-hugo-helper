package config

import "path/filepath"

var (
	// NewsiteQueueName is the name for the NewSite queue
	NewsiteQueueName = "new_site"
	// RegenerateQueueName is the name for the ReGenerate queue
	RegenerateQueueName = "re_generate"
	// DataDir is the root directory of the persistent volume
	DataDir = "/data"
	// SitesDir is the path for sites
	SitesDir = filepath.Join(DataDir, "sites")
	// ThemesDir is the path for themes
	ThemesDir = filepath.Join(DataDir, "themes")
	// ConfigsDir is the path for initial configs of each theme
	ConfigsDir = filepath.Join(DataDir, "configs")
	// Dev indicates whether the environment is dev or prod; it's set during compiling
	Dev = "unset"
)
