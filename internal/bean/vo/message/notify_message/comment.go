package notify_message

// 某人 (someone)
func NewCommentOper(in *UserData) *Subject {
	subject := newSubject()
	subject.Type = SubjectType_user
	subject.Data = in
	return subject
}

// 某任务
func NewCommentWorkItemObject(in *WorkItemData) *Object {
	object := newObject()
	object.Type = ObjectType_workItem
	object.Data = in
	return object
}

// 某评论
func NewCommentObject(in *WorkItemCommentData) *Object {
	subject := newObject()
	subject.Type = ObjectType_workItemComment
	subject.Data = in
	return subject
}
