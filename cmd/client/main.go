package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"train-booking-service/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Define flags
	Operation := flag.String("Operation", "", "Flag to indicate Operation: PurchaseTicket, GetReceipt, ModifySeat, GetUsersBySection, RemoveUser")
	Data := flag.String("Data", "", "Flag to indicate data required for Operations: PurchaseTicket, GetReceipt, ModifySeat, GetUsersBySection, RemoveUser")

	flag.Parse()

	// Connect to the gRPC server
	conn, err := grpc.Dial("localhost:7001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create the client for TrainService
	client := proto.NewTrainServiceClient(conn)

	// Switch based on the operation type
	switch *Operation {
	case "PurchaseTicket":
		// Parse the PurchaseTicketRequest
		var req proto.PurchaseTicketRequest
		if err := json.Unmarshal([]byte(*Data), &req); err != nil {
			log.Fatalf("could not unmarshal PurchaseTicketRequest JSON: %v", err)
		}

		// Call the PurchaseTicket method
		resp, err := client.PurchaseTicket(context.Background(), &req)
		if err != nil {
			log.Fatalf("could not purchase ticket: %v", err)
		}

		// Output the response
		fmt.Printf("Ticket purchased: %s, Seat: %s\n", resp.Message, resp.Ticket.Seat)

	case "GetReceipt":
		// Parse the GetReceiptRequest
		var req proto.GetReceiptRequest
		if err := json.Unmarshal([]byte(*Data), &req); err != nil {
			log.Fatalf("could not unmarshal GetReceiptRequest JSON: %v", err)
		}

		// Call the GetReceipt method
		resp, err := client.GetReceipt(context.Background(), &req)
		if err != nil {
			log.Fatalf("could not get receipt: %v", err)
		}

		// Output the receipt details
		fmt.Printf("Receipt: %s to %s, Seat: %s, Price Paid: %.2f\n", resp.Ticket.From, resp.Ticket.To, resp.Ticket.Seat, resp.Ticket.PricePaid)

	case "ModifySeat":
		// Parse the ModifySeatRequest
		var req proto.ModifySeatRequest
		if err := json.Unmarshal([]byte(*Data), &req); err != nil {
			log.Fatalf("could not unmarshal ModifySeatRequest JSON: %v", err)
		}

		// Call the ModifySeat method
		resp, err := client.ModifySeat(context.Background(), &req)
		if err != nil {
			log.Fatalf("could not modify seat: %v", err)
		}

		// Output the modified seat details
		fmt.Printf("New Receipt: %s to %s, Seat: %s, Price Paid: %.2f\n", resp.NewTicket.From, resp.NewTicket.To, resp.NewTicket.Seat, resp.NewTicket.PricePaid)

	case "RemoveUser":
		// Parse the RemoveUserRequest
		var req proto.RemoveUserRequest
		if err := json.Unmarshal([]byte(*Data), &req); err != nil {
			log.Fatalf("could not unmarshal RemoveUserRequest JSON: %v", err)
		}

		// Call the RemoveUser method
		resp, err := client.RemoveUser(context.Background(), &req)
		if err != nil {
			log.Fatalf("could not remove user: %v", err)
		}

		// Output the removed user details
		fmt.Printf("Removed User: Name: %s %s, Email: %s\n", resp.User.FirstName, resp.User.LastName, resp.User.Email)

	case "GetUsersBySection":
		// Parse the GetUsersBySectionRequest
		var req proto.GetUsersBySectionRequest
		if err := json.Unmarshal([]byte(*Data), &req); err != nil {
			log.Fatalf("could not unmarshal GetUsersBySectionRequest JSON: %v", err)
		}

		// Call the GetUsersBySection method
		resp, err := client.GetUsersBySection(context.Background(), &req)
		if err != nil {
			log.Fatalf("could not get users by section: %v", err)
		}

		// Output the list of users and their seats
		fmt.Println("User List:")
		for _, userSeat := range resp.UserSeats {
			fmt.Printf("\nUser: %s %s (%s)\nSeat: %s\n", userSeat.User.FirstName, userSeat.User.LastName, userSeat.User.Email, userSeat.Seat)
		}

	default:
		log.Fatalf("No valid operation selected")
	}
}
