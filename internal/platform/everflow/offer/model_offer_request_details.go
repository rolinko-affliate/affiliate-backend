package offer

type Details []struct {
	Total   int         `json:"total"`
	Entries EntriesInfo `json:"entries"`
}
