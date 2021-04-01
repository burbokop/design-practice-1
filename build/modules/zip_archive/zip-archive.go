package zip_archive

import (
	"fmt"
	"path"

	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
)

var (
	pctx = blueprint.NewPackageContext("github.com/burbokop/design-practice-1/build/modules/zip_archive")

	zip = pctx.StaticRule("binaryBuild", blueprint.RuleParams{
		Command:     "cd $workDir && go build -o $outputPath $pkg",
		Description: "build go command $pkg",
	}, "workDir", "outputPath", "pkg")
)

type goTestedBinaryModuleType struct {
	blueprint.SimpleName

	properties struct {
		Pkg         string
		Srcs        []string
		SrcsExclude []string

		TestPkg         string
		TestSrcs        []string
		TestSrcsExclude []string

		VendorFirst bool

		Deps []string
	}
}

func (gtb *goTestedBinaryModuleType) DynamicDependencies(blueprint.DynamicDependerModuleContext) []string {
	return gtb.properties.Deps
}

func (gtb *goTestedBinaryModuleType) GenerateBuildActions(ctx blueprint.ModuleContext) {
	name := ctx.ModuleName()
	config := bood.ExtractConfig(ctx)
	config.Debug.Printf("Adding build actions for go binary module '%s'", name)

	outputPath := path.Join(config.BaseOutputDir, "bin", name)
	testOutputPath := path.Join(config.BaseOutputDir, "reports", name)

	var inputs []string
	inputErors := false
	for _, src := range gtb.properties.Srcs {
		matches, err := ctx.GlobWithDeps(src, append(gtb.properties.SrcsExclude, gtb.properties.TestSrcs...))
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

	if gtb.properties.VendorFirst {
		vendorDirPath := path.Join(ctx.ModuleDir(), "vendor")
		ctx.Build(pctx, blueprint.BuildParams{
			Description: fmt.Sprintf("Vendor dependencies of %s", name),
			Rule:        goVendor,
			Outputs:     []string{vendorDirPath},
			Implicits:   []string{path.Join(ctx.ModuleDir(), "go.mod")},
			Optional:    true,
			Args: map[string]string{
				"workDir": ctx.ModuleDir(),
				"name":    name,
			},
		})
		inputs = append(inputs, vendorDirPath)
	}

	ctx.Build(pctx, blueprint.BuildParams{
		Description: fmt.Sprintf("Build %s as Go binary", name),
		Rule:        goBuild,
		Outputs:     []string{outputPath},
		Implicits:   inputs,
		Args: map[string]string{
			"outputPath": outputPath,
			"workDir":    ctx.ModuleDir(),
			"pkg":        gtb.properties.Pkg,
		},
	})

	for _, testSrc := range gtb.properties.TestSrcs {
		if matches, err := ctx.GlobWithDeps(testSrc, gtb.properties.TestSrcsExclude); err == nil {
			inputs = append(inputs, matches...)
		} else {
			ctx.PropertyErrorf("testSrcs", "Cannot resolve files that match pattern %s", testSrc)
			inputErors = true
		}
	}
	if inputErors {
		return
	}
	ctx.Build(pctx, blueprint.BuildParams{
		Description: "Test my module",
		Rule:        goTest,
		Outputs:     []string{testOutputPath},
		Implicits:   inputs,
		Args: map[string]string{
			"outputPath": testOutputPath,
			"workDir":    ctx.ModuleDir(),
			"testPkg":    gtb.properties.TestPkg,
		},
	})

}

func SimpleBinFactory() (blueprint.Module, []interface{}) {
	mType := &goTestedBinaryModuleType{}
	return mType, []interface{}{&mType.SimpleName.Properties, &mType.properties}
}
