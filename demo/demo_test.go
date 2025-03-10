package main

import (
	_ "embed"
	"testing"
)

func TestDemoDef(t *testing.T) {
	run("def.tem", def)
}

func TestDemoSemicolon(t *testing.T) {
	run("def_semicolon", semicolon)
}

// func TestDemoTmpl(t *testing.T) {
// 	run("template.tem", tmpl)
// }

// func TestDemoTmplSemicolon(t *testing.T) {
// 	run("template_semicolon.tem", tmplSemicolon)
// }
