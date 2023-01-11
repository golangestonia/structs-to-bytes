package main

type MessageDef struct {
	Name   string
	Fields []FieldDef
}

type FieldDef struct {
	Name string
	Type string
}

func Field(name, typ string) FieldDef {
	return FieldDef{
		Name: name,
		Type: typ,
	}
}

func Message(name string, fields ...FieldDef) MessageDef {
	return MessageDef{
		Name:   name,
		Fields: fields,
	}
}
