package http

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

type APIServer struct {
	listenAddr string
}

func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}
func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/token/{userID}", makeHTTPHandleFunc(s.handleGetToken)).Methods("GET")
	router.HandleFunc("/token/{userID}", makeHTTPHandleFunc(s.handleRefreshToken)).Methods("POST")

	http.ListenAndServe(s.listenAddr, router)
}
