package zip_archive

import (
	"fmt"
	"path"
	"strings"

	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
)

var (
	pctx = blueprint.NewPackageContext("github.com/burbokop/design-practice-1/build/modules/zip_archive")

	zip = pctx.StaticRule("zip", blueprint.RuleParams{
		Command:     "zip $resultPath $srcPaths",
		Description: "make zip archive $resultPath",
	}, "resultPath", "srcPaths")
)

type zipModuleType struct {
	blueprint.SimpleName

	properties struct {
		Files    []string
		Excluded []string
	}
}

func (gtb *zipModuleType) GenerateBuildActions(ctx blueprint.ModuleContext) {
	name := ctx.ModuleName()
	config := bood.ExtractConfig(ctx)
	config.Debug.Printf("Adding build actions for go binary module '%s'", name)

	outputPath := path.Join(config.BaseOutputDir, "zip", name)

	var inputs []string
	inputErors := false
	for _, src := range gtb.properties.Files {
		matches, err := ctx.GlobWithDeps(src, gtb.properties.Excluded)
		if err == nil {
			inputs = append(inputs, matches...)
		} else {
			ctx.PropertyErrorf("srcs", "Cannot resolve files that match pattern %s", src)
			inputErors = true
		}
	}
	if inputErors {
		return
	}

	ctx.Build(pctx, blueprint.BuildParams{
		Description: fmt.Sprintf("Zip %s archive", name),
		Rule:        zip,
		Outputs:     []string{outputPath},
		Implicits:   inputs,
		Args: map[string]string{
			"resultPath": outputPath,
			"srcPaths":   strings.Join(inputs, " "),
		},
	})
}

func SimpleBinFactory() (blueprint.Module, []interface{}) {
	mType := &zipModuleType{}
	return mType, []interface{}{&mType.SimpleName.Properties, &mType.properties}
}
