package data

import (
	job_model "go-cs/internal/dwh/model/job"
)

type JobVariablesRepo struct {
	data *DwhData
}

func NewJobVariablesRepo(
	data *DwhData,
) *JobVariablesRepo {
	return &JobVariablesRepo{
		data: data,
	}
}

func (j *JobVariablesRepo) GetVariablesByName(jobName string, variableName string) (*job_model.JobVariables, error) {
	var v *job_model.JobVariables

	err := j.data.Db().Table("job_variables").
		Where(job_model.JobVariables{JobName: jobName, VariableName: variableName}).
		Attrs(job_model.JobVariables{VariableValue: ""}).
		FirstOrCreate(&v).Error
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (j *JobVariablesRepo) SaveVariables(jobVariables *job_model.JobVariables) error {
	err := j.data.Db().Table("job_variables").Where("job_name = ? AND variable_name = ?", jobVariables.JobName, jobVariables.VariableName).Save(&jobVariables).Error
	if err != nil {
		return nil
	}
	return nil
}
