package homehub

import "strings"

var responseTypeMappings = map[string]string{
	// Standard types
	"xmo:boolean": "bool",
	"xmo:int":     "int",
	"xmo:int32":   "int",
	"xmo:int64":   "int64",
	"xmo:number":  "int",
	"xmo:str":     "string",
	"xmo:time":    "string",
	"xmo:uint32":  "int",
	"xmo:uint64":  "int64",

	// Extended types
	"deviceconfig:LastSuccesfulWanType":     "string",
	"interface:Interface:Status":            "string",
	"managers:LedEnable":                    "string",
	"wifi:AccessPoint:Security:ModeEnabled": "string",
}

func getTypeMapping(typeNames string) string {
	types := strings.Split(typeNames, " ")
	for _, typeName := range types {
		mapping, exists := responseTypeMappings[typeName]
		if exists {
			return mapping
		}
	}

	//TODO: This should probably be treated as an error
	return "string"
}
