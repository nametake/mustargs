package mustargs_test

import (
	"path/filepath"
	"testing"

	"github.com/gostaticanalysis/testutil"
	"github.com/nametake/mustargs"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	tests := []struct {
		configFile string
		patterns   []string
	}{
		{
			configFile: "testdata/src/primitive/config.yaml",
			patterns:   []string{"primitive"},
		},
		{
			configFile: "testdata/src/argindex/config.yaml",
			patterns:   []string{"argindex"},
		},
		{
			configFile: "testdata/src/pkgtype/config.yaml",
			patterns:   []string{"pkgtype"},
		},
	}

	testdata := testutil.WithModules(t, analysistest.TestData(), nil)
	for _, tt := range tests {
		configFile := tt.configFile
		defaultPath, err := filepath.Abs(configFile)
		if err != nil {
			t.Error(err)
			return
		}
		if err := mustargs.Analyzer.Flags.Set("config", defaultPath); err != nil {
			t.Error(err)
			return
		}
		analysistest.Run(t, testdata, mustargs.Analyzer, tt.patterns...)
	}
}
