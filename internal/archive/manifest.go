package archive

import (
	"encoding/json"
	"sort"

	"independentweeklylog/internal/domain"
)

type Manifest struct {
	EntryID     string   `json:"entry_id"`
	Title       string   `json:"title"`
	Week        int      `json:"week"`
	ResourceIDs []string `json:"resource_ids"`
	Tags        []string `json:"tags"`
}

func BuildManifest(entry domain.JournalEntry, resources []domain.Resource) Manifest {
	ids := make([]string, 0, len(resources))
	for _, resource := range resources {
		ids = append(ids, resource.ID)
	}
	sort.Strings(ids)
	return Manifest{EntryID: entry.ID, Title: entry.Title, Week: entry.Week, ResourceIDs: ids, Tags: domain.NormalizeList(entry.Tags)}
}

func EncodeManifest(manifest Manifest) ([]byte, error) { return json.Marshal(manifest) }

func DecodeManifest(data []byte) (Manifest, error) {
	var manifest Manifest
	err := json.Unmarshal(data, &manifest)
	return manifest, err
}

func ManifestMatches(manifest Manifest, entry domain.JournalEntry) bool {
	return manifest.EntryID == entry.ID && manifest.Week == entry.Week && manifest.Title == entry.Title
}
