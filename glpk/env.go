package glpk

/*
#cgo LDFLAGS: -lglpk
#include <glpk.h>
#include <stdlib.h>

int term_hook(char *s);

static int hook(void *info, const char *s)
{
	return term_hook((char *)s);
}

static void set_glp_term_hook()
{
	glp_term_hook(hook, NULL);
}

static void reset_glp_term_hook()
{
	glp_term_hook(NULL, NULL);
}
*/
import "C"

// TermHookFunc represents a terminal output callback function pointer
type TermHookFunc func(s string) bool

var (
	termHookCallback TermHookFunc
)

//export term_hook
func term_hook(s *C.char) C.int {
	if termHookCallback != nil {
		str := C.GoString(s)
		if ret := termHookCallback(str); ret {
			return C.int(1)
		}
		return C.int(0)
	} else {
		C.reset_glp_term_hook()
		return C.int(0)
	}
}

// SetTermHook installs the user-defined hook routine to intercept all terminal
// output performed by GLPK routines.
func SetTermHook(callback TermHookFunc) {
	// C.glp_term_hook(, nil)
	if callback != nil {
		termHookCallback = callback
		C.set_glp_term_hook()
	} else {
		termHookCallback = nil
		C.reset_glp_term_hook()
	}
}

// SetTermOut enables or disables terminal output performed by glpk routines.
func SetTermOut(flag bool) {
	if flag {
		C.glp_term_out(C.int(1))
	} else {
		C.glp_term_out(C.int(0))
	}
}
