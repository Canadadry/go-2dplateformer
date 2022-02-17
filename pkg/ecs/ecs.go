package ecs

import (
	"sync"
)

type SystemId int

type EntityId int

type ComponentKind string

type System interface {
	MatchKinds(map[ComponentKind]interface{}) bool
	Update(map[ComponentKind]interface{})
}

type World struct {
	lastEntityId       EntityId
	availableEntityIds []EntityId
	systems            []System
	components         map[EntityId]map[ComponentKind]interface{}
	entitiesBySystem   map[SystemId][]EntityId
	mutex              sync.Mutex
}

func New() *World {
	return &World{
		components:       map[EntityId]map[ComponentKind]interface{}{},
		entitiesBySystem: map[SystemId][]EntityId{},
	}
}

func (w *World) AddEntity(cs map[ComponentKind]interface{}) EntityId {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	e, ok := w.getAvailableEntityId()
	if !ok {
		w.lastEntityId++
		e = w.lastEntityId
	}
	w.updateComponent_AddEntity(e, cs)
	w.updateSystem_AddEntity(e, cs)
	return e
}

func (w *World) getAvailableEntityId() (EntityId, bool) {
	available := len(w.availableEntityIds)
	if available == 0 {
		return 0, false
	}
	e := w.availableEntityIds[available-1]
	w.availableEntityIds = w.availableEntityIds[:available-1]
	return e, true
}

func (w *World) updateComponent_AddEntity(e EntityId, cs map[ComponentKind]interface{}) {
	bucket := map[ComponentKind]interface{}{}
	for k, c := range cs {
		bucket[k] = c
	}
	w.components[e] = bucket
}

func (w *World) updateSystem_AddEntity(e EntityId, cs map[ComponentKind]interface{}) {
	for id, s := range w.systems {
		sId := SystemId(id)
		if s.MatchKinds(cs) {
			w.entitiesBySystem[sId] = append(w.entitiesBySystem[sId], e)
		}
	}
}

func (w *World) RemoveEntity(e EntityId) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.availableEntityIds = append(w.availableEntityIds, e)
	delete(w.components, e)
	w.updateSystem_RemoveEntity(e)
}

func (w *World) updateSystem_RemoveEntity(e EntityId) {
	for id := range w.systems {
		sId := SystemId(id)
		w.entitiesBySystem[sId] = remove(w.entitiesBySystem[sId], e)
	}
}

func remove(s []EntityId, eId EntityId) []EntityId {
	found := false
	for i, e := range s {
		if e == eId {
			s[i] = s[len(s)-1]
			found = true
			break
		}
	}
	if found {
		return s[:len(s)-1]
	}
	return s
}

func (w *World) AddSystem(s System) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	sId := SystemId(len(w.systems))
	w.systems = append(w.systems, s)
	for entityId, cmpts := range w.components {
		if s.MatchKinds(cmpts) {
			w.entitiesBySystem[sId] = append(w.entitiesBySystem[sId], entityId)
		}
	}
}

func (w *World) Update() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	for sId, entities := range w.entitiesBySystem {
		for _, e := range entities {
			w.systems[int(sId)].Update(w.components[e])
		}
	}
}
