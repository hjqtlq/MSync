package mongosync

import (
	"strings"
)

type Namespace struct {
	Ns   string
	Db   string
	Coll string
}

func ParseNamespace(ns string) (db, coll string) {
	temp := strings.SplitN(ns, ".", 2)
	if len(temp) == 2 {
		return temp[0], temp[1]
	}
	return "", ""
}

var nsm = &namespaceManager{}

type namespaceManager struct {
}

func (n *namespaceManager) shouldSkipNamespace(db, coll string) bool {
	return false
}
