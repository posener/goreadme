// Package config holds configuration options for goreadme.
package config

type Config struct {
	// Override readme title. Default is package name.
	Title string `json:"title"`
	// ImportPath is used to override the import path. For example: github.com/user/project,
	// github.com/user/project/package or github.com/user/project/version.
	ImportPath string `json:"import_path"`
	// Consts will make constants documentation to be added to the README.
	// If Types is specified, constants for each type will also be added to the README.
	Consts bool `json:"consts"`
	// Vars will make exported variables documentation to be added to the README.
	// If Types is specified, exported variables for each type will also be added to the README.
	Vars bool `json:"vars"`
	// Functions will make functions documentation to be added to the README.
	Functions bool `json:"functions"`
	// Types will make types documentation to be added to the README.
	Types bool `json:"types"`
	// Factories will make functions returning a type to be added to the README, if Types is also specified.
	// Has no effect if Types is not specified.
	Factories bool `json:"factories"`
	// Methods will make the methods for a type to be added to the README, if Types is also specified.
	// Has no effect if Types is not specified.
	Methods bool `json:"methods"`
	// SkipExamples will omit the examples section from the README.
	SkipExamples bool `json:"skip_examples"`
	// SkipSubPackages will omit the sub packages section from the README.
	SkipSubPackages bool `json:"skip_sub_packages"`
	// NoDiffBlocks disables marking code blocks as diffs if they start with minus or plus signes.
	NoDiffBlocks bool `json:"no_diff_blocks"`
	// RecursiveSubPackages will retrieved subpackages information recursively.
	// If false, only one level of subpackages will be retrieved.
	RecursiveSubPackages bool `json:"recursive_sub_packages"`
	Badges               struct {
		TravisCI     bool `json:"travis_ci"`
		CodeCov      bool `json:"code_cov"`
		GolangCI     bool `json:"golang_ci"`
		GoDoc        bool `json:"go_doc"`
		GoReportCard bool `json:"go_report_card"`
	} `json:"badges"`
	Credit bool `json:"credit"`
}
