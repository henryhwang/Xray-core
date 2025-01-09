package main

import (
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func main() {
	protogen.Options{}.Run(func(gen *protogen.Plugin) error {
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			generateFile(gen, f)
		}
		return nil
	})
}

func generateFile(gen *protogen.Plugin, file *protogen.File) {
	filename := file.GeneratedFilenamePrefix + ".schema.json"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)

	g.P(`{`)
	g.P(`"definitions": {`)
	for _, message := range file.Messages {
		generateMessageSchema(g, message)
	}
	g.P(`}`)
	g.P(`}`)
}

func generateMessageSchema(g *protogen.GeneratedFile, message *protogen.Message) {
	g.P(fmt.Sprintf(`"%s": {`, message.GoIdent.GoName))
	g.P(`"type": "object",`)
	g.P(`"properties": {`)
	for _, field := range message.Fields {
		generateFieldSchema(g, field)
	}
	g.P(`}`)
	g.P(`},`)
}

func generateFieldSchema(g *protogen.GeneratedFile, field *protogen.Field) {
	jsonName := field.Desc.JSONName()
	g.P(fmt.Sprintf(`"%s": {`, jsonName))
	
	if field.Desc.IsList() {
		g.P(`"type": "array",`)
		g.P(`"items": {`)
		generateTypeSchema(g, field.Desc.Kind(), field.Message)
		g.P(`}`)
	} else {
		generateTypeSchema(g, field.Desc.Kind(), field.Message)
	}
	g.P(`},`)
}

func generateTypeSchema(g *protogen.GeneratedFile, kind protoreflect.Kind, message *protogen.Message) {
	switch kind {
	case protoreflect.BoolKind:
		g.P(`"type": "boolean"`)
	case protoreflect.Int32Kind, protoreflect.Int64Kind,
		protoreflect.Uint32Kind, protoreflect.Uint64Kind,
		protoreflect.Sint32Kind, protoreflect.Sint64Kind,
		protoreflect.Fixed32Kind, protoreflect.Fixed64Kind,
		protoreflect.Sfixed32Kind, protoreflect.Sfixed64Kind:
		g.P(`"type": "integer"`)
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		g.P(`"type": "number"`)
	case protoreflect.StringKind:
		g.P(`"type": "string"`)
	case protoreflect.BytesKind:
		g.P(`"type": "string", "contentEncoding": "base64"`)
	case protoreflect.MessageKind:
		if message != nil {
			g.P(fmt.Sprintf(`"$ref": "#/definitions/%s"`, message.GoIdent.GoName))
		}
	case protoreflect.EnumKind:
		// TODO: handle enums
		g.P(`"type": "string"`)
	}
}
