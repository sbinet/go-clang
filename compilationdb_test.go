package clang_test

import (
	"testing"

	clang "github.com/sbinet/go-clang"
)

func TestCompilationDatabaseError(t *testing.T) {
	_, err := clang.NewCompilationDatabase("testdata-not-there")
	if err == nil {
		t.Fatalf("expected an error")
	}

	if err.(clang.CompilationDatabaseError) != clang.CompilationDatabase_CanNotLoadDatabase {
		t.Fatalf("expected %v", clang.CompilationDatabase_CanNotLoadDatabase)
	}
}

func TestCompilationDatabase(t *testing.T) {
	db, err := clang.NewCompilationDatabase("testdata")
	if err != nil {
		t.Fatalf("error loading compilation database: %v", err)
	}
	defer db.Dispose()

	table := []struct {
		directory string
		args      []string
	}{
		{
			directory: "/home/user/llvm/build",
			args: []string{
				"/usr/bin/clang++",
				"-Irelative",
				//FIXME: bug in clang ?
				//`-DSOMEDEF="With spaces, quotes and \-es.`,
				"-DSOMEDEF=With spaces, quotes and -es.",
				"-c",
				"-o",
				"file.o",
				"file.cc",
			},
		},
		{
			directory: "@TESTDIR@",
			args:      []string{"g++", "-c", "-DMYMACRO=a", "subdir/a.cpp"},
		},
	}

	cmds := db.GetAllCompileCommands()
	if cmds.GetSize() != len(table) {
		t.Errorf("expected #cmds=%d. got=%d", len(table), cmds.GetSize())
	}

	for i := 0; i < cmds.GetSize(); i++ {
		cmd := cmds.GetCommand(i)
		if cmd.GetDirectory() != table[i].directory {
			t.Errorf("expected dir=%q. got=%q", table[i].directory, cmd.GetDirectory())
		}

		nargs := cmd.GetNumArgs()
		if nargs != len(table[i].args) {
			t.Errorf("expected #args=%d. got=%d", len(table[i].args), nargs)
		}
		if nargs > len(table[i].args) {
			nargs = len(table[i].args)
		}
		for j := 0; j < nargs; j++ {
			arg := cmd.GetArg(j)
			if arg != table[i].args[j] {
				t.Errorf("expected arg[%d]=%q. got=%q", j, table[i].args[j], arg)
			}
		}
	}
}
