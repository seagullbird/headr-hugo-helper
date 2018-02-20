package main

import "path/filepath"

var (
	newsiteQueueName = "new_site"
	initialThemeName = "gohugo-theme-ananke"
	dataDir = "/data"
	sitesDir = filepath.Join(dataDir, "sites")
	themesDir = filepath.Join(dataDir, "themes")
	configsDir = filepath.Join(dataDir, "configs")
)