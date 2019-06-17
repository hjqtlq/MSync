package mongosync

import "strings"

type Namespace struct {
	Ns   string
	Db   string
	Coll string
}

func NewNamespace(ns string) *Namespace {
	n := &Namespace{}
	temp := strings.Split(ns, ".")
	if len(temp) == 2 {
		n.Db = temp[0]
		n.Coll = temp[1]
	}

	return n
}

var nsm = &namespaceManager{}

type namespaceManager struct {
}

func (n *namespaceManager) shouldSkipNamespace(db, coll string) bool {
	return false
}
