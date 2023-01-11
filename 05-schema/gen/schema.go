package main

type SchemaDef struct {
	Messages []MessageDef
}

type MessageDef struct {
	Name    string
	SliceOf string
	Fields  []FieldDef
}

type FieldDef struct {
	Name string
	Type string
}

func Schema(messages ...MessageDef) SchemaDef {
	return SchemaDef{Messages: messages}
}

func Message(name string, fields ...FieldDef) MessageDef {
	return MessageDef{
		Name:   name,
		Fields: fields,
	}
}

func SliceOf(name, typ string) MessageDef {
	return MessageDef{
		Name:    name,
		SliceOf: typ,
	}
}

func Field(name, typ string) FieldDef {
	return FieldDef{
		Name: name,
		Type: typ,
	}
}
