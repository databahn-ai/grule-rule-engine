package examples

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/databahn-ai/grule-rule-engine/ast"
	"github.com/databahn-ai/grule-rule-engine/builder"
	"github.com/databahn-ai/grule-rule-engine/engine"
	"github.com/databahn-ai/grule-rule-engine/pkg"
)

const (
	bracketNotationRule = `
	rule TestRule "" {
		when
			Event["accountname.0"].ToLower() == "new_user" || IsNil(Event["accountname.0"])
		then
			R.Result = "ok";
			Retract("TestRule");
	}
	`
)

func TestBracketNotationMissingField(t *testing.T) {
	result := &ObjectResult{
		Result: "NoResult",
	}

	// build rules
	lib := ast.NewKnowledgeLibrary()
	rb := builder.NewRuleBuilder(lib)
	err := rb.BuildRuleFromResource("Test", "0.0.1", pkg.NewBytesResource([]byte(bracketNotationRule)))
	assert.NoError(t, err)

	// add JSON fact without the "accountname.0" field
	json := []byte(`{"blabla":"bla","name":{"first":"john"}}`)
	kb, err := lib.NewKnowledgeBaseInstance("Test", "0.0.1")
	assert.NoError(t, err)
	dcx := ast.NewDataContext()

	err = dcx.Add("R", result)
	err = dcx.AddJSON("Event", json)
	if err != nil {
		fmt.Println(err.Error())
	}

	// This should not error even though "accountname.0" field is missing
	err = engine.NewGruleEngine().Execute(dcx, kb)
	assert.NoError(t, err)

	if result.Result != "ok" {
		t.Errorf("Expected result to be ok, got %s", result.Result)
	}
}

const (
	bracketNotationRuleWithField = `
	rule TestRule "" {
		when
			Event["accountname.0"].ToLower() == "new_user"
		then
			R.Result = "ok";
			Retract("TestRule");
	}
	`
)

func TestBracketNotationWithField(t *testing.T) {
	result := &ObjectResult{
		Result: "NoResult",
	}

	// build rules
	lib := ast.NewKnowledgeLibrary()
	rb := builder.NewRuleBuilder(lib)
	err := rb.BuildRuleFromResource("Test", "0.0.1", pkg.NewBytesResource([]byte(bracketNotationRuleWithField)))
	assert.NoError(t, err)

	// add JSON fact with the "accountname.0" field
	json := []byte(`{"accountname.0":"NEW_USER","blabla":"bla"}`)
	kb, err := lib.NewKnowledgeBaseInstance("Test", "0.0.1")
	assert.NoError(t, err)
	dcx := ast.NewDataContext()

	err = dcx.Add("R", result)
	err = dcx.AddJSON("Event", json)
	if err != nil {
		fmt.Println(err.Error())
	}

	// This should work and match the rule
	err = engine.NewGruleEngine().Execute(dcx, kb)
	assert.NoError(t, err)

	if result.Result != "ok" {
		t.Errorf("Expected result to be ok, got %s", result.Result)
	}
}
