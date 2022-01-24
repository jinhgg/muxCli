package types

type User struct {
	Id   uint8  `bson:"id"`
	Name string `bson:"name"`
	Age  uint8  `bson:"age"`
}
