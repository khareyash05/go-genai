package genai

import (
	"testing"

	"errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFormatMap_ReplacesPlaceholders_321 tests that the formatMap function correctly replaces placeholders in a template string with values from a map.
func TestFormatMap_ReplacesPlaceholders_321(t *testing.T) {
	template := "Hello, {name}! Welcome to {place}."
	variables := map[string]any{
		"name":  "Alice",
		"place": "Wonderland",
	}
	result, err := formatMap(template, variables)
	require.NoError(t, err)
	assert.Equal(t, "Hello, Alice! Welcome to Wonderland.", result)
}

// TestSetValueByPath_SetsValue_456 tests that the setValueByPath function correctly sets a value in a nested map structure.
func TestSetValueByPath_SetsValue_456(t *testing.T) {
	data := make(map[string]any)
	keys := []string{"level1", "level2", "level3"}
	value := "testValue"
	setValueByPath(data, keys, value)
	require.NotNil(t, data["level1"])
	nestedMap, ok := data["level1"].(map[string]any)
	require.True(t, ok)
	require.NotNil(t, nestedMap["level2"])
	nestedMap, ok = nestedMap["level2"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, value, nestedMap["level3"])
}

// TestApplyConverterToSlice_AppliesFunction_654 tests that the applyConverterToSlice function correctly applies a converter function to each element in a slice.
func TestApplyConverterToSlice_AppliesFunction_654(t *testing.T) {
	mockConverter := func(ac *apiClient, input map[string]any, _ map[string]any) (map[string]any, error) {
		input["converted"] = true
		return input, nil
	}
	inputs := []any{
		map[string]any{"key1": "value1"},
		map[string]any{"key2": "value2"},
	}
	outputs, err := applyConverterToSlice(nil, inputs, mockConverter)
	require.NoError(t, err)
	require.Len(t, outputs, 2)
	assert.Equal(t, true, outputs[0]["converted"])
	assert.Equal(t, true, outputs[1]["converted"])
}

// TestApplyItemTransformerToSlice_ErrorHandling_123 tests that applyItemTransformerToSlice correctly handles errors returned by the transformer function.
func TestApplyItemTransformerToSlice_ErrorHandling_123(t *testing.T) {
	mockTransformer := func(ac *apiClient, input string) (string, error) {
		if input == "error" {
			return "", errors.New("mock error")
		}
		return input + "_transformed", nil
	}
	inputs := []string{"valid", "error", "anotherValid"}
	outputs, err := applyItemTransformerToSlice(nil, inputs, mockTransformer)
	require.Error(t, err)
	assert.Nil(t, outputs)
}

// TestGetValueByPath_RetrievesValue_789 tests that the getValueByPath function correctly retrieves a value from a nested map structure.
func TestGetValueByPath_RetrievesValue_789(t *testing.T) {
	data := map[string]any{
		"level1": map[string]any{
			"level2": map[string]any{
				"level3": "testValue",
			},
		},
	}
	keys := []string{"level1", "level2", "level3"}
	value := getValueByPath(data, keys)
	assert.Equal(t, "testValue", value)
}

// TestDeepMarshal_MarshalsAndUnmarshals_987 tests that the deepMarshal function correctly marshals and unmarshals a map.
func TestDeepMarshal_MarshalsAndUnmarshals_987(t *testing.T) {
	input := map[string]any{"key": "value"}
	var output map[string]any
	err := deepMarshal(input, &output)
	require.NoError(t, err)
	assert.Equal(t, input, output)
}

// TestPtr_ReturnsPointer_123 tests that the Ptr function correctly returns a pointer to the provided value.
func TestPtr_ReturnsPointer_123(t *testing.T) {
	value := 42
	ptr := Ptr(value)
	require.NotNil(t, ptr)
	assert.Equal(t, value, *ptr)
}

// TestApplyItemTransformerToSlice_SuccessfulTransformation_456 tests that applyItemTransformerToSlice correctly transforms all inputs when no errors occur.
func TestApplyItemTransformerToSlice_SuccessfulTransformation_456(t *testing.T) {
	mockTransformer := func(ac *apiClient, input string) (string, error) {
		return input + "_transformed", nil
	}
	inputs := []string{"input1", "input2", "input3"}
	outputs, err := applyItemTransformerToSlice(nil, inputs, mockTransformer)
	require.NoError(t, err)
	require.Len(t, outputs, 3)
	assert.Equal(t, "input1_transformed", outputs[0])
	assert.Equal(t, "input2_transformed", outputs[1])
	assert.Equal(t, "input3_transformed", outputs[2])
}
