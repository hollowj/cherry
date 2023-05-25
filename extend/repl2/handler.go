package repl2

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// routeMessage creates msg interface from LocalHandler
func routeMessage(route string) (interface{}, error) {
	lastDotIndex := strings.LastIndex(route, ".")
	msgName := fmt.Sprintf("gamemsg.%v", route[lastDotIndex+1:])
	msgType, err := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(msgName))
	if err != nil {
		if err == protoregistry.NotFound {
			data := make(map[string]interface{})
			return &data, nil
		}
		return nil, err
	}
	m := msgType.New().Interface()
	return &m, nil
}
