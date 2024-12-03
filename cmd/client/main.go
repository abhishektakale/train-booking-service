package main

import (
	"context"
	"fmt"
	"log"
	"train-booking-service/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:7001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := proto.NewTrainServiceClient(conn)

	// Purchase a ticket
	user := &proto.User{
		FirstName: "Abhishek",
		LastName:  "Takale",
		Email:     "abhishek.takale1995@gmail.com",
	}
	req := &proto.PurchaseTicketRequest{
		From: "London",
		To:   "Paris",
		User: user,
	}
	resp, err := client.PurchaseTicket(context.Background(), req)
	if err != nil {
		log.Fatalf("could not purchase ticket: %v", err)
	}
	fmt.Printf("Ticket purchased: %s, Seat: %s\n", resp.Message, resp.Ticket.Seat)

	// Get receipt
	receiptReq := &proto.GetReceiptRequest{UserEmail: user.Email}
	receiptResp, err := client.GetReceipt(context.Background(), receiptReq)
	if err != nil {
		log.Fatalf("could not get receipt: %v", err)
	}
	fmt.Printf("Receipt: %s to %s, Seat: %s, Price Paid: %.2f\n", receiptResp.Ticket.From, receiptResp.Ticket.To, receiptResp.Ticket.Seat, receiptResp.Ticket.PricePaid)
}
