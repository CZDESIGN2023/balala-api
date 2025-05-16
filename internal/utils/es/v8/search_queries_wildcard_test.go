package v8

import (
	"encoding/json"
	"testing"
)

func TestWildcardQuery(t *testing.T) {
	q := NewWildcardQuery("user", "ki*y??")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"wildcard":{"user":{"value":"ki*y??"}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestWildcardQueryWithBoost(t *testing.T) {
	q := NewWildcardQuery("user", "ki*y??").Boost(1.2)
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"wildcard":{"user":{"boost":1.2,"value":"ki*y??"}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestWildcardQueryWithCaseInsensitive(t *testing.T) {
	q := NewWildcardQuery("user", "ki*y??").CaseInsensitive(true)
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"wildcard":{"user":{"case_insensitive":true,"value":"ki*y??"}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
