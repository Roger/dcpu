package main

import (
    "fmt"
    "github.com/tncardoso/gocurses"
    )

const (
    windowLines = 12
    windowColumns  = 32
)

type Display struct {
    main *gocurses.Window
    debug *gocurses.Window
    memory *gocurses.Window
    registers *gocurses.Window
    cpu *Cpu
}

func (display *Display) update_main() {
    display.main.Box(0, 0)
    display.main.Mvaddstr(0, 2, "DCPU-16")
    for i := int(0); i < windowLines*windowColumns; i++ {
        ch := int(display.cpu.Mem[0x8000+i-1]&0xff)
        if ch == 0 {
            ch = 96
        }
        display.main.Mvaddch((i/windowColumns)+1, (i%windowColumns)+1, rune(ch))
    }
    display.main.Refresh()
}

func (display *Display) log(str string) {
    display.debug.Addstr(str)
    display.debug.Refresh()
}

func (display *Display) update_debug() {
    //display.debug.Box(0, 0)
    //display.debug.Mvaddstr(0, 2, "Debug")
    display.debug.Refresh()
}

func (display *Display) update_memory() {
    display.memory.Box(0, 0)
    display.memory.Mvaddstr(0, 2, "Memory")
    display.memory.Refresh()
}

func (display *Display) update_registers() {
    display.registers.Box(0, 0)
    display.registers.Mvaddstr(0, 2, "Info/Registers")
    display.registers.Mvaddstr(1, 1, fmt.Sprintf("A: %X B: %X C: %X",
                         display.cpu.A, display.cpu.B, display.cpu.C))
    display.registers.Mvaddstr(2, 1, fmt.Sprintf("X: %X Y: %X Z: %X",
                         display.cpu.X, display.cpu.Y, display.cpu.Z))
    display.registers.Mvaddstr(3, 1, fmt.Sprintf("I: %X J: %X, O: %X",
                                 display.cpu.I, display.cpu.J, display.cpu.O))
    display.registers.Mvaddstr(4, 1, fmt.Sprintf("PC: %X SP: %X", 
        display.cpu.PC, display.cpu.PC))
    display.registers.Mvaddstr(5, 1, fmt.Sprintf("Cycles: %.f", display.cpu.cycles))
    display.registers.Refresh()
}

func NewDisplay(cpu *Cpu) *Display {
    display := &Display {cpu: cpu}

    gocurses.Initscr()
    defer gocurses.End()

    gocurses.Cbreak()
    gocurses.Noecho()

    gocurses.Stdscr.Keypad(false)
    gocurses.Refresh()

    // Main Window
    display.main = gocurses.NewWindow(windowLines+2, windowColumns+2, 1, 1)

    // Info Window
    display.registers = gocurses.NewWindow(windowLines/2+1, windowColumns+2,
                                                        1, windowColumns+3)

    // Memory Window
    display.memory = gocurses.NewWindow(windowLines/2+1, windowColumns+2,
                                        windowLines/2+2, windowColumns+3)

    //// Debug Window
    display.debug = gocurses.NewWindow(windowLines+2, 80, windowLines+3, 1)
    display.debug.Scrollok(true)

    // Refresh All Windows
    display.update_main()
    display.update_memory()
    display.update_registers()
    display.update_debug()

    //gocurses.Getch()
    return display
}
