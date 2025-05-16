package v8

import (
	"encoding/json"
	"testing"
)

func TestTermsSetQueryWithField(t *testing.T) {
	q := NewTermsSetQuery("codes", "abc", "def", "ghi").MinimumShouldMatchField("required_matches")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"terms_set":{"codes":{"minimum_should_match_field":"required_matches","terms":["abc","def","ghi"]}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestTermsSetQueryWithScript(t *testing.T) {
	q := NewTermsSetQuery("codes", "abc", "def", "ghi").
		MinimumShouldMatchScript(
			NewScript(`Math.min(params.num_terms, doc['required_matches'].value)`),
		)
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"terms_set":{"codes":{"minimum_should_match_script":{"source":"Math.min(params.num_terms, doc['required_matches'].value)"},"terms":["abc","def","ghi"]}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
