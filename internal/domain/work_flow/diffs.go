package work_flow

import shared "go-cs/internal/pkg/domain"

var (
	Diff_Ranking      shared.PropDiff = "ranking"
	Diff_Name         shared.PropDiff = "name"
	Diff_Status       shared.PropDiff = "status"
	Diff_LastTemplate shared.PropDiff = "lastTemplate"
	Diff_Version      shared.PropDiff = "version"

	Diff_DeletedAt shared.PropDiff = "deletedAt"

	// 以下为模版字段
	Diff_TemplateSetting shared.PropDiff = "setting"
)
