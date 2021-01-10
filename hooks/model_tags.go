package main

import (
	"fmt"
	"log"

	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/plugin/modelgen"
)

func addMongoDbTags(b *modelgen.ModelBuild) *modelgen.ModelBuild {
	for _, model := range b.Models {
		for _, field := range model.Fields {
			field.Tag += fmt.Sprintf(` bson:"%s"`, field.Name)
		}
	}
	return b
}

func main() {
	cfg, err := config.LoadConfigFromDefaultLocations()
	if err != nil {
		log.Fatal(err)
	}

	p := modelgen.Plugin{
		MutateHook: addMongoDbTags,
	}

	err = api.Generate(cfg,
		api.NoPlugins(),
		api.AddPlugin(&p),
	)

	if err != nil {
		log.Fatal(err)
	}
}
