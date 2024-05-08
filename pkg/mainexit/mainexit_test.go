package mainexit

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestMainExit(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), Analyzer, "./...")
}
