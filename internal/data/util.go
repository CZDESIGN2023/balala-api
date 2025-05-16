package data

import (
	"strings"

	"github.com/spf13/cast"
)

type UserKey struct {
	builder strings.Builder
}

func (u *UserKey) User(userId int64) *UserKey {
	u.builder.WriteString("balala:user:")
	u.builder.WriteString(cast.ToString(userId))
	return u
}

func (u *UserKey) Key() string {
	return u.builder.String()
}

type SpaceKey struct {
	builder strings.Builder
}

func (s *SpaceKey) Space(spaceId int64) *SpaceKey {
	s.builder.WriteString("balala:space:")
	s.builder.WriteString(cast.ToString(spaceId))
	return s
}

func (s *SpaceKey) Member(userId int64) *SpaceKey {
	s.builder.WriteString(":mem:")
	s.builder.WriteString(cast.ToString(userId))
	return s
}

func (s *SpaceKey) Key() string {
	return s.builder.String()
}

func (s *SpaceKey) Wildcard() string {
	s.builder.WriteString(":*")
	return s.builder.String()
}

type WorkItemKey struct {
	builder strings.Builder
}

func (s *WorkItemKey) WorkItem(workItemId int64) *WorkItemKey {
	s.builder.WriteString("balala:workItem:")
	s.builder.WriteString(cast.ToString(workItemId))
	return s
}

func (s *WorkItemKey) Key() string {
	return s.builder.String()
}

type WorkObjectKey struct {
	builder strings.Builder
}

func (s *WorkObjectKey) WorkObject(workObjectId int64) *WorkObjectKey {
	s.builder.WriteString("balala:workObject:")
	s.builder.WriteString(cast.ToString(workObjectId))
	return s
}

func (s *WorkObjectKey) Key() string {
	return s.builder.String()
}

type TagKey struct {
	builder strings.Builder
}

func (s *TagKey) Tag(tagId int64) *TagKey {
	s.builder.WriteString("balala:tag:")
	s.builder.WriteString(cast.ToString(tagId))
	return s
}

func (s *TagKey) Key() string {
	return s.builder.String()
}

func NewTagKey(tagId int64) *TagKey {
	v := &TagKey{}
	return v.Tag(tagId)
}

type WorkVersionKey struct {
	builder strings.Builder
}

func (s *WorkVersionKey) WorkVersion(workVersionId int64) *WorkVersionKey {
	s.builder.WriteString("balala:workVersionId:")
	s.builder.WriteString(cast.ToString(workVersionId))
	return s
}

func (s *WorkVersionKey) Key() string {
	return s.builder.String()
}

func NewWorkObjectKey(workObjectId int64) *WorkObjectKey {
	v := &WorkObjectKey{}
	return v.WorkObject(workObjectId)
}

func NewWorkItemKey(workItemId int64) *WorkItemKey {
	u := &WorkItemKey{}
	return u.WorkItem(workItemId)
}
func NewUserKey(userId int64) *UserKey {
	u := &UserKey{}
	return u.User(userId)
}

func NewSpaceKey(spaceId int64) *SpaceKey {
	v := &SpaceKey{}
	return v.Space(spaceId)
}

func NewWorkVersionKey(workVersionId int64) *WorkVersionKey {
	v := &WorkVersionKey{}
	return v.WorkVersion(workVersionId)
}
