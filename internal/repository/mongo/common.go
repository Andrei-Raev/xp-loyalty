package mongo

import "go.mongodb.org/mongo-driver/bson/primitive"

func toPrimitives(ids []string) []primitive.ObjectID {
	ids_ := make([]primitive.ObjectID, len(ids))
	for i := range ids_ {
		ids_[i], _ = primitive.ObjectIDFromHex(ids[i])
	}
	return ids_
}
