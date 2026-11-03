// Code generated-adjacent ABI fix for the Hike float64 lowering gap.
#include "textflag.h"

TEXT ·WindowAttributesWithDimensionsFloat64(SB), 0, $32-24
	MOVQ windowAttributes+0(FP), CX
	// Windows x64 positions floating arguments by argument slot: the first
	// slot is the windowAttributes pointer, so width/height use XMM1/XMM2.
	MOVSD width+8(FP), X1
	MOVSD height+16(FP), X2
	MOVQ SP, R11
	ANDQ $15, R11
	SUBQ R11, SP
	SUBQ $48, SP
	MOVQ R11, 40(SP)
	CALL winit_window_attributes_with_dimensions(SB)
	MOVQ 40(SP), R11
	ADDQ $48, SP
	ADDQ R11, SP
	RET
