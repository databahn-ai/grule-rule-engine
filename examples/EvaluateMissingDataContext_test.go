package examples

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/databahn-ai/grule-rule-engine/ast"
	"github.com/databahn-ai/grule-rule-engine/builder"
	"github.com/databahn-ai/grule-rule-engine/engine"
	"github.com/databahn-ai/grule-rule-engine/pkg"
)

const (
	inputRule = `
	rule TestRule "" {
		when
			R.Result == 'NoResult' &&
			inputs.i_am_missing == 'abc' &&
                        inputs.name.first == 'john'
		then
			R.Result = "ok";
	}
	`
)

func TestDataContextMissingFact(t *testing.T) {

	oresult := &ObjectResult{
		Result: "NoResult",
	}

	// build rules
	lib := ast.NewKnowledgeLibrary()
	rb := builder.NewRuleBuilder(lib)
	err := rb.BuildRuleFromResource("Test", "0.0.1", pkg.NewBytesResource([]byte(inputRule)))

	// 	add JSON fact
	json := []byte(`{"blabla":"bla","name":{"first":"john","last":"doe"}}`)
	kb, err := lib.NewKnowledgeBaseInstance("Test", "0.0.1")
	assert.NoError(t, err)
	dcx := ast.NewDataContext()

	err = dcx.Add("R", oresult)
	err = dcx.AddJSON("inputs", json)
	if err != nil {
		fmt.Println(err.Error())
	}

	// results in panic
	engine.NewGruleEngine().Execute(dcx, kb)

}

const (
	matchingRule = `
	rule TestRule "" {
		when
			inputs.i_am_missing == 'abc' || R.Result == 'NoResult' || inputs.name.first == 'john'
		then
			R.Result = "ok";
			Retract("TestRule");
	}
	`
)

func TestDataContextWithFactMissingFieldsAndMatching(t *testing.T) {
	result := &ObjectResult{
		Result: "NoResult",
	}

	// build rules
	lib := ast.NewKnowledgeLibrary()
	rb := builder.NewRuleBuilder(lib)
	err := rb.BuildRuleFromResource("Test", "0.0.1", pkg.NewBytesResource([]byte(matchingRule)))

	// 	add JSON fact
	json := []byte(`{"blabla":"bla","name":{"first":"john","last":"doe"}}`)
	kb, err := lib.NewKnowledgeBaseInstance("Test", "0.0.1")
	assert.NoError(t, err)
	dcx := ast.NewDataContext()

	err = dcx.Add("R", result)
	err = dcx.AddJSON("inputs", json)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = engine.NewGruleEngine().Execute(dcx, kb)
	assert.NoError(t, err)

	if result.Result != "ok" {
		t.Errorf("Expected result to be ok, got %s", result.Result)
	}

}

const (
	matchingRuleWithFunctionCall = `
	rule TestRule "" {
		when
			R.IsMatching("matching", inputs.missing) && inputs.i_am_missing != 'abc' && R.Result == 'NoResult' && inputs.name.first == 'john'
		then
			R.PassedWith("ok", inputs.description_missing);
			Retract("TestRule");
	}
	`
)

type ObjectResultWithMethod struct {
	Result      string
	Description string
}

func (o *ObjectResultWithMethod) PassedWith(result string, description string) {
	o.Result = result
	o.Description = description
}

func (o *ObjectResultWithMethod) IsMatching(query string, anything string) bool {
	return query == "matching" || anything == "any"
}

func TestDataContextWithFactMissingFieldsAndMatchingWithFunctionCalls(t *testing.T) {
	result := &ObjectResultWithMethod{
		Result: "NoResult",
	}

	// build rules
	lib := ast.NewKnowledgeLibrary()
	rb := builder.NewRuleBuilder(lib)
	err := rb.BuildRuleFromResource("Test", "0.0.1", pkg.NewBytesResource([]byte(matchingRuleWithFunctionCall)))

	// 	add JSON fact
	json := []byte(`{"blabla":"bla","name":{"first":"john","last":"doe"}}`)
	kb, err := lib.NewKnowledgeBaseInstance("Test", "0.0.1")
	assert.NoError(t, err)
	dcx := ast.NewDataContext()

	err = dcx.Add("R", result)
	err = dcx.AddJSON("inputs", json)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = engine.NewGruleEngine().Execute(dcx, kb)
	assert.NoError(t, err)

	if result.Result != "ok" {
		t.Errorf("Expected result to be ok, got %s", result.Result)
	}

	if result.Description != "" {
		t.Errorf("Expected Description to be empty string, got %s", result.Description)
	}

}

const (
	isZeroIsNilRule = `
	rule TestRule "" {
		when
			IsZero(inputs.i_am_missing) && IsNil(inputs.another_missing)
		then
			R.Result = "ok";
			Retract("TestRule");
	}
	`
)

func TestDataContextWithFactMissingFieldsWithIsZeroIsNil(t *testing.T) {
	result := &ObjectResult{
		Result: "NoResult",
	}

	// build rules
	lib := ast.NewKnowledgeLibrary()
	rb := builder.NewRuleBuilder(lib)
	err := rb.BuildRuleFromResource("Test", "0.0.1", pkg.NewBytesResource([]byte(isZeroIsNilRule)))

	// 	add JSON fact
	json := []byte(`{"blabla":"bla","name":{"first":"john","last":"doe"}}`)
	kb, err := lib.NewKnowledgeBaseInstance("Test", "0.0.1")
	assert.NoError(t, err)
	dcx := ast.NewDataContext()

	err = dcx.Add("R", result)
	err = dcx.AddJSON("inputs", json)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = engine.NewGruleEngine().Execute(dcx, kb)
	assert.NoError(t, err)

	if result.Result != "ok" {
		t.Errorf("Expected result to be ok, got %s", result.Result)
	}

}

const (
	functionCallRule = `
	rule TestRule "" {
		when
			inputs.i_am_missing.ToLower().HasPrefix("miss") || IsNil(inputs.another_missing)
		then
			R.Result = "ok";
			Retract("TestRule");
	}
	`
)

func TestDataContextWithFactMissingFieldsWithFunctionCall(t *testing.T) {
	result := &ObjectResult{
		Result: "NoResult",
	}

	// build rules
	lib := ast.NewKnowledgeLibrary()
	rb := builder.NewRuleBuilder(lib)
	err := rb.BuildRuleFromResource("Test", "0.0.1", pkg.NewBytesResource([]byte(functionCallRule)))

	// 	add JSON fact
	json := []byte(`{"blabla":"bla","name":{"first":"john","last":"doe"}}`)
	kb, err := lib.NewKnowledgeBaseInstance("Test", "0.0.1")
	assert.NoError(t, err)
	dcx := ast.NewDataContext()

	err = dcx.Add("R", result)
	err = dcx.AddJSON("inputs", json)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = engine.NewGruleEngine().Execute(dcx, kb)
	assert.NoError(t, err)

	if result.Result != "ok" {
		t.Errorf("Expected result to be ok, got %s", result.Result)
	}

}
