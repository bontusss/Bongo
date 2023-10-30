package template

import (
	"os"
	"testing"

	"github.com/CloudyKit/jet/v6"
)

var views = jet.NewSet(jet.NewOSFileSystemLoader("testData/views"), jet.InDevelopmentMode())

var testEngine = Template {
	Engine: "",
	RootPath: "",
	JetViews: views,
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}