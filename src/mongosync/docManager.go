package mongosync

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"mongosync/managers"
	"sync"
)

var DocManager = &docManager{mutex: &sync.Mutex{}, managers: make(map[string]managers.IDocManager)}

type docManager struct {
	managers map[string]managers.IDocManager
	mutex    *sync.Mutex
	shardId  string

	buffChan chan oplog
}

func (dm *docManager) init() {
	go func() {
		for {
			select {
			case entry := <-dm.buffChan:
				switch entry.Op {
				case OplogOpDelete:
					if id, ok := entry.O["_id"]; ok {
						dm.Remove(id.(primitive.ObjectID), entry.Ns, entry.Ts)
					}
				case OplogOpInsert:
					dm.Upsert(entry.O, entry.Ns, entry.Ts)
				case OplogOpUpdate:
					if id, ok := entry.O2["_id"].(primitive.ObjectID); ok {
						dm.Update(id, entry.O, entry.Ns, entry.Ts)
					}
				case OplogOpCommand:
					dm.HandCommand(entry.O, entry.Ns, entry.Ts)
				}
			}
		}
	}()
}

func (dm *docManager) Process(entry oplog) {
	dm.buffChan <- entry
	Checkpoint.set(dm.shardId, entry.Ts)
}

func (dm *docManager) BulkUpsert(docs []interface{}, ns string, ts *primitive.Timestamp) {
	for _, manager := range dm.managers {
		manager.BulkUpsert(docs, ns, ts)
	}
}

func (dm *docManager) Update(docId primitive.ObjectID, updateSpec interface{}, ns string, ts *primitive.Timestamp) {
	for _, manager := range dm.managers {
		manager.Update(docId, updateSpec, ns, ts)

	}
}
func (dm *docManager) Upsert(doc interface{}, ns string, ts *primitive.Timestamp) {
	for _, manager := range dm.managers {
		manager.Upsert(doc, ns, ts)
	}
}
func (dm *docManager) Remove(docId primitive.ObjectID, ns string, ts *primitive.Timestamp) {
	for _, manager := range dm.managers {
		manager.Remove(docId, ns, ts)
	}
}
func (dm *docManager) HandCommand(doc interface{}, ns string, ts *primitive.Timestamp) {
	for _, manager := range dm.managers {
		manager.HandCommand(doc, ns, ts)
	}
}
func (dm *docManager) Commit() {
	for _, manager := range dm.managers {
		manager.Commit()
	}
}
func (dm *docManager) Stop() {
	for name, manager := range dm.managers {
		log.Println("Stopping doc manager:", name)
		manager.Stop()
	}
}

func (dm *docManager) Register(managers ...managers.IDocManager) {
	for _, manager := range managers {
		if _, ok := Config.DocManagers[manager.GetName()]; ok {
			manager.Init()
			dm.managers[manager.GetName()] = manager
		} else {
			log.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		}
	}
}

func (dm *docManager) getManagers() map[string]managers.IDocManager {
	ms := make(map[string]managers.IDocManager, len(DocManager.managers))
	for name, manager := range DocManager.managers {
		manager.SetShardId(dm.shardId)
		ms[name] = manager
	}
	return ms
}

func NewDocManager(shardId string) *docManager {
	dm := &docManager{
		shardId:  shardId,
		buffChan: make(chan oplog, 1000),
		managers: DocManager.getManagers(),
	}
	dm.init()
	return dm
}
