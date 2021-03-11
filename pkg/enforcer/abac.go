package enforcer

import "reflect"

type evaluatorFunc func(subject LoggedInUser, object interface{}) bool

var abacPermissionEvaluator map[Permission]evaluatorFunc

func init() {
	abacPermissionEvaluator = make(map[Permission]evaluatorFunc)
	abacPermissionEvaluator[ReadContribution] = func(subject LoggedInUser, object interface{}) bool {
		switch subject.Role {
		case Student:
			userId := reflect.ValueOf(object).FieldByName("UserId").Int()
			return subject.Id == int(userId)
		case MarketingCoordinator:
			facultyId := reflect.ValueOf(object).FieldByName("User").FieldByName("FacultyId").Int()
			return *subject.FacultyId == int(facultyId)
		default:
			return false
		}
	}
}

func AbacEvaluate(subject LoggedInUser, permission Permission, object interface{}) bool {
	if v, ok := abacPermissionEvaluator[permission]; ok {
		return v(subject, object)
	} else {
		return false
	}
}
