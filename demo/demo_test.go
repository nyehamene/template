package main

import (
	_ "embed"
	"strings"
	"testing"
)

func TestDemoDef(t *testing.T) {
	str := getString("def.tem", def)
	if strings.Contains(str, "ERROR") {
		t.Error("parser failed unexpectedly")
	}
}

func TestDemoSemicolon(t *testing.T) {
	str := getString("def_semicolon", semicolon)
	if strings.Contains(str, "ERROR") {
		t.Error("parser failed unexpectedly")
	}
}

// func TestDemoTmpl(t *testing.T) {
// 	run("template.tem", tmpl)
// }

// func TestDemoTmplSemicolon(t *testing.T) {
// 	run("template_semicolon.tem", tmplSemicolon)
// }
