package main

import (
    "image/color"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
    "github.com/faiface/pixel/text"
)

type Menu struct {
    Options []string
    Active int
    Color color.RGBA
    ActiveColor color.RGBA
    Atlas *text.Atlas
    Text *text.Text
}

func NewMenu(atlas *text.Atlas, options []string, color, activeColor color.RGBA) *Menu {
    return &Menu{options, 0, color, activeColor, atlas, nil}
}

func (self *Menu) Write() {
    optText := text.New(pixel.ZV, self.Atlas)
    optText.Color = self.Color
    for i, opt := range self.Options {
        if i == self.Active {
            optText.Color = self.ActiveColor
        } else {
            optText.Color = self.Color
        }
        optText.WriteString(opt)
        optText.WriteString("\n")
    }
    self.Text = optText
}

func (self *Menu) Draw(win *pixelgl.Window, mat pixel.Matrix) {
    self.Text.Draw(win, mat)
}


func (self *Menu) Down() {
    self.Active = (self.Active + 1) % len(self.Options)
    if self.Active < 0 {
        self.Active += len(self.Options)
    }
}

func (self *Menu) Up() {
    self.Active = (self.Active - 1) % len(self.Options)
    if self.Active < 0 {
        self.Active += len(self.Options)
    }

}

