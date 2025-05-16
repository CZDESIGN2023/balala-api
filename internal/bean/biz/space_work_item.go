package db

import "encoding/json"

type DbSpaceWorkItem struct {
	SpaceWorkItemV2
	SpaceWorkItemDocV2
}

func (m *DbSpaceWorkItem) TableName() string {
	return m.SpaceWorkItemV2.TableName()
}

func (m *DbSpaceWorkItem) DirectorSlice() []string {
	var ids []string
	if m.SpaceWorkItemDocV2.Directors == "" || m.SpaceWorkItemDocV2.Directors == "[]" {
		return ids
	}

	json.Unmarshal([]byte(m.SpaceWorkItemDocV2.Directors), &ids)
	return ids
}

func (m *DbSpaceWorkItem) FollowerSlice() []string {
	var ids []string
	if m.SpaceWorkItemDocV2.Followers == "" || m.SpaceWorkItemDocV2.Followers == "[]" {
		return ids
	}

	json.Unmarshal([]byte(m.SpaceWorkItemDocV2.Followers), &ids)
	return ids
}

func (m *DbSpaceWorkItem) ParticipatorSlice() []string {
	var ids []string
	if m.SpaceWorkItemDocV2.Participators == "" || m.SpaceWorkItemDocV2.Participators == "[]" {
		return ids
	}

	json.Unmarshal([]byte(m.SpaceWorkItemDocV2.Participators), &ids)
	return ids
}

func (m *DbSpaceWorkItem) TagSlice() []string {
	var ids []string
	if m.SpaceWorkItemDocV2.Tags == "" || m.SpaceWorkItemDocV2.Tags == "[]" {
		return ids
	}

	json.Unmarshal([]byte(m.SpaceWorkItemDocV2.Tags), &ids)
	return ids
}
