package router

// Router is exported and used in main.go
//func Router() *mux.Router {
//
//	router := mux.NewRouter()
//
//	router.HandleFunc("/api/ticket", middleware.CreateTicketOption).Methods("POST", "OPTIONS")
//	router.HandleFunc("/api/ticket/{id}", middleware.GetTicketOption).Methods("GET", "OPTIONS")
//	router.HandleFunc("/api/ticket/{id}/purchases", middleware.PurchasesFromTicketOptions).Methods("POST", "OPTIONS")
//
//	return router
//}