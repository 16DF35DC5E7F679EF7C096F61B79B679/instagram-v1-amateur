package user

import "go.mongodb.org/mongo-driver/bson"

func createFilterToFindByHandle(handle string) bson.M {
	return bson.M{
		"handle": handle,
	}
}

func createFilterToFindByUsername(username string) bson.M {
	return bson.M{
		"username" : username,
	}
}

func createFilterToFindByEmail(email string) bson.M {
	return bson.M{
		"email" : email,
	}
}