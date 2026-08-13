package vm

type Op byte

const (
	OpHalt  Op = iota // stop Execute
	OpPush            // operand: u32 literal
	OpPop             // discard top of stack
	OpAdd             // pop a, pop b, push b+a
	OpLoad            // pop addr, push mem[addr]
	OpStore           // pop val, pop addr, mem[addr]=val
	OpAlloc           // pop nbytes, push ptr
	OpFree            // pop ptr
	OpCall            // operand: target ip; nargs already on stack per your convention
	OpRet             // restore FP/SP; jump to saved return addr
	OpDup             // duplicate top of stack
	OpPrint           // pop and log — demo only
)

type Instr struct {
	Op  Op
	Arg uint32
}
