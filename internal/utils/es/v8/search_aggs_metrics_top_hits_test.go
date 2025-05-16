package v8

import (
	"encoding/json"
	"testing"
)

func TestTopHitsAggregation(t *testing.T) {
	fsc := NewFetchSourceContext(true).Include("title")
	agg := NewTopHitsAggregation().
		Sort("last_activity_date", false).
		FetchSourceContext(fsc).
		Size(1)
	src, err := agg.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"top_hits":{"_source":{"includes":["title"]},"size":1,"sort":[{"last_activity_date":{"order":"desc"}}]}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
