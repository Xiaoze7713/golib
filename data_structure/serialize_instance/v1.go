package serialize_instance

type SerializationFunc func(i interface{}) string
type DeserializationFunc func(s string, i interface{}) error
