package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"os"

	"github.com/Alfazal007/mail-worker/send_mail"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

type mynotificationService struct {
	send_mail.UnimplementedSendMailServiceServer
}

var fromAddress string
var password string

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	portNumber := os.Getenv("PORT")
	if portNumber == "" {
		log.Fatal("Error getting port number from the env variables")
	}
	fromAddress = os.Getenv("FROM_ADDRESS")
	if fromAddress == "" {
		log.Fatal("Error getting from address from the env variables")
	}
	password = os.Getenv("PASSWORD")
	if password == "" {
		log.Fatal("Error getting password from the env variables")
	}
	lis, err := net.Listen("tcp", ":"+portNumber)
	if err != nil {
		log.Fatalf("Cannot create grpc listener: %s", err)
	}
	serverRegistrar := grpc.NewServer()
	service := mynotificationService{}
	send_mail.RegisterSendMailServiceServer(serverRegistrar, &service)
	err = serverRegistrar.Serve(lis)
	if err != nil {
		log.Fatalf("impossible to serve: %s", err)
	}
}

func (s *mynotificationService) SendMail(ctx context.Context, req *send_mail.CreateRequest) (*send_mail.CreateResponse, error) {

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Message.
	message := []byte("Subject: New tweet by " + req.Fromusername + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
		"\r\n" +
		"<html><body>" +
		"<p>The user you follow has created a new tweet, click the below link to visit.</p>" +
		"<p><a href=" + req.Link + ">Click to view</a></p>" +
		"</body></html>")

	// Authentication.
	auth := smtp.PlainAuth("", fromAddress, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, fromAddress, req.To, message)
	if err != nil {
		fmt.Println(err)
		return &send_mail.CreateResponse{
			Done: false,
		}, err
	}
	fmt.Println("Email Sent Successfully!")
	return &send_mail.CreateResponse{
		Done: true,
	}, nil
}
