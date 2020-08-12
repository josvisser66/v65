package asm

type opcodeMap map[string][]int

var opcodes = opcodeMap{
	"adc": {-1, -1, 0x6d, 0x65, 0x69, 0x7d, 0x79, 0x61, 0x71, 0x75, -1, -1, -1},
	"and": {-1, -1, 0x2d, 0x25, 0x29, 0x3d, 0x39, 0x21, 0x31, 0x35, -1, -1, -1},
	"asl": {-1, 0x0a, 0x0e, 0x06, -1, 0x1e, -1, -1, -1, 0x16, -1, -1, -1},
	"bcc": {-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 0x90, -1},
	"bcs": {-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 0x80, -1},
	"beq": {-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 0xf0, -1},
	"bit": {-1, -1, 0x2c, 0x24, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	"bmi": {-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 0x30, -1},
	"bne": {-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 0xD0, -1},
	"bpl": {-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 0x10, -1},
	"brk": {0x00, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	"bvc": {-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 0x50, -1},
	"bvs": {-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 0x70, -1},
	"clc": {0x18, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	"cld": {0xd8, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	"cli": {0x58, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	"clv": {0xb8, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	"cmp": {-1, -1, 0xcd, 0xc5, 0xc9, 0xdd, 0xd9, 0xc1, 0xd1, 0xd5, -1, -1, -1},
	"cpx": {-1, -1, 0xec, 0xe4, 0xe0, -1, -1, -1, -1, -1, -1, -1, -1},
	"cpy": {-1, -1, 0xcc, 0xc4, 0xc0, -1, -1, -1, -1, -1, -1, -1, -1},
	"dec": {-1, -1, 0xce, 0xc6, -1, 0xde, -1, -1, -1, 0xd6, -1, -1, -1},
	"dex": {0xca, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	"dey": {0x88, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	"eor": {-1, -1, 0x4d, 0x45, 0x49, 0x5d, 0x59, 0x41, 0x51, 0x55, -1, -1, -1},
	"inc": {-1, -1, 0xee, 0xe6, -1, 0xfe, -1, -1, -1, 0xf6, -1, -1, -1},
	"inx": {0xe8, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	"iny": {0xc8, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	"jmp": {-1, -1, 0x4c, -1, -1, -1, -1, -1, -1, -1, -1, -1, 0x6c},
	"jsr": {-1, -1, 0x20, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	"lda": {-1, -1, 0xad, 0xa5, 0xa9, 0xbd, 0xb9, 0xa1, 0xb1, 0xb5, -1, -1, -1},
	"ldx": {-1, -1, 0xae, 0xa6, 0xa2, -1, 0xbe, -1, -1, -1, 0xb6, -1, -1},
	"ldy": {-1, -1, 0xac, 0xa4, 0xa0, 0xbc, -1, -1, -1, 0xb4, -1, -1, -1},
	"lsr": {-1, 0x4a, 0x4e, 0x46, -1, 0x5e, -1, -1, -1, 0x56, -1, -1, -1},
	"nop": {0xea, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	"ora": {-1, -1, 0x0d, 0x05, 0x09, 0x1d, 0x19, 0x01, 0x11, 0x15, -1, -1, -1},
	"pha": {0x48, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	"php": {0x08, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	"pla": {0x68, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	"plp": {0x28, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	"rol": {-1, 0x2a, 0x2e, 0x26, -1, 0x3e, -1, -1, -1, 0x36, -1, -1, -1},
	"ror": {-1, 0x6a, 0x6e, 0x66, -1, 0x7e, -1, -1, -1, 0x76, -1, -1, -1},
	"rti": {0x40, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	"rts": {0x60, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	"sbc": {-1, -1, 0xed, 0xe5, 0xe9, 0xfd, 0xf9, 0xe1, 0xf1, 0xf5, -1, -1, -1},
	"sec": {0x38, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	"sed": {0xf8, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	"sei": {0x78, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	"sta": {-1, -1, 0x8d, 0x85, -1, 0x9d, 0x99, 0x81, 0x91, 0x95, -1, -1, -1},
	"stx": {-1, -1, 0x8e, 0x86, -1, -1, -1, -1, -1, -1, 0x96, -1, -1},
	"sty": {-1, -1, 0x8c, 0x84, -1, -1, -1, -1, -1, 0x94, -1, -1, -1},
	"tax": {0xaa, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	"tay": {0xa8, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	"tsx": {0xba, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	"txa": {0x8a, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	"txs": {0x9a, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	"tya": {0x98, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
}
