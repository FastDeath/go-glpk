package glpk

// #cgo LDFLAGS: -lglpk
// #include <glpk.h>
// #include <stdlib.h>
import "C"
import (
	"unsafe"
)

// SolType is used to specify which solution should be copied from the problem object to the workspace.
type SolType int

// Allowed values of type SolType (solution to be copied to workspace).
const (
	SOL = SolType(C.GLP_SOL) // basic solution
	IPT = SolType(C.GLP_IPT) // interior-point solution
	MIP = SolType(C.GLP_MIP) // mixed integer solution
)

type workspace struct {
	w *C.glp_tran
}

// Workspace represenst a MathProg translator workspace. Use mpl.New() to create a new workspace.
type Workspace struct {
	w *workspace
}

// NewMPL creates a new optimization problem.
func NewMPL() *Workspace {
	p := &workspace{C.glp_mpl_alloc_wksp()}
	return &Workspace{p}
}

// Free frees all the memory allocated to the translator workspace.
func (w *Workspace) Free() {
	if w.w.w != nil {
		C.glp_mpl_free_wksp(w.w.w)
		w.w.w = nil
	}
}

// InitRand initializes a pseudo-random number generator used by theMathProg translator.
func (w *Workspace) InitRand(seed int) {
	if w.w.w == nil {
		panic("MathProg method called on a deleted workspace")
	}
	C.glp_mpl_init_rand(w.w.w, C.int(seed))
}

// ReadModel reads model section and, optionally, data section, which may follow the model section,
// from a text file, whose name is the character string fname, performs translation of model statements
// and data blocks, and stores all the information in the workspace.
// The parameter skip is a flag. If the input file contains the data section and this flag is true,
// the data section is not read as if there were no data section and a warning message is printed.
// This allows reading data section(s) from other file(s).
func (w *Workspace) ReadModel(filename string, skip bool) {
	if w.w.w == nil {
		panic("MathProg method called on a deleted workspace")
	}

	fname := C.CString(filename)
	defer C.free(unsafe.Pointer(fname))

	var iskip int
	if skip {
		iskip = 1
	} else {
		iskip = 0
	}

	if C.glp_mpl_read_model(w.w.w, fname, C.int(iskip)) != 0 {
		panic("glp_mpl_read_model failed")
	}
}

// ReadData reads data section from a text file, whose name is the character string fname,
// performs translation of data blocks, and stores the data read in the translator workspace.
// If necessary, this routine may be called more than once.
func (w *Workspace) ReadData(filename string) {
	if w.w.w == nil {
		panic("MathProg method called on a deleted workspace")
	}
	fname := C.CString(filename)
	defer C.free(unsafe.Pointer(fname))
	if C.glp_mpl_read_data(w.w.w, fname) != 0 {
		panic("glp_mpl_read_data failed")
	}
}

// GenerateTo generates the model using its description stored in the translator workspace.
// This operation means generating all variables, constraints, and objectives, executing
// check and display statements, which precede the solve statement (if it is presented).
// The character string fname specifies the name of an output text file, to which output produced
// by display statements should be written.
func (w *Workspace) GenerateTo(filename string) {
	if w.w.w == nil {
		panic("MathProg method called on a deleted workspace")
	}

	fname := C.CString(filename)
	defer C.free(unsafe.Pointer(fname))
	if C.glp_mpl_generate(w.w.w, fname) != 0 {
		panic("glp_mpl_generate failed")
	}
}

// Generate generates the model using its description stored in the translator workspace.
// This operation means generating all variables, constraints, and objectives, executing
// check and display statements, which precede the solve statement (if it is presented).
// Output produced by display statements is sent to the terminal.
func (w *Workspace) Generate() {
	if w.w.w == nil {
		panic("MathProg method called on a deleted workspace")
	}

	if C.glp_mpl_generate(w.w.w, nil) != 0 {
		panic("glp_mpl_generate failed")
	}
}

// BuildProb obtains all necessary information from the translator workspace and stores
// it in the specified problem object p
func (w *Workspace) BuildProb(p *Prob) {
	if w.w.w == nil {
		panic("MathProg method called on a deleted workspace")
	}
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}

	C.glp_mpl_build_prob(w.w.w, p.p.p)
}

// PostSolve copies the solution from the specified problem object prob to
// the translator workspace and then executes all the remaining model statements,
// which follow the solve statement.
// The parameter sol specifies which solution should be copied from the problem object to the workspace.
func (w *Workspace) PostSolve(p *Prob, sol SolType) {
	if w.w.w == nil {
		panic("MathProg method called on a deleted workspace")
	}
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}

	if C.glp_mpl_postsolve(w.w.w, p.p.p, C.int(sol)) != 0 {
		panic("glp_mpl_postsolve failed")
	}
}
