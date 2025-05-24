package modal

type Test struct {
	Status string `bson:"status"`
	Type   string `bson:"type"`
	Time   string `bson:"time"`
	Name   string `bson:"name"`
	Value  string `bson:"value"`
}
