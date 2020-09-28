package user

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const MinimumAllowedAge = 14

func getEpochXYearsAgo(x int) int64 {
	return time.Now().AddDate((-1) * x, 0, 0).Unix()
}

func GetUserByHandle(ctx context.Context, client *mongo.Client, handle string) (*User, error){
	newCtx, cancelFunc := context.WithTimeout(ctx, 10 * time.Second)
	defer cancelFunc()
	collection := client.Database("instagram-v1").Collection("user")
	filter := &bson.M{
		"handle": handle,
	}
	fmt.Println("Finding user by handle: "+ handle)
	singleResult := collection.FindOne(newCtx, filter)
	user := &User{}
	err := singleResult.Decode(user)
	if err != nil {
		fmt.Printf("Error occurred in finding users by handle: %s, %e \n", handle, err)
		return nil, err
	}
	fmt.Println("User is ")
	fmt.Println(user)
	return user, nil
}

func CreateUser(ctx context.Context, client *mongo.Client, userCreationDTO UserCreationRequestDTO) (error, *UserResponseDTO) {
	newCtx, cancelFunc := context.WithTimeout(ctx, time.Second * 10)
	defer cancelFunc()
	collection := client.Database("instagram-v1").Collection("user")
	err := validateUserCreation(ctx, collection, userCreationDTO)
	if err != nil {
		return err, nil
	}
	user := createUserFromRequestDTO(userCreationDTO)
	userCreationResult, err := collection.InsertOne(newCtx, user)
	if err != nil {
		fmt.Println("Error occurred in inserting user")
		fmt.Println(err)
		return err, nil
	}
	userResponseDTO := &UserResponseDTO{
		Id:   userCreationResult.InsertedID.(primitive.ObjectID).Hex(),
		Name: user.Name,
	}
	return nil, userResponseDTO
}

func validateUserCreation(ctx context.Context, userCollection *mongo.Collection, dto UserCreationRequestDTO) error {
	err := validateIfDOBIsValid(dto.DOB)
	if err != nil {
		return err
	}
	err = validateIfHandleIsUnique(ctx, userCollection, dto.Handle)
	if err != nil {
		return err
	}
	err = validateIfUsernameIsUnique(ctx, userCollection, dto.Username)
	if err != nil {
		return err
	}
	err = validateIfEmailIsUnique(ctx, userCollection, dto.Email)
	if err != nil {
		return err
	}
	err = validateIfPasswordIsValid(dto.Password)
	if err != nil {
		return err
	}
	return nil
}

func validateIfPasswordIsValid(password string) error {
	if len(password) <= 7 {
		return &CreationValidationError{
			InvalidField: "password",
			ErrorReason:  "Should be more than 7 characters",
		}
	}
	return nil
}

func validateIfDOBIsValid(dob int64) error {
	if dob > getEpochXYearsAgo(MinimumAllowedAge) {
		return &CreationValidationError{
			InvalidField: "DOB",
			ErrorReason:  fmt.Sprintf("Get outta here kiddo! You should be at least %d years old to use this app", MinimumAllowedAge),
		}
	}
	return nil
}

func validateIfUsernameIsUnique(ctx context.Context, userCollection *mongo.Collection, username string) error {
	newCtx, cancelFunc := context.WithTimeout(ctx, time.Second * 5)
	defer cancelFunc()
	existingUserCountBySameUsername, err := userCollection.CountDocuments(newCtx, createFilterToFindByUsername(username))
	if err != nil {
		return err
	}
	if existingUserCountBySameUsername != 0 {
		return &CreationValidationError{
			InvalidField: "username",
			ErrorReason:  "Already Occupied",
		}
	}
	return nil
}

func validateIfHandleIsUnique(ctx context.Context, userCollection *mongo.Collection, handle string) error {
	newCtx, cancelFunc := context.WithTimeout(ctx, time.Second * 5)
	defer cancelFunc()
	existingUserCountBySameHandle, err := userCollection.CountDocuments(newCtx, createFilterToFindByHandle(handle))
	if err != nil {
		return err
	}
	if existingUserCountBySameHandle != 0 {
		return &CreationValidationError{
			InvalidField: "handle",
			ErrorReason:  "Already Occupied",
		}
	}
	return nil
}

func validateIfEmailIsUnique(ctx context.Context, userCollection *mongo.Collection, email string) error {
	newCtx, cancelFunc := context.WithTimeout(ctx, time.Second * 5)
	defer cancelFunc()
	existingUserCountBySameEmail, err := userCollection.CountDocuments(newCtx, createFilterToFindByEmail(email))
	if err != nil {
		return err
	}
	if existingUserCountBySameEmail != 0 {
		return &CreationValidationError{
			InvalidField: "email",
			ErrorReason:  "Already Occupied",
		}
	}
	return nil
}

func createUserFromRequestDTO(dto UserCreationRequestDTO) *User {
	return &User{
		Id: primitive.NewObjectID(),
		Name:     dto.Name,
		Username: dto.Username,
		Email: dto.Email,
		Handle:   dto.Handle,
		DOB: dto.DOB,
		Password: dto.Password,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
}
