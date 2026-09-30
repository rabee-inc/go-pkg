package codegen

const errorInvalidKind = `invalid kind: %v`
const errorInvalidDefs = `defs should be specified as follows
defs:
	var_key: name value

or

defs:
	var_key:
		name: name value
		other: value`
