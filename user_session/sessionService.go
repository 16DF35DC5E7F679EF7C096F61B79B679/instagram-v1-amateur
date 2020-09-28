package user_session

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)


func CreateSession(ctx context.Context, mongoClient *mongo.Client, createSessionRequestDTO *CreateSessionRequestDTO) (*SessionResponseDTO, error) {
	session, err := validateAndCreateSession(ctx, mongoClient, createSessionRequestDTO)
	if err != nil {
		return nil, err
	}
	newCtx, cancelFunc := context.WithTimeout(ctx, time.Second * 10)
	defer cancelFunc()
	invalidateOtherSessionsInTheDeviceAndBrowser(ctx, mongoClient, createSessionRequestDTO)
	collection := mongoClient.Database("instagram-v1").Collection("session")
	sessionCreationResult, err := collection.InsertOne(newCtx, session)
	if err != nil {
		return nil, err
	}
	return createSessionResponseDTO(session, sessionCreationResult.InsertedID.(primitive.ObjectID)), nil
}

func validateAndCreateSession(ctx context.Context, client *mongo.Client, dto *CreateSessionRequestDTO) (*Session, error) {
	err := validatePassword(ctx, client, dto.Handle, dto.Password)
	if err != nil {
		return nil, err
	}
	err = validateDeviceId(ctx, client, dto.DeviceId)
	if err != nil {
		return nil, err
	}
	return createSession(dto)
}

func createSession(dto *CreateSessionRequestDTO) (*Session, error) {
	return &Session{
		Id:           primitive.NewObjectID(),
		SessionToken: generateToken(dto.Handle, dto.DeviceId, dto.BrowserType),
		Handle:       dto.Handle,
		DeviceId:     dto.DeviceId,
		BrowserType:  dto.BrowserType,
		ActiveTill:   time.Now().Unix() + 60 * 60 * 24 * 15,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}, nil
}


func validatePassword(ctx context.Context, client *mongo.Client, handle string, password string) error {
	return verifyPassword(ctx, client, handle, password)
}

func validateDeviceId(ctx context.Context, client *mongo.Client, deviceId string) error {
	collection := client.Database("instagram-v1").Collection("device")
	newCtx, cancelFunc := context.WithTimeout(ctx, time.Second * 5)
	defer cancelFunc()
	deviceIdObjectId, err := primitive.ObjectIDFromHex(deviceId)
	if err != nil {
		return &InvalidSessionCreationError{
			InvalidField: "device_id",
			ErrorReason:  "Invalid Hex String",
		}
	}
	deviceByDeviceIdCount, err := collection.CountDocuments(newCtx, &bson.M{"_id" : deviceIdObjectId})
	if err != nil {
		return err
	}
	if deviceByDeviceIdCount == 0 {
		return &InvalidSessionCreationError{
			InvalidField: "device_id",
			ErrorReason:  "Please register the device first",
		}
	}
	return nil
}

func createSessionResponseDTO(session *Session, id primitive.ObjectID) *SessionResponseDTO {
	return &SessionResponseDTO{Id: id.Hex(), Handle: session.Handle, SessionToken: session.SessionToken, ActiveTill: session.ActiveTill}
}

func invalidateOtherSessionsInTheDeviceAndBrowser(ctx context.Context, client *mongo.Client, dto *CreateSessionRequestDTO) {
	collection := client.Database("instagram-v1").Collection("session")
	filterToBeFindCurrentlyActiveSessionsInTheSameDeviceAndBrowser := &bson.M{
		"handle" : dto.Handle,
		"device_id" : dto.DeviceId,
		"browser_type" : dto.BrowserType,
		"active_till" : bson.M{
			"$gte" : time.Now().Unix(),
		},
	}
	filterToInvalidteCurrentSession := &bson.M{
		"$set" : bson.M{
			"active_till" : time.Now().Unix(),
		},
	}
	newCtx, cancelFunc := context.WithTimeout(ctx, time.Second * 10)
	defer cancelFunc()
	invalidateSessionsResult, err := collection.UpdateMany(newCtx,
			filterToBeFindCurrentlyActiveSessionsInTheSameDeviceAndBrowser,
			filterToInvalidteCurrentSession)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Number of sessions invalidated: %d ", invalidateSessionsResult.MatchedCount)
	//cursor, err := collection.Find(newCtx, filter)
	//if err != nil {
	//	fmt.Println("Couldn't find existing sessions in database")
	//	return
	//}
	//for cursor.Next(context.TODO()) {
	//	var activeSession Session
	//	err := cursor.Decode(&activeSession)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	activeSession.ActiveTill = time.Now().Unix()
	//	go updateSession(newCtx, collection, activeSession)
	//}
}

//func invalidateSession(ctx context.Context, sessionCollection *mongo.Collection, session Session) {
//	sessionCollection.UpdateOne(context.TODO(), bson.M{"_id": session.Id}, )
//}

func verifyPassword(ctx context.Context, client *mongo.Client, handle string, passwordText string) error {
	passwordMatchingFilter := &bson.M{
		"handle" : handle,
		"password": passwordText,
	}
	newCtx, cancelFunc := context.WithTimeout(ctx, time.Second * 10)
	defer cancelFunc()
	count, err := client.Database("instagram-v1").Collection("user").CountDocuments(newCtx, passwordMatchingFilter)
	if err != nil {
		return err
	}
	if count != 1 {
		return &PasswordMismatchError{
			RootCause: "Invalid credentials",
		}
	}
	return nil
}

func RegisterDevice(ctx context.Context, client *mongo.Client, dto RegisterDeviceRequestDTO) (*RegisterDeviceResponseDTO, error) {
	newCtx, cancelFunc := context.WithTimeout(ctx, time.Second * 10)
	defer cancelFunc()
	collection := client.Database("instagram-v1").Collection("device")
	singleResult := collection.FindOne(newCtx, &bson.M{"device_name": dto.DeviceName, "device_ip": dto.DeviceIP})
	if singleResult.Err() == nil {
		alreadyExistingDevice := &Device{}
		err := singleResult.Decode(&alreadyExistingDevice)
		if err != nil {
			return nil, err
		}
		if alreadyExistingDevice != nil {
			return &RegisterDeviceResponseDTO{Id: alreadyExistingDevice.Id.Hex()}, nil
		}
	}
	//TODO Verify IP
	singleInsertionResult, err := collection.InsertOne(newCtx, &Device{
		Id:         primitive.NewObjectID(),
		DeviceName: dto.DeviceName,
		DeviceIP: dto.DeviceIP,
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	})
	if err != nil {
		return nil, err
	}
	return &RegisterDeviceResponseDTO{Id: singleInsertionResult.InsertedID.(primitive.ObjectID).Hex()}, nil
}