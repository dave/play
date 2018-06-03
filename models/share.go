package models

// SharePack is the structure of the data persisted on src.jsgo.io as json, so best to use json tags
// to lower-case the names.
type SharePack struct {
	Version int                          `json:"version"`
	Source  map[string]map[string]string `json:"source"` // Source packages for this build: map[<package>]map[<filename>]<contents>
	Tags    []string                     `json:"tags"`   // Build tags
}
