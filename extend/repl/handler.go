package repl

// routeMessage creates msg interface from LocalHandler
func routeMessage(route string) (interface{}, error) {
	//handler, err := localHandler.RouteHandler(route)
	//if err != nil {
	//	return nil, fmt.Errorf("unexpected route:%s, can not find it's route handler", route)
	//}
	//return reflect.New(handler.Type.Elem()).Interface(), nil
	m := make(map[string]interface{})
	return &m, nil
}
