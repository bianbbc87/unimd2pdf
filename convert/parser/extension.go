package parser

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

// GoldmarkExtension wraps a goldmark.Extender with a name for config toggling.
type GoldmarkExtension struct {
	ExtName string
	Ext     goldmark.Extender
}

func (e *GoldmarkExtension) Name() string               { return e.ExtName }
func (e *GoldmarkExtension) Extender() goldmark.Extender { return e.Ext }

// Built-in extension constructors

func NewFootnote() *GoldmarkExtension {
	return &GoldmarkExtension{ExtName: "footnote", Ext: extension.NewFootnote()}
}

func NewDefinitionList() *GoldmarkExtension {
	return &GoldmarkExtension{ExtName: "definitionlist", Ext: extension.DefinitionList}
}

func NewTypographer() *GoldmarkExtension {
	return &GoldmarkExtension{ExtName: "typographer", Ext: extension.NewTypographer()}
}

func NewCJK() *GoldmarkExtension {
	return &GoldmarkExtension{ExtName: "cjk", Ext: extension.NewCJK()}
}
