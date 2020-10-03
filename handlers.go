package main

import (
	"com.harsha/go-practices/instagram-v1/user"
	"com.harsha/go-practices/instagram-v1/user_session"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func createSuccessfulResponse(message string, data interface{}) *GenericResponse{
	return &GenericResponse{
		Error:   false,
		Message: message,
		Data:    data,
	}
}

func extractFromRequestURL(r *http.Request, key string) (string, error) {
	vars := mux.Vars(r)
	return vars[key], nil
}

func createBadRequestResponse(message string) *GenericResponse {
	return &GenericResponse{
		Error:   true,
		Message: message,
		Data:    nil,
	}
}

func createInternalServerErrorResponse(errorMessage string) *GenericResponse {
	return &GenericResponse{
		Error:   true,
		Message: errorMessage,
		Data:    nil,
	}
}


func respondWithJson(w http.ResponseWriter, response *GenericResponse, responseCode int) {
	responseBody, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseCode)
	w.Write(responseBody)
}

func (server *Server) index(w http.ResponseWriter, r *http.Request) {
	respondWithJson(w,
		createSuccessfulResponse("Homepage Showing", "Welcome to homepage of instagram"),
		http.StatusOK,
	)
	return
}

func (server *Server) findUserByHandle(w http.ResponseWriter, r *http.Request) {
	handle, err := extractFromRequestURL(r, "handle")
	if err != nil {
		fmt.Printf("Handle is not present in request %e \n", err)
		respondWithJson(w,
			createBadRequestResponse("Handle not given: "+ err.Error()),
			http.StatusBadRequest,
			)
		return
	}
	userByHandle, err := user.GetUserByHandle(context.TODO(), server.MongoClient, handle)
	if err != nil {
		fmt.Println("Handle is not present in db ")
		fmt.Println(err)
		respondWithJson(w,
			createBadRequestResponse("Handle not present in DB: "+ err.Error()),
			http.StatusBadRequest,
		)
		return
	}
	respondWithJson(w,
		createSuccessfulResponse("User with handle found ", userByHandle),
		http.StatusOK,
	)
	return
}

func (server *Server) createUser(w http.ResponseWriter, r *http.Request) {
	var userCreationRequestDTO user.UserCreationRequestDTO
	err := json.NewDecoder(r.Body).Decode(&userCreationRequestDTO)
	if err != nil {
		respondWithJson(w, createBadRequestResponse("Please supply a valid body: Errors:\n" + err.Error()), http.StatusBadRequest)
		return
	}
	err, responseDTO := user.CreateUser(context.Background(), server.MongoClient, userCreationRequestDTO)
	if err != nil {
		fmt.Println("Error occurred in user creation")
		fmt.Println(err)
		respondWithJson(w, createBadRequestResponse("Couldn't create user: "+ err.Error()), http.StatusBadRequest)
		return
	}
	sessionResponseDTO, err := user_session.CreateSession(context.Background(), server.MongoClient, createSessionCreationDTOForSignUp(&userCreationRequestDTO))
	if err != nil {
		fmt.Printf("User created but error in creating session: %e ", err)
		respondWithJson(w, createInternalServerErrorResponse("User created. Please sign-in again"), http.StatusInternalServerError)
		return
	}
	responseDTO.SessionToken = sessionResponseDTO.SessionToken
	respondWithJson(w, createSuccessfulResponse("Successfully created user", responseDTO), http.StatusAccepted)
}

func (server *Server) signIn(w http.ResponseWriter, r *http.Request) {
	var signInRequestDTO user.SignInRequestDTO
	err := json.NewDecoder(r.Body).Decode(&signInRequestDTO)
	if err != nil {
		respondWithJson(w, createBadRequestResponse("Please supply a valid body: Errors:\n" + err.Error()), http.StatusBadRequest)
		return
	}
	sessionResponseDTO, err := user_session.CreateSession(context.Background(), server.MongoClient, createSessionCreationDTOForSignIn(signInRequestDTO))
	if err != nil {
		fmt.Printf("Error in creating session: %e ", err)
		respondWithJson(w, createInternalServerErrorResponse(err.Error()), http.StatusInternalServerError)
		return
	}
	respondWithJson(w, createSuccessfulResponse("Successfully signed-in", sessionResponseDTO), http.StatusAccepted)
}

func (server *Server) registerDevice(w http.ResponseWriter, r *http.Request) {
	var registerDeviceRequestDTO user_session.RegisterDeviceRequestDTO
	err := json.NewDecoder(r.Body).Decode(&registerDeviceRequestDTO)
	if err != nil {
		respondWithJson(w, createBadRequestResponse("Please supply a valid body: Errors:\n" + err.Error()), http.StatusBadRequest)
		return
	}
	registerDeviceResponseDTO, err := user_session.RegisterDevice(context.Background(), server.MongoClient, registerDeviceRequestDTO)
	if err != nil {
		fmt.Printf("Error occurred in registering device: %e ", err)
		respondWithJson(w, createInternalServerErrorResponse(err.Error()), http.StatusInternalServerError)
		return
	}
	respondWithJson(w, createSuccessfulResponse("Successfully registered device", registerDeviceResponseDTO), http.StatusOK)
}

func (server *Server) signOut(w http.ResponseWriter, r *http.Request) {
	var invalidateSessionRequestDTO user_session.InvalidateSessionRequestDTO
	err := json.NewDecoder(r.Body).Decode(&invalidateSessionRequestDTO)
	if err != nil {
		respondWithJson(w, createBadRequestResponse("Please supply a valid body: Errors:\n" + err.Error()), http.StatusBadRequest)
		return
	}
	err = user_session.InvalidateSession(context.Background(), server.MongoClient, &invalidateSessionRequestDTO)
	if err != nil {
		fmt.Printf("Error in invalidating sessions: %e ", err)
		respondWithJson(w, createSuccessfulResponse("Error occurred in singing out", ""), http.StatusInternalServerError)
		return
	}
	respondWithJson(w, createSuccessfulResponse("Successfully signed out", ""), http.StatusOK)
}

func (server *Server) getAllActiveSessions(w http.ResponseWriter, r *http.Request) {
	handle, err := extractFromRequestURL(r, "handle")
	if err != nil {
		fmt.Printf("[getAllActiveSessions]Error in extracting handle from request %e ", err)
		respondWithJson(w, createBadRequestResponse("Please provide a valid handle"), http.StatusBadRequest)
		return
	}
	allActiveSessions, err := user_session.GetAllActiveSessions(context.Background(), server.MongoClient, handle)
	if err != nil {
		fmt.Printf("Error in fetching active sessoins for handle %s %e ", handle, err)
		respondWithJson(w, createBadRequestResponse(err.Error()), http.StatusBadRequest)
		return
	}
	respondWithJson(w, createSuccessfulResponse("Successfully fetched active sessions", allActiveSessions), http.StatusOK)
}

func createSessionCreationDTOForSignIn(dto user.SignInRequestDTO) *user_session.CreateSessionRequestDTO {
	return &user_session.CreateSessionRequestDTO{
		Handle:      dto.Handle,
		Password:    dto.Password,
		DeviceId:    dto.DeviceId,
		BrowserType: dto.BrowserType,
		Timestamp:   0,
	}
}

func createSessionCreationDTOForSignUp(userCreationRequestDTO *user.UserCreationRequestDTO) *user_session.CreateSessionRequestDTO {
	return &user_session.CreateSessionRequestDTO{
		Handle:      userCreationRequestDTO.Handle,
		Password:    userCreationRequestDTO.Password,
		DeviceId:    userCreationRequestDTO.DeviceId,
		BrowserType: userCreationRequestDTO.BrowserType,
		Timestamp:   time.Now().Unix(),
	}
}