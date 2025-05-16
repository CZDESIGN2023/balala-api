package data

import (
	"context"
	"go-cs/pkg/pprint"
	"testing"
)

func Test_SearchHistoryTaskWorkFlowTemplateByOwnerRule(t *testing.T) {
	rule, err := WorkFlowRepo.SearchHistoryTaskWorkFlowTemplateByOwnerRule(context.Background(), 87, "21")
	if err != nil {
		t.Error(err)
	}

	pprint.Println(rule)
}

func Test_ClearHistoryTemplate(t *testing.T) {
	err := WorkFlowRepo.ClearHistoryTemplate(context.Background(), 13866)
	if err != nil {
		t.Error(err)
	}

}
