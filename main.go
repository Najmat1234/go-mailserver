package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/smtp"
	"sync"
)

var mu sync.Mutex
var logger = log.Default()

// Email details
const (
	smtpHost       = "smtp.gmail.com"           // Replace with your SMTP server
	smtpPort       = "587"                      // Typically 587 for TLS
	senderEmail    = "shashank897949@gmail.com" // Replace with your email
	senderPassword = "azhb mjiw lwve uyye"      // Replace with your password
	recipientEmail = "test771df@yopmail.com"    // Replace with recipient email
)

// Request
type EmailRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// Response
type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func JsonResponse(w http.ResponseWriter, status int, message string, data any) {
	response := &Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)

}

func SendEmail(req EmailRequest) (bool, error) {
	// Set up authentication information
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)

	// Compose the email
	subject := req.Subject
	body := req.Message + "\n" + "From: " + req.Name + "(" + req.Email + ")"
	message := []byte(subject + "\n" + body)

	logger.Println("Sending email...")
	// Send the email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, []string{recipientEmail}, message)
	if err != nil {
		logger.Printf("Failed to send email: %v", err)
		return false, err
	}
	return true, nil
}

func SenderHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	if r.Method != http.MethodPost {
		logger.Printf("Bad request %v", http.StatusBadRequest)
		JsonResponse(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), nil)
		return
	}
	var req EmailRequest
	jsonErr := json.NewDecoder(r.Body).Decode(&req)
	if jsonErr != nil {
		logger.Printf("Error during decode body %v", jsonErr)
		JsonResponse(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), nil)
		return
	}
	status, err := SendEmail(req)
	if err != nil && !status {
		logger.Printf("Error during sending email %v", err)
		JsonResponse(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), nil)
		return
	}
	logger.Printf("Contact Request from %s", req.Email)
	JsonResponse(w, http.StatusOK, http.StatusText(http.StatusOK), "Email sent successfully!")

}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	logger.Println("Health route hit")
	JsonResponse(w, http.StatusOK, http.StatusText(http.StatusOK), "Healthy")
}

func main() {

	defer func() {
		r := recover()
		if r != nil {
			logger.Printf("Recover from panic %v", r)
		}
	}()

	http.HandleFunc("/health", HealthHandler)
	http.HandleFunc("/send-mail", SenderHandler)
	// Log when the server starts
	log.Println("Server is starting on port 8080...")

	// Start the server and log any errors
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
