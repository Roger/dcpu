package main

import (
    "time"
    "fmt"
    )

const MEM_SIZE = 0x10000
const DEBUG = false

type Cpu struct {
    display *Display
    Mem [MEM_SIZE]uint16
    A,B,C,X,Y,Z,I,J,PC,SP,O uint16
    cycles float32
    registers map[uint16] *uint16
    instructions map[uint16] func(cpu *Cpu, a *uint16, b *uint16)
    skip bool
}


func (cpu *Cpu) Printf(format string, v ...interface{}) {
    if DEBUG {
        cpu.display.log(fmt.Sprintf(format, v...))
    }
}

func (cpu *Cpu) Log(format string, v ...interface{}) {
    cpu.display.log(fmt.Sprintf(format, v...))
}

func NewCpu() *Cpu {
    cpu := &Cpu {}
    display := NewDisplay(cpu)
    cpu.display = display
    cpu.registers = map[uint16] *uint16 {
          0x0: &cpu.A,
          0x1: &cpu.B,
          0x2: &cpu.C,
          0x3: &cpu.X,
          0x4: &cpu.Y,
          0x5: &cpu.Z,
          0x6: &cpu.I,
          0x7: &cpu.J,
    }
    cpu.instructions = map[uint16] func(cpu *Cpu, a *uint16, b *uint16) {
          0x1: (*Cpu).set,
          0x2: (*Cpu).add,
          0x3: (*Cpu).sub,
          0x4: (*Cpu).mul,
          0x5: (*Cpu).div,
          0x6: (*Cpu).mod,
          0x7: (*Cpu).shl,
          0x8: (*Cpu).shr,
          0x9: (*Cpu).and,
          0xa: (*Cpu).bor,
          0xb: (*Cpu).xor,
          0xc: (*Cpu).ife,
          0xd: (*Cpu).ifn,
          0xe: (*Cpu).ifg,
          0xf: (*Cpu).ifb,
          0x80 | 0x1: (*Cpu).jsr,
    }
    return cpu
}

func (cpu *Cpu) set(a *uint16, b *uint16) {
    cpu.Printf("SET [%X] [%X]\n", *a, *b)
    *a = *b
}
func (cpu *Cpu) add(a *uint16, b *uint16) {
    cpu.Printf("ADD [%X] [%X]\n", *a, *b)
    res := uint32(*a) + uint32(*b)
    // limit to a Word
    *a = uint16(res & 0xffff)
    // 0 or 1
    cpu.O = uint16((res >> 16) & 0xffff)
}
func (cpu *Cpu) sub(a *uint16, b *uint16) {
    cpu.Printf("SUB [%X] [%X]\n", *a, *b)
    res := *a - *b
    // catch undeflow
    *a = uint16(res & 0xffff)
    // 0 or 0xffff
    cpu.O = uint16((res >> 16) & 0xffff)
}
func (cpu *Cpu) mul(a *uint16, b *uint16) {
    cpu.Printf("MUL [%X] [%X]\n", *a, *b)
    *a = *a * *b
    cpu.O =((*a * *b)>>16)&0xffff
}
func (cpu *Cpu) div(a *uint16, b *uint16) {
    cpu.Printf("DIV [%X] [%X]\n", *a, *b)
    cpu.O = ((*a<<16) / *b) & 0xffff
    if *b == 0 {
        *a = cpu.O
    }
}
func (cpu *Cpu) mod(a *uint16, b *uint16) {
    cpu.Printf("MOD [%X] [%X]\n", *a, *b)
    *a = *a % *b
    if *b == 0 {
        cpu.O = 0
    }
}
func (cpu *Cpu) shl(a *uint16, b *uint16) {
    cpu.Printf("SHL [%X] [%X]\n", *a, *b)
    *a = *a << *b
    cpu.O = (*a >> 16)&0xffff
}
func (cpu *Cpu) shr(a *uint16, b *uint16) {
    cpu.Printf("SHR [%X] [%X]\n", *a, *b)
    cpu.O = ((*a<<16)>>*b)&0xffff
    *a = *a >> *b
}
func (cpu *Cpu) and(a *uint16, b *uint16) {
    cpu.Printf("AND [%X] [%X]\n", *a, *b)
    *a = *a & *b
}
func (cpu *Cpu) bor(a *uint16, b *uint16) {
    cpu.Printf("BOR [%X] [%X]\n", *a, *b)
    *a = *a | *b
}
func (cpu *Cpu) xor(a *uint16, b *uint16) {
    cpu.Printf("XOR [%X] [%X]\n", *a, *b)
    *a = *a ^ *b
}
func (cpu *Cpu) ife(a *uint16, b *uint16) {
    cpu.Printf("IFE [%X] [%X]\n", *a, *b)
    if *a != *b {
        cpu.skip = true
    }
}
func (cpu *Cpu) ifn(a *uint16, b *uint16) {
    cpu.Printf("IFN [%X] [%X]\n", *a, *b)
    if *a == *b {
        cpu.skip = true
    }
}
func (cpu *Cpu) ifg(a *uint16, b *uint16) {
    cpu.Printf("IFG [%X] [%X]\n", *a, *b)
    if *a <= *b {
        cpu.skip = true
    }
}
func (cpu *Cpu) ifb(a *uint16, b *uint16) {
    cpu.Printf("IFB [%X] [%X]\n", *a, *b)
    if *a >= *b {
        cpu.skip = true
    }
}
func (cpu *Cpu) jsr(a *uint16, b *uint16) {
    cpu.Printf("JSR [%X]\n", *b)
    cpu.SP--
    cpu.Mem[cpu.SP] = cpu.PC
    cpu.PC = *b
}

func (cpu *Cpu) fill_ram() {
    ram := []uint16{
     0x7c01, 0x0030, 0x7de1, 0x1000, 0x0020, 0x7803, 0x1000, 0xc00d,
     0x7dc1, 0x0041, 0xa861, 0x7c01, 0x2000, 0x2161, 0x2000, 0x8463,
     0x806d, 0x7dc1, 0x000d, 0x9031, 0x7c10, 0x0018, 0x7dc1, 0x001a,
     0x9037, 0x61c1, 0x7de1, 0x8000, 0x0048, 0x7de1, 0x8001, 0x0065,
     0x7de1, 0x8002, 0x006c, 0x7de1, 0x8003, 0x006c, 0x7de1, 0x8004,
     0x006f, 0x7de1, 0x8005, 0x002c, 0x7de1, 0x8006, 0x0020, 0x7de1,
     0x8007, 0x0077, 0x7de1, 0x8008, 0x006f, 0x7de1, 0x8009, 0x0072,
     0x7de1, 0x800a, 0x006c, 0x7de1, 0x800b, 0x0064, 0x7de1, 0x800c,
     0x0021, 0x7dc1, 0x0041, 0x0000, 0x0000, 0x0000, 0x0000, 0x0000,
    }
    var i int
    for i=0; i<len(ram); i++ {
        cpu.Mem[i] = ram[i]
    }
}

func (cpu *Cpu) read_ram() {
    prnt := false
    var line string
    for i := range cpu.Mem {
        if i == int(cpu.PC) {
            line += fmt.Sprintf("[%.4X] ", cpu.Mem[i])
        } else {
            line += fmt.Sprintf("%.4X ", cpu.Mem[i])
        }

        if cpu.Mem[i] > 0 {
            prnt = true
        }

        if int(i+1) % 8 == 0 {
            if prnt {
                //cpu.Log("%.4X: %s\n", i-7,line)
                prnt = false
            }
            line = ""
        }
    }
    //cpu.Log("\n\n")
}

func (cpu *Cpu) get_operand(addr uint16) *uint16{
    // register
    if(addr < 0x08){
        cpu.Printf("register [%X]\n", addr)
        return cpu.registers[addr]
    // [register]
    }else if(addr < 0x10){
        cpu.Printf("[register] [%X]\n", addr)
        return &cpu.Mem[*cpu.registers[addr-0x8]]
    // [next word + register]
    }else if(addr < 0x18){
        next_word := cpu.Mem[cpu.PC]
        cpu.PC++
        cpu.cycles++
        cpu.Printf("[next word + register] [%X] [%X]\n", addr, cpu.Mem[next_word+*cpu.registers[addr-0xf]])
        return &cpu.Mem[next_word+*cpu.registers[addr-0xf]]
    // POP / [SP++]
    }else if(addr == 0x18){
        cpu.Printf("[SP++] [%X] [%X]\n", addr, cpu.SP)
        ret := cpu.Mem[cpu.SP]
        cpu.SP++
        return &ret
    // PEEK / [SP]
    }else if(addr == 0x19){
        return &cpu.Mem[*cpu.registers[cpu.SP]]
    // PUSH / [--SP]
    }else if(addr == 0x1a){
        cpu.SP--
        return &cpu.Mem[*cpu.registers[cpu.SP]]
    // SP
    }else if(addr == 0x1b){
        return &cpu.SP
    // PC
    }else if(addr == 0x1c){
        cpu.Printf("[PC] [%X]\n", addr)
        return &cpu.PC
    // O
    }else if(addr == 0x1d){
        return &cpu.O
    // [next word]
    }else if(addr == 0x1e){
        cpu.Printf("[next_word] [%X]\n", addr)
        cpu.PC++
        return &cpu.Mem[cpu.Mem[cpu.PC-1]]
    // next word (literal)
    }else if(addr == 0x1f){
        cpu.Printf("next_word(literal) [%X]\n", addr)
        cpu.PC++
        return &cpu.Mem[cpu.PC-1]
    // literal value
    }else if(addr >= 0x20 && addr < 0x3f){
        cpu.Printf("literal val [%X]\n", addr)
        addr -= 0x20
        return &addr
    }
    panic("AAAAH")
}

func (cpu *Cpu) do_cycle() {
    inst := cpu.Mem[cpu.PC]
    cpu.PC++
    cpu.cycles++

    opcode := inst & 0xf

    var a *uint16
    var b *uint16

    a_addr := (inst & 0x3f0) >> 4
    b_addr := (inst & 0xfc00) >> 10

    // Non basic opcode
    if opcode == 0 {
        opcode = a_addr | 0x80
        cpu.Printf("Non Basic Op\n")
    } else {
        a = cpu.get_operand(a_addr)
    }

    cpu.Printf("aa:%X\n", a_addr)
    cpu.Printf("ba:%X\n", b_addr)

    b = cpu.get_operand(b_addr)

    if cpu.skip {
        cpu.cycles += 2
        cpu.skip = false
    } else {
        cpu.instructions[opcode](cpu, a, b)
    }
}

func main() {
    cpu := NewCpu()
    cpu.fill_ram()
    for i:=0;i<2000;i++ {
        cpu.do_cycle()
        time.Sleep(100 * time.Millisecond)
        //if curses.Getch() == 0x71 {
        //    // clean up screen
        //    curses.End()
        //    return
        //}
        cpu.display.update_registers()
        cpu.display.update_main()
        cpu.read_ram()
    }
}
