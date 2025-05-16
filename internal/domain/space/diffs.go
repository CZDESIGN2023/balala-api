package space

import shared "go-cs/internal/pkg/domain"

var (
	Diff_SpaceName shared.PropDiff = "spaceName"
	Diff_Describe  shared.PropDiff = "describe"
	Diff_UserId    shared.PropDiff = "userId"
	Diff_Notify    shared.PropDiff = "notify"

	Diff_SpaceConfig_WorkingDay                   shared.PropDiff = "spaceConfig.workingDay"
	Diff_SpaceConfig_CommentDeletable             shared.PropDiff = "spaceConfig.commentDeletable"
	Diff_SpaceConfig_CommentDeletableWhenArchived shared.PropDiff = "spaceConfig.commentDeletableWhenArchived"
	Diff_SpaceConfig_CommentShowPos               shared.PropDiff = "spaceConfig.commentShowPos"
)
