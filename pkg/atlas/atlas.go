package atlas

import (
	"encoding/xml"
	"io"
)

type Frame struct {
	X      int
	Y      int
	Width  int
	Height int
}

func Load(r io.Reader) (map[string]Frame, error) {
	t := struct {
		XMLName     string `xml:"TextureAtlas"`
		ImagePath   string `xml:"imagePath,attr"`
		SubTextures []struct {
			Name   string `xml:"name,attr"`
			X      int    `xml:"x,attr"`
			Y      int    `xml:"y,attr"`
			Width  int    `xml:"width,attr"`
			Height int    `xml:"height,attr"`
		} `xml:"SubTexture"`
	}{}
	err := xml.NewDecoder(r).Decode(&t)
	if err != nil {
		return nil, err
	}
	ret := map[string]Frame{}
	for _, f := range t.SubTextures {
		ret[f.Name] = Frame{X: f.X, Y: f.Y, Width: f.Width, Height: f.Height}
	}
	return ret, err
}
