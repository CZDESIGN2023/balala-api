package domain

type PropDiff string
type PropDiffs []PropDiff

type PropDiffSet map[PropDiff]any

func (s PropDiffSet) HasProp(pd PropDiff) bool {
	_, ok := s[pd]
	return ok
}
