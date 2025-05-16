package job

type JobVariabless []*JobVariables

type JobVariables struct {
	JobName       string `json:"job_name" gorm:"uniqueIndex,column:job_name"`
	VariableName  string `json:"variable_name" gorm:"uniqueIndex,column:variable_name"`
	VariableValue string `json:"variable_value" gorm:"column:variable_value"`
}
