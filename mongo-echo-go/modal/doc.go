package modal

type Doc struct {
	Name  string `bson:"name"`
	Value string `bson:"value"`
	CC    string `bson:"cc,omitempty"` // Optional field, can be empty
}
