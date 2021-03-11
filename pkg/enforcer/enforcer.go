package enforcer

func Enforce(subject LoggedInUser, permission Permission, object ...interface{}) bool {
	if len(object) == 0 {
		return Can(subject.Role, permission)
	}

	return Can(subject.Role, permission) && AbacEvaluate(subject, permission, object[0])
}
