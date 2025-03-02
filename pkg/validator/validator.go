package validator

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

type Validator struct {
	compiledSchema *jsonschema.Schema
	compiler       *jsonschema.Compiler
}

func (v *Validator) LoadJsonSchemas(dirPath ...string) error {
	path := "json-schema"
	if len(dirPath) > 0 {
		path = dirPath[0]
	}

	// properties from all .json files will be merged at this var
	var allProperties map[string]interface{}

	// all json validation schema files must contains this fields:
	schema := map[string]interface{}{
		"$id":        "schemas",
		"$schema":    "https://json-schema.org/draft/2020-12/schema",
		"type":       "object",
		"properties": allProperties,
	}

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.Compare(info.Name(), "gererated-final-schemas.json") != 0 {
			// Read each JSON file
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			err = json.Unmarshal(data, &schema)
			if err != nil {
				return fmt.Errorf("not valid JSON passed to LoadJsonSchemas" + path)
			}

			// Check if the "properties" field exists in the file
			if properties, ok := schema["properties"].(map[string]interface{}); ok {
				if allProperties == nil {
					allProperties = make(map[string]interface{})
				}
				// Merge the properties
				for key, value := range properties {
					allProperties[key] = value
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Write the merged schema into a final file
	finalSchema := map[string]interface{}{
		"$id":        "schemas",
		"$schema":    "https://json-schema.org/draft/2020-12/schema",
		"type":       "object",
		"properties": allProperties,
	}

	finalSchemaData, err := json.MarshalIndent(finalSchema, "", "  ")
	if err != nil {
		return err
	}

	finalFilePath := filepath.Join(path, "gen", "gererated-final-schemas.json")
	err = os.MkdirAll(filepath.Dir(finalFilePath), 0644)
	if err != nil {
		return err
	}

	err = os.WriteFile(finalFilePath, finalSchemaData, 0644)
	if err != nil {
		return err
	}

	v.compiler = jsonschema.NewCompiler()
	v.compiledSchema, err = v.compiler.Compile(finalFilePath)
	if err != nil {
		return err
	}

	return nil
}

func (v *Validator) Validate(body []byte) (map[string]interface{}, error) {
	var data map[string]interface{}           // return this JSON
	var dataToValidate map[string]interface{} // adjuct JSON for validation

	if !json.Valid(body) {
		return nil, fmt.Errorf("not valid JSON passed to Validator")
	}

	err := json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	// every webhook must contain typeWebhook field
	if data["typeWebhook"] == nil {
		return nil, fmt.Errorf("no typeWebhook field")
	}
	typeWebhook := fmt.Sprintf("%s", data["typeWebhook"])

	// adjuct JSON for validation schemas
	jsonToValidateStr := "{\r\n\"" + typeWebhook + "\":\r\n" + string(body) + "\r\n}"
	err = json.Unmarshal([]byte(jsonToValidateStr), &dataToValidate)
	if err != nil {
		return nil, err
	}

	// json-schema validation
	result := v.compiledSchema.Validate(dataToValidate)
	if result != nil {
		log.Print(result.Error())
		return nil, result
	}
	return data, nil
}
