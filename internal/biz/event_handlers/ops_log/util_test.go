package ops_log

import (
	"go-cs/internal/utils"
	"testing"
)

const htmlExample = `
<div style="color: #cccccc; font-weight: normal;">
<div><span style="color: #808080;">&lt;</span><span style="color: #569cd6;">script</span><span style="color: #cccccc;"> </span><span style="color: #9cdcfe;">lang</span><span style="color: #cccccc;">=</span><span style="color: #ce9178;">"ts"</span><span style="color: #cccccc;"> </span><span style="color: #9cdcfe;">setup</span><span style="color: #808080;">&gt;</span></div>
<div><span style="color: #c586c0;">import</span><span style="color: #cccccc;"> </span><span style="color: #c586c0;">type</span><span style="color: #cccccc;"> {</span></div>
<div><span style="color: #cccccc;"> </span><span style="color: #9cdcfe;">ConditionGroup</span><span style="color: #cccccc;">,</span></div>
<div><span style="color: #cccccc;"> </span><span style="color: #9cdcfe;">QueryConditionGroup</span><span style="color: #cccccc;">,</span></div>
<div><span style="color: #cccccc;">} </span><span style="color: #c586c0;">from</span><span style="color: #cccccc;"> </span><span style="color: #ce9178;">'./interface'</span></div>
<div><span style="color: #c586c0;">import</span><span style="color: #cccccc;"> { </span><span style="color: #9cdcfe;">getFilterTag</span><span style="color: #cccccc;">, </span><span style="color: #9cdcfe;">setFilterTag</span><span style="color: #cccccc;"> } </span><span style="color: #c586c0;">from</span><span style="color: #cccccc;"> </span><span style="color: #ce9178;">'@/api/project'</span></div>
<div><span style="color: #c586c0;">import</span><span style="color: #cccccc;"> { </span><span style="color: #9cdcfe;">isInternal</span><span style="color: #cccccc;"> } </span><span style="color: #c586c0;">from</span><span style="color: #cccccc;"> </span><span style="color: #ce9178;">'@/config/internal'</span></div>
<div><span style="color: #c586c0;">import</span><span style="color: #cccccc;"> { </span><span style="color: #9cdcfe;">useProjectKey</span><span style="color: #cccccc;"> } </span><span style="color: #c586c0;">from</span><span style="color: #cccccc;"> </span><span style="color: #ce9178;">'@/hooks/useProjectCacheKey'</span></div>
<div><span style="color: #c586c0;">import</span><span style="color: #cccccc;"> { </span><span style="color: #9cdcfe;">compareFilter</span><span style="color: #cccccc;">, </span><span style="color: #9cdcfe;">filterInValidConditionGroup</span><span style="color: #cccccc;">, </span><span style="color: #9cdcfe;">formatConditionData</span><span style="color: #cccccc;">, </span><span style="color: #9cdcfe;">generateUniqueConditionName</span><span style="color: #cccccc;">, </span><span style="color: #9cdcfe;">getConditionsCount</span><span style="color: #cccccc;">, </span><span style="color: #9cdcfe;">getFirstValidConditionOfGroup</span><span style="color: #cccccc;">, </span><span style="color: #9cdcfe;">parseToUIData</span><span style="color: #cccccc;">, </span><span style="color: #9cdcfe;">shouldIgnoreClick</span><span style="color: #cccccc;"> } </span><span style="color: #c586c0;">from</span><span style="color: #cccccc;"> </span><span style="color: #ce9178;">'@/mixins/condition'</span></div>
<div><span style="color: #c586c0;">import</span><span style="color: #cccccc;"> </span><span style="color: #9cdcfe;">router</span><span style="color: #cccccc;"> </span><span style="color: #c586c0;">from</span><span style="color: #cccccc;"> </span><span style="color: #ce9178;">'@/router'</span></div>
<div><span style="color: #c586c0;">import</span><span style="color: #cccccc;"> { </span><span style="color: #9cdcfe;">FilterManager</span><span style="color: #cccccc;"> } </span><span style="color: #c586c0;">from</span><span style="color: #cccccc;"> </span><span style="color: #ce9178;">'@/services/filterData/FilterManager'</span></div>
</div>
`

func Test(t *testing.T) {
	tag := utils.ClearRichTextToPlanText(htmlExample, true)

	t.Log(tag)
}
