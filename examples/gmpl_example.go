package main

import (
	"fmt"
	"log"

	"github.com/fastdeath/glpk/glpk"
)

func onOutput(s string) bool {
	// fmt.Print("/")
	fmt.Print(s)
	return true
}

func main() {
	glpk.SetTermOut(true)
	glpk.SetTermHook(onOutput)
	// glpk.SetTermHook(nil)

	lp := glpk.New()
	defer lp.Delete()

	tran := glpk.NewMPL()
	defer tran.Free()

	tran.ReadModel("sudoku.mod", true)
	tran.ReadData("sudoku.dat")
	tran.Generate()
	tran.BuildProb(lp)

	smcp := glpk.NewSmcp()
	smcp.SetMsgLev(glpk.MSG_ERR)
	lp.Simplex(smcp)

	iocp := glpk.NewIocp()
	iocp.SetMsgLev(glpk.MSG_ERR)
	iocp.SetPresolve(false)

	if err := lp.Intopt(iocp); err != nil {
		log.Fatalf("Mip error: %v", err)
	}

	tran.PostSolve(lp, glpk.MIP)
}
