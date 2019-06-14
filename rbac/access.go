package rbac

import (
	"regexp"
	"sync"
)

type Role string

type AccessList []Access

// accessCacheMap holds the accesses of a specified role
var accessCacheMap sync.Map

// Access controls the access to the path of the service by a specified http method
type Access struct {
	PathPattern string
	Method      string
}

// RegisterAccess registers new access list for a specified role
func RegisterAccess(role Role, accessList AccessList) {
	accessCacheMap.Store(role, accessList)
}

// CheckAccess checks whether the request access is valid or not, which returns error only when regexp pattern is invalid
func CheckAccess(role Role, reqPath string, reqMethod string) (valid bool, err error) {
	value, ok := accessCacheMap.Load(role)
	if !ok {
		// no access for the role
		return
	}

	accessList := value.(AccessList)
	var matched bool
	for _, access := range accessList {
		matched, err = regexp.MatchString(access.PathPattern, reqPath)
		if err != nil {
			return
		}

		if matched && access.Method == reqMethod {
			valid = true
			return
		}
	}

	return
}
