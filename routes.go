package main

func (server *Server) initializeRoutes() {
	server.Router.HandleFunc("/", loggedAction(server.index)).Methods("GET")
	server.Router.HandleFunc("/handle/{handle}", loggedAction(server.findUserByHandle)).Methods("GET")
	server.Router.HandleFunc("/user", loggedAction(server.createUser)).Methods("POST")
	server.Router.HandleFunc("/signIn", loggedAction(server.signIn)).Methods("POST")
	server.Router.HandleFunc("/register/device", loggedAction(server.registerDevice)).Methods("POST")
	server.Router.HandleFunc("/signOut", loggedAction(server.signOut)).Methods("POST")
	server.Router.HandleFunc("/sessions/{handle}", loggedAction(server.getAllActiveSessions)).Methods("GET")
}
