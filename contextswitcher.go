package main

import (
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
        //l := len(self.Stack)
        if !self.InMenu {
            self.Stack = append([]Context{self.Contexts["menu"]}, self.Stack...)
            //self.Stack = append(self.Stack, self.Contexts["menu"])
            self.InMenu = true
        } else {
            if self.Stack[0] == self.Contexts["menu"] {
                self.Stack = self.Stack[1:]
                self.InMenu = false
            }
        }
        //l = len(self.Stack)
        //self.Current = self.Stack[0]
    }
    /*
    if win.JustPressed(pixelgl.KeyEscape) {
        if self.Current != self.Contexts["menu"] {
            // saved previous context
            self.Saved = self.Current
            // start menu context
            self.Current = self.Contexts["menu"]
        } else {
            // load saved context
            self.Current = self.Saved
        }
    }
    */
    if win.JustPressed(pixelgl.KeyC) {
        if self.Current() == self.Contexts["main"] {
            self.Stack = append([]Context{self.Contexts["crafting"]}, self.Stack...)
            //self.Saved = self.Current
            //self.Current = self.Contexts["crafting"]
        } else if self.Current() == self.Contexts["crafting"] {
            self.Stack = self.Stack[1:]
            //self.Current = self.Saved
        }
        //self.Current = self.Stack[0]
    }
}

func NewContextSwitcher(ctxts map[string]Context) *ContextSwitcher {
    main := ctxts["main"]
    //save := ctxts["main"]
    stack := []Context{main}
    return &ContextSwitcher{ctxts, stack, false}
}
