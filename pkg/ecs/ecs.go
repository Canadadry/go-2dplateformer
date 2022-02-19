package ecs

import (
	"sync"
)

type ComponentKind string
type Entity map[ComponentKind]interface{}

type SystemId int

type EntityId int

type System interface {
	Match(Entity) bool
	Update(Entity)
	Draw(Entity)
}

type World struct {
	lastEntityId       EntityId
	availableEntityIds []EntityId
	systems            []System
	entities           map[EntityId]Entity
	entitiesBySystem   map[SystemId][]EntityId
	mutex              sync.Mutex
}

func New() *World {
	return &World{
		entities:         map[EntityId]Entity{},
		entitiesBySystem: map[SystemId][]EntityId{},
	}
}

func (w *World) AddEntity(e Entity) EntityId {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	eId, ok := w.getAvailableEntityId()
	if !ok {
		w.lastEntityId++
		eId = w.lastEntityId
	}
	w.entities[eId] = e
	w.updateSystem_AddEntity(eId, e)
	return eId
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

func (w *World) updateSystem_AddEntity(e EntityId, cs Entity) {
	for id, s := range w.systems {
		sId := SystemId(id)
		if s.Match(cs) {
			w.entitiesBySystem[sId] = append(w.entitiesBySystem[sId], e)
		}
	}
}

func (w *World) RemoveEntity(e EntityId) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.availableEntityIds = append(w.availableEntityIds, e)
	delete(w.entities, e)
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
	for entityId, entity := range w.entities {
		if s.Match(entity) {
			w.entitiesBySystem[sId] = append(w.entitiesBySystem[sId], entityId)
		}
	}
}

func (w *World) Update() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	for i, s := range w.systems {
		entities, _ := w.entitiesBySystem[SystemId(i)]
		for _, e := range entities {
			s.Update(w.entities[e])
		}
	}
}

func (w *World) Draw() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	for i, s := range w.systems {
		entities, _ := w.entitiesBySystem[SystemId(i)]
		for _, e := range entities {
			s.Draw(w.entities[e])
		}
	}
}
