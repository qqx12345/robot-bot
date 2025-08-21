package message

var HandlerFunc = map[string]func(data map[string]interface{},ID string){
	"C2C_MESSAGE_CREATE":Chat,
}

func Message(data map[string]interface{},ID string,T string) (interface{}, error) {
	HandlerFunc[T](data,ID)
	return nil,nil
}