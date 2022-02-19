package ecs

import (
	"fmt"
	"testing"
)

type fakeSystem struct {
	MatchOnCountGreaterThan int
	Trace                   []string
}

func (f *fakeSystem) Match(e Entity) bool {
	return len(e) > f.MatchOnCountGreaterThan
}
func (f *fakeSystem) Update(e Entity) {
	f.Trace = append(f.Trace, fmt.Sprintf("update <%v>", e))
}
func (f *fakeSystem) Draw(e Entity, s interface{}) {
	f.Trace = append(f.Trace, fmt.Sprintf("draw <%v>", e))
}

func buildCmpt(count int) Entity {
	out := Entity{}
	for i := 0; i < count; i++ {
		out[ComponentKind(fmt.Sprintf("cmpt%d", i))] = struct{}{}
	}
	return out
}

func TestWorld_AddOne(t *testing.T) {
	w := New()
	w.AddEntity(buildCmpt(1))
	exp := ExpectWorld{
		LastEntityId:             1,
		AvailableEntityIds:       nil,
		LenForSystems:            0,
		LenForComponents:         1,
		LenForComponentsOfEntity: map[EntityId]int{1: 1},
		LenForEntitiesBySystem:   0,
	}
	testWorld(t, w, exp)
}

func TestWorld_AddTwo(t *testing.T) {
	w := New()
	w.AddEntity(buildCmpt(1))
	w.AddEntity(buildCmpt(3))
	exp := ExpectWorld{
		LastEntityId:             2,
		AvailableEntityIds:       nil,
		LenForSystems:            0,
		LenForComponents:         2,
		LenForComponentsOfEntity: map[EntityId]int{1: 1, 2: 3},
		LenForEntitiesBySystem:   0,
	}
	testWorld(t, w, exp)
}

func TestWorld_AddThreeAndRemoveSecond(t *testing.T) {
	w := New()
	w.AddEntity(buildCmpt(1))
	w.AddEntity(buildCmpt(3))
	w.AddEntity(buildCmpt(5))
	w.RemoveEntity(2)
	exp := ExpectWorld{
		LastEntityId:             3,
		AvailableEntityIds:       []EntityId{2},
		LenForSystems:            0,
		LenForComponents:         2,
		LenForComponentsOfEntity: map[EntityId]int{1: 1, 3: 5},
		LenForEntitiesBySystem:   0,
	}
	testWorld(t, w, exp)
}

func TestWorld_AddThreeAndRemoveSecondAndReAddOne(t *testing.T) {
	w := New()
	w.AddEntity(buildCmpt(1))
	w.AddEntity(buildCmpt(3))
	w.AddEntity(buildCmpt(5))
	w.RemoveEntity(2)
	w.AddEntity(buildCmpt(0))
	exp := ExpectWorld{
		LastEntityId:             3,
		AvailableEntityIds:       nil,
		LenForSystems:            0,
		LenForComponents:         3,
		LenForComponentsOfEntity: map[EntityId]int{1: 1, 2: 0, 3: 5},
		LenForEntitiesBySystem:   0,
	}
	testWorld(t, w, exp)
}

func TestWorld_WithFakeSystem_AddOne(t *testing.T) {
	w := New()
	w.AddSystem(&fakeSystem{MatchOnCountGreaterThan: 2})
	w.AddEntity(buildCmpt(1))
	exp := ExpectWorld{
		LastEntityId:             1,
		AvailableEntityIds:       nil,
		LenForSystems:            1,
		LenForComponents:         1,
		LenForComponentsOfEntity: map[EntityId]int{1: 1},
		LenForEntitiesBySystem:   0,
	}
	testWorld(t, w, exp)
}

func TestWorld_WithFakeSystem_AddTwo(t *testing.T) {
	w := New()
	w.AddSystem(&fakeSystem{MatchOnCountGreaterThan: 2})
	w.AddEntity(buildCmpt(1))
	w.AddEntity(buildCmpt(3))
	exp := ExpectWorld{
		LastEntityId:             2,
		AvailableEntityIds:       nil,
		LenForSystems:            1,
		LenForComponents:         2,
		LenForComponentsOfEntity: map[EntityId]int{1: 1, 2: 3},
		LenForEntitiesBySystem:   1,
	}
	testWorld(t, w, exp)
}

func TestWorld_WithFakeSystem_AddThreeAndRemoveSecond(t *testing.T) {
	w := New()
	w.AddSystem(&fakeSystem{MatchOnCountGreaterThan: 2})
	w.AddEntity(buildCmpt(1))
	w.AddEntity(buildCmpt(3))
	w.AddEntity(buildCmpt(5))
	w.RemoveEntity(2)
	exp := ExpectWorld{
		LastEntityId:             3,
		AvailableEntityIds:       []EntityId{2},
		LenForSystems:            1,
		LenForComponents:         2,
		LenForComponentsOfEntity: map[EntityId]int{1: 1, 3: 5},
		LenForEntitiesBySystem:   1,
	}
	testWorld(t, w, exp)
}

func TestWorld_WithFakeSystem_AddThreeAndRemoveSecondAndReAddOne(t *testing.T) {
	w := New()
	w.AddSystem(&fakeSystem{MatchOnCountGreaterThan: 2})
	w.AddEntity(buildCmpt(1))
	w.AddEntity(buildCmpt(3))
	w.AddEntity(buildCmpt(5))
	w.RemoveEntity(2)
	w.AddEntity(buildCmpt(0))
	exp := ExpectWorld{
		LastEntityId:             3,
		AvailableEntityIds:       nil,
		LenForSystems:            1,
		LenForComponents:         3,
		LenForComponentsOfEntity: map[EntityId]int{1: 1, 2: 0, 3: 5},
		LenForEntitiesBySystem:   1,
	}
	testWorld(t, w, exp)
}

func TestWorld_WithFakeSystem_AddThreeAndRemoveSecondAndReAddOne_ButSystemIsAddedAfter(t *testing.T) {
	w := New()
	w.AddEntity(buildCmpt(1))
	w.AddEntity(buildCmpt(3))
	w.AddEntity(buildCmpt(5))
	w.RemoveEntity(2)
	w.AddEntity(buildCmpt(0))
	w.AddSystem(&fakeSystem{MatchOnCountGreaterThan: 2})
	exp := ExpectWorld{
		LastEntityId:             3,
		AvailableEntityIds:       nil,
		LenForSystems:            1,
		LenForComponents:         3,
		LenForComponentsOfEntity: map[EntityId]int{1: 1, 2: 0, 3: 5},
		LenForEntitiesBySystem:   1,
	}
	testWorld(t, w, exp)
}

func TestWorldUpdate_WithFakeSystem_AddThreeAndRemoveSecondAndReAddOne_ButSystemIsAddedAfter(t *testing.T) {
	s := (&fakeSystem{MatchOnCountGreaterThan: 1})
	w := New()
	w.AddEntity(buildCmpt(1))
	w.AddEntity(buildCmpt(2))
	w.AddEntity(buildCmpt(3))
	w.AddSystem(s)
	w.Update()

	expected := []string{
		"update <map[cmpt0:{} cmpt1:{}]>",
		"update <map[cmpt0:{} cmpt1:{} cmpt2:{}]>",
	}
	expectedStr := fmt.Sprintf("%v", expected)
	resultStr := fmt.Sprintf("%v", s.Trace)
	if resultStr != expectedStr {
		t.Fatalf("s.UpdateTrace \nexp %s\ngot %s\n", expectedStr, resultStr)
	}
}

func TestWorldDraw_WithFakeSystem_AddThreeAndRemoveSecondAndReAddOne_ButSystemIsAddedAfter(t *testing.T) {
	s := (&fakeSystem{MatchOnCountGreaterThan: 1})
	w := New()
	w.AddEntity(buildCmpt(1))
	w.AddEntity(buildCmpt(2))
	w.AddEntity(buildCmpt(3))
	w.AddSystem(s)
	w.Draw(nil)

	expected := []string{
		"draw <map[cmpt0:{} cmpt1:{}]>",
		"draw <map[cmpt0:{} cmpt1:{} cmpt2:{}]>",
	}
	expectedStr := fmt.Sprintf("%v", expected)
	resultStr := fmt.Sprintf("%v", s.Trace)
	if resultStr != expectedStr {
		t.Fatalf("s.UpdateTrace \nexp %s\ngot %s\n", expectedStr, resultStr)
	}
}

type ExpectWorld struct {
	LastEntityId             EntityId
	AvailableEntityIds       []EntityId
	LenForSystems            int
	LenForComponents         int
	LenForComponentsOfEntity map[EntityId]int
	LenForEntitiesBySystem   int
}

func testWorld(t *testing.T, w *World, expected ExpectWorld) {
	t.Helper()

	if expected.LastEntityId != w.lastEntityId {
		t.Fatalf("lastEntityId should be equal to %d got %d", expected.LastEntityId, w.lastEntityId)
	}

	if len(expected.AvailableEntityIds) != len(w.availableEntityIds) {
		t.Fatalf("len(w.availableEntityIds) should be equal to %d got %d", len(expected.AvailableEntityIds), len(w.availableEntityIds))
	}
	if len(expected.AvailableEntityIds) > 0 {
		resultStr := fmt.Sprintf("%v", w.availableEntityIds)
		expectedStr := fmt.Sprintf("%v", expected.AvailableEntityIds)
		if resultStr != expectedStr {
			t.Fatalf("w.availableEntityIds \nexp %s\ngot %s\n", expectedStr, resultStr)
		}
	}

	if expected.LenForSystems != len(w.systems) {
		t.Fatalf("len(w.systems) should be equal to %d got %d", expected.LenForSystems, len(w.systems))
	}

	if expected.LenForComponents != len(w.entities) {
		t.Fatalf("len(w.entities) should be equal to %d got %d", expected.LenForComponents, len(w.entities))
	}

	for e, cmpts := range w.entities {
		if expected.LenForComponentsOfEntity[e] != len(cmpts) {
			t.Fatalf("len(cmpts) should be equal to %d got %d", expected.LenForComponentsOfEntity, len(cmpts))
		}
	}

	if expected.LenForEntitiesBySystem != len(w.entitiesBySystem) {
		t.Fatalf("len(w.entitiesBySystem) should be equal to %d got %d", expected.LenForEntitiesBySystem, len(w.entitiesBySystem))
	}
}

func TestRemove(t *testing.T) {
	tests := []struct {
		slice    []EntityId
		item     EntityId
		expected []EntityId
	}{
		{
			slice:    []EntityId{1, 2, 3, 4},
			item:     2,
			expected: []EntityId{1, 4, 3},
		},
		{
			slice:    []EntityId{4, 3, 2, 1},
			item:     2,
			expected: []EntityId{4, 3, 1},
		},
		{
			slice:    []EntityId{4, 3, 2, 1},
			item:     5,
			expected: []EntityId{4, 3, 2, 1},
		},
	}

	for i, tt := range tests {
		result := remove(tt.slice, tt.item)
		resultStr := fmt.Sprintf("%v", result)
		expectedStr := fmt.Sprintf("%v", tt.expected)
		if resultStr != expectedStr {
			t.Fatalf("%d failed \nexp %s\ngot %s\n", i, expectedStr, resultStr)
		}
	}
}
