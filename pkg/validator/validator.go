package validator

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/kaptinlin/jsonschema"
)

type jsonData struct {
	Properties interface{} `json:"properties"` // Adjust type if you know the structure of properties
}

type Validator struct {
	schemas        string
	compiledSchema *jsonschema.Schema
	compiler       *jsonschema.Compiler
}

func (v *Validator) LoadJsonSchemas(dirPath ...string) error {
	path := "json-schema"
	if len(dirPath) > 0 {
		path = dirPath[0]
	}

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			var jsonData jsonData
			err = json.Unmarshal(data, &jsonData)
			if err != nil {
				return err
			}

			byteData, err := json.Marshal(jsonData.Properties)
			if err != nil {
				return err
			}
			var strData string
			if len(v.schemas) > 0 {
				if v.schemas[len(v.schemas)-1] == '}' {
					v.schemas = v.schemas[:len(v.schemas)-1]
				}
				v.schemas += ","
				if byteData[0] == '{' {
					strData = string(byteData[1:])
				}
			} else {
				strData = string(byteData)
			}
			v.schemas += strData
		}
		return nil
	})
	if err != nil {
		return err
	}

	//fmt.Println(v.schemas)

	v.compiler = jsonschema.NewCompiler()
	v.compiledSchema, err = v.compiler.Compile([]byte(v.schemas))
	if err != nil {
		log.Fatalf("Failed to compile schema: %v", err)
	}

	return err
}

func (v *Validator) Validate(body []byte) error {
	var data map[string]interface{}

	// json validation
	if !json.Valid(body) {
		return fmt.Errorf("not valid JSON passed to Validator")
	}

	err := json.Unmarshal(body, &data)
	if err != nil {
		return err
	}

	fmt.Println(v.compiledSchema)
	fmt.Println(data)

	/*typeWebhook := ""
	typeWebhook = fmt.Sprintf("%s", data["typeWebhook"])
	t := `{"` + typeWebhook + `":` + string(body) + `}`
	fmt.Println(t)*/

	// json-schema validation
	result := v.compiledSchema.Validate(data)
	if !result.IsValid() {
		_, err := json.MarshalIndent(result.ToList(), "", "  ")
		return err
		//fmt.Println(string(details))
		//return fmt.Errorf("%s", details)
	}
	return nil
}
