package atlas

import (
	"reflect"
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	in := `<TextureAtlas>
    	<SubTexture name="alienBeige_climb1" x="910" y="1290" width="128" height="256"/>
    	<SubTexture name="alienBeige_climb2" x="910" y="1032" width="128" height="256"/>
    	<SubTexture name="alienBeige_duck" x="910" y="774" width="128" height="256"/>
</TextureAtlas>`
	expected := map[string]Frame{
		"alienBeige_climb1": Frame{X: 910, Y: 1290, Width: 128, Height: 256},
		"alienBeige_climb2": Frame{X: 910, Y: 1032, Width: 128, Height: 256},
		"alienBeige_duck":   Frame{X: 910, Y: 774, Width: 128, Height: 256},
	}

	result, err := Load(strings.NewReader(in))
	if err != nil {
		t.Fatalf("failed %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("\ngot %#v\nexp %#v", result, expected)
	}
}

func TestLoad_Err(t *testing.T) {
	in := `<TextureAtlas>
    	<SubTexture name="alienBeige_climb1" x="910" y="1290" width="128" height="256"/>
    	<SubTexture name="alienBeige_cl`
	_, err := Load(strings.NewReader(in))
	if err == nil {
		t.Fatalf("should have failed")
	}
}
