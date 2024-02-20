package printer

import (
	"io"
	"os"

	"github.com/k0kubun/pp/v3"
	"github.com/mattn/go-isatty"
)

func initPP() {
	out := os.Stdout
	pp.SetDefaultOutput(out)

	if !isatty.IsTerminal(out.Fd()) {
		prettyPrinter := pp.New()
		prettyPrinter.SetColoringEnabled(false)
		prettyPrinter.SetExportedOnly(true)
	}
}

func Println(args ...any) (n int, err error) {
	return pp.Println(args...)
}

func Printf(format string, args ...any) (n int, err error) {
	return pp.Printf(format, args...)
}

func Print(args ...any) (n int, err error) {
	return pp.Print(args...)
}

func Sprint(args ...any) string {
	return pp.Sprint(args...)
}

func Sprintf(format string, args ...any) string {
	return pp.Sprintf(format, args...)
}

func Sprintln(args ...any) string {
	return pp.Sprintln(args...)
}

func Fprint(w io.Writer, args ...any) (n int, err error) {
	return pp.Fprint(w, args...)
}

func Fprintf(w io.Writer, format string, args ...any) (n int, err error) {
	return pp.Fprintf(w, format, args...)
}

func Fprintln(w io.Writer, args ...any) (n int, err error) {
	return pp.Fprintln(w, args...)
}

func Errorf(format string, args ...any) error {
	return pp.Errorf(format, args...)
}

func Fatal(args ...any) {
	pp.Fatal(args...)
}
