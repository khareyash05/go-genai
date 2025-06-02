package genai

import (
	"testing"

	"errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFormatMap_ReplacesPlaceholders_004 tests the formatMap function to ensure it replaces placeholders in a template string with values from a map.
func TestFormatMap_ReplacesPlaceholders_004(t *testing.T) {
	template := "Hello, {name}! Welcome to {place}."
	variables := map[string]any{
		"name":  "Alice",
		"place": "Wonderland",
	}
	result, err := formatMap(template, variables)
	require.NoError(t, err)
	assert.Equal(t, "Hello, Alice! Welcome to Wonderland.", result)
}

// TestSetValueByPath_SetsValue_002 tests the setValueByPath function to ensure it correctly sets a value in a nested map structure.
func TestSetValueByPath_SetsValue_002(t *testing.T) {
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

// TestApplyItemTransformerToSlice_Success_010 tests applyItemTransformerToSlice when the transformer function succeeds for all elements.
func TestApplyItemTransformerToSlice_Success_010(t *testing.T) {
	mockTransformer := func(ac *apiClient, input string) (string, error) {
		return input + "_transformed", nil
	}
	inputs := []string{"item1", "item2"}
	ac := &apiClient{}
	outputs, err := applyItemTransformerToSlice(ac, inputs, mockTransformer)
	require.NoError(t, err)
	require.Len(t, outputs, 2)
	assert.Equal(t, "item1_transformed", outputs[0])
	assert.Equal(t, "item2_transformed", outputs[1])
}

// TestApplyConverterToSlice_Success_008 tests applyConverterToSlice when the converter function succeeds for all elements.
func TestApplyConverterToSlice_Success_008(t *testing.T) {
	mockConverter := func(ac *apiClient, input map[string]any, _ map[string]any) (map[string]any, error) {
		return map[string]any{"converted": true}, nil
	}
	inputs := []any{
		map[string]any{"key1": "value1"},
		map[string]any{"key2": "value2"},
	}
	ac := &apiClient{}
	outputs, err := applyConverterToSlice(ac, inputs, mockConverter)
	require.NoError(t, err)
	require.Len(t, outputs, 2)
	assert.Equal(t, map[string]any{"converted": true}, outputs[0])
	assert.Equal(t, map[string]any{"converted": true}, outputs[1])
}

// TestGetValueByPath_RetrievesValue_003 tests the getValueByPath function to ensure it retrieves a value from a nested map structure.
func TestGetValueByPath_RetrievesValue_003(t *testing.T) {
	data := map[string]any{
		"level1": map[string]any{
			"level2": map[string]any{
				"level3": "testValue",
			},
		},
	}
	keys := []string{"level1", "level2", "level3"}
	result := getValueByPath(data, keys)
	assert.Equal(t, "testValue", result)
}

// TestDeepMarshal_MarshalsAndUnmarshals_006 tests the deepMarshal function to ensure it correctly marshals and unmarshals a map.
func TestDeepMarshal_MarshalsAndUnmarshals_006(t *testing.T) {
	input := map[string]any{"key": "value"}
	var output map[string]any
	err := deepMarshal(input, &output)
	require.NoError(t, err)
	assert.Equal(t, input, output)
}

// TestPtr_ReturnsPointer_001 tests the Ptr function to ensure it returns a pointer to the given argument.
func TestPtr_ReturnsPointer_001(t *testing.T) {
	value := 42
	ptr := Ptr(value)
	require.NotNil(t, ptr)
	assert.Equal(t, value, *ptr)
}

// TestFormatMap_UnsupportedTypeError_005 tests the formatMap function to ensure it returns an error when a placeholder value is not a string.
func TestFormatMap_UnsupportedTypeError_005(t *testing.T) {
	template := "Hello, {name}!"
	variables := map[string]any{
		"name": 123, // Invalid type
	}
	result, err := formatMap(template, variables)
	assert.Error(t, err)
	assert.Equal(t, "", result)
}

// TestDeepMarshal_MarshalError_007 tests the deepMarshal function to ensure it returns an error when the input map cannot be marshaled.
func TestDeepMarshal_MarshalError_007(t *testing.T) {
	input := map[string]any{"key": make(chan int)} // Invalid type
	var output map[string]any
	err := deepMarshal(input, &output)
	assert.Error(t, err)
}

// TestApplyConverterToSlice_Error_009 tests applyConverterToSlice when the converter function fails for one element.
func TestApplyConverterToSlice_Error_009(t *testing.T) {
	mockConverter := func(ac *apiClient, input map[string]any, _ map[string]any) (map[string]any, error) {
		if input["key"] == "error" {
			return nil, errors.New("conversion error")
		}
		return map[string]any{"converted": true}, nil
	}
	inputs := []any{
		map[string]any{"key": "value1"},
		map[string]any{"key": "error"},
	}
	ac := &apiClient{}
	outputs, err := applyConverterToSlice(ac, inputs, mockConverter)
	require.Error(t, err)
	assert.Nil(t, outputs)
}

// TestApplyItemTransformerToSlice_Error_011 tests applyItemTransformerToSlice when the transformer function fails for one element.
func TestApplyItemTransformerToSlice_Error_011(t *testing.T) {
	mockTransformer := func(ac *apiClient, input string) (string, error) {
		if input == "error" {
			return "", errors.New("transformation error")
		}
		return input + "_transformed", nil
	}
	inputs := []string{"item1", "error"}
	ac := &apiClient{}
	outputs, err := applyItemTransformerToSlice(ac, inputs, mockTransformer)
	require.Error(t, err)
	assert.Nil(t, outputs)
}
