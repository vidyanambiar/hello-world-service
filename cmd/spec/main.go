// Copyright Red Hat

package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
	"github.com/ghodss/yaml"
	"github.com/identitatem/idp-configs-api/pkg/errors"
	"github.com/identitatem/idp-configs-api/pkg/models"
)

// Used to generate openapi yaml file for components.
func main() {
	components := openapi3.NewComponents()
	components.Schemas = make(map[string]*openapi3.SchemaRef)

	authRealm, _, err := openapi3gen.NewSchemaRefForValue(&models.AuthRealm{})
	if err != nil {
		panic(err)
	}
	components.Schemas["v1.AuthRealm"] = authRealm	

	authRealmBody, _, err := openapi3gen.NewSchemaRefForValue(&models.AuthRealmBody{})
	if err != nil {
		panic(err)
	}
	components.Schemas["v1.AuthRealmBody"] = authRealmBody		

	internalServerError, _, err := openapi3gen.NewSchemaRefForValue(&errors.InternalServerError{})
	if err != nil {
		panic(err)
	}
	components.Schemas["v1.InternalServerError"] = internalServerError

	badRequest, _, err := openapi3gen.NewSchemaRefForValue(&errors.BadRequest{})
	if err != nil {
		panic(err)
	}
	components.Schemas["v1.BadRequest"] = badRequest

	notFound, _, err := openapi3gen.NewSchemaRefForValue(&errors.NotFound{})
	if err != nil {
		panic(err)
	}
	components.Schemas["v1.NotFound"] = notFound

	conflict, _, err := openapi3gen.NewSchemaRefForValue(&errors.Conflict{})
	if err != nil {
		panic(err)
	}
	components.Schemas["v1.Conflict"] = conflict	

	forbidden, _, err := openapi3gen.NewSchemaRefForValue(&errors.Forbidden{})
	if err != nil {
		panic(err)
	}
	components.Schemas["v1.Forbidden"] = forbidden	

	type Swagger struct {
		Components openapi3.Components `json:"components,omitempty" yaml:"components,omitempty"`
	}

	swagger := Swagger{}
	swagger.Components = components

	b := &bytes.Buffer{}
	err = json.NewEncoder(b).Encode(swagger)
	checkErr(err)

	schema, err := yaml.JSONToYAML(b.Bytes())
	checkErr(err)

	paths, err := ioutil.ReadFile("./cmd/spec/path.yaml")
	checkErr(err)

	b = &bytes.Buffer{}
	b.Write(schema)
	b.Write(paths)

	doc, err := openapi3.NewLoader().LoadFromData(b.Bytes())
	checkErr(err)

	jsonB, err := doc.MarshalJSON()
	checkErr(err)
	err = ioutil.WriteFile("./cmd/spec/openapi.json", jsonB, 0666)
	checkErr(err)
	err = ioutil.WriteFile("./cmd/spec/openapi.yaml", b.Bytes(), 0666)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
