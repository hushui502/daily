package codegen

import . "luago/binchunk"

func toProto(fi *funcInfo) *Prototype {
	proto := &Prototype{
		NumParams:    byte(fi.numParams),
		MaxStackSize: byte(fi.maxRegs),
		Code:         fi.insts,
		Constants:    getConstant(fi),
		Upvalues:     getUpvalues(fi),
		Protos:       toProtos(fi.subFuncs),
		LineInfo:     []uint32{},
		LocVars:      []LocVar{},
		UpvalueNames: []string{},
	}

	if proto.MaxStackSize < 2 {
		proto.MaxStackSize = 2
	}
	if fi.isVararg {
		proto.IsVararg = 1
	}

	return proto
}

func toProtos(fis []*funcInfo) []*Prototype {
	protos := make([]*Prototype, len(fis))
	for i, fi := range fis {
		protos[i] = toProto(fi)
	}

	return protos
}

func getConstant(fi *funcInfo) []interface{} {
	consts := make([]interface{}, len(fi.constants))
	for k, idx := range fi.constants {
		consts[idx] = k
	}

	return consts
}

func getUpvalues(fi *funcInfo) []Upvalue {
	upvals := make([]Upvalue, len(fi.upvalues))
	for _, uv := range fi.upvalues {
		if uv.locVarSlot >= 0 {
			upvals[uv.index] = Upvalue{1, byte(uv.locVarSlot)}
		} else {
			upvals[uv.index] = Upvalue{0, byte(uv.locVarSlot)}
		}
	}

	return upvals
}
