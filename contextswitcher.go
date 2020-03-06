package main

import (
    //"log"
	"github.com/faiface/pixel/pixelgl"
)

type ContextSwitcher struct {
    Contexts map[string]Context
    //Current Context
    //Saved Context
    Stack []Context
    InMenu bool
}

func (self *ContextSwitcher) Pop() {
    if len(self.Stack) > 0 {
        self.Stack = self.Stack[1:]
    }
}

func (self *ContextSwitcher) PopMenu() {
    for self.InMenu {
        if self.Current() != self.Contexts["menu"] {
            self.Pop()
        } else {
            self.Pop()
            self.InMenu = false
        }
    }
}

func (self *ContextSwitcher) Current() Context {
    return self.Stack[0]
}

func (self *ContextSwitcher) Switch(win *pixelgl.Window) {
    // ESC is a context switcher
    if win.JustPressed(pixelgl.KeyEscape) {
        if !self.InMenu {
            self.Stack = append([]Context{self.Contexts["menu"]}, self.Stack...)
            self.InMenu = true
        } else {
            if !self.Stack[0].HandleEscape() {
                if self.Stack[0] == self.Contexts["menu"] {
                    self.Stack = self.Stack[1:]
                    self.InMenu = false
                }
            }
        }
    }

    // c switches from main context to crafting context (menu context ignores c)
    //if win.JustPressed(pixelgl.KeyC) {
    //    if self.Current() == self.Contexts["main"] {
    //        self.Stack = append([]Context{self.Contexts["crafting"]}, self.Stack...)
    //    } else if self.Current() == self.Contexts["crafting"] {
    //        self.Stack = self.Stack[1:]
    //    }
    //}
}

func NewContextSwitcher(ctxts map[string]Context) *ContextSwitcher {
    main := ctxts["main"]
    //save := ctxts["main"]
    stack := []Context{main}
    return &ContextSwitcher{ctxts, stack, false}
}
