package test

import "go-cs/internal/dwh/pkg"

type testJob struct {
	id   string
	name string
}

func (m *testJob) Id() string {
	return m.id
}

func (m *testJob) Name() string {
	return m.name
}

func (m *testJob) FullName() string {
	return m.name + ":" + m.id
}

func (m *testJob) Run() {}

func MockTestJob(id string, name string) pkg.Job {
	return &testJob{
		id:   id,
		name: name,
	}
}
