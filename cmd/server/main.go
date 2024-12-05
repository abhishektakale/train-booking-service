package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"train-booking-service/dao"
	"train-booking-service/proto"
	"google.golang.org/grpc/reflection"

	"google.golang.org/grpc"
)

type TrainServiceServer struct {
	proto.UnimplementedTrainServiceServer
	dao *dao.TrainDAO
}

func NewTrainServiceServer() *TrainServiceServer {
	return &TrainServiceServer{
		dao: dao.NewTrainDAO(),
	}
}

func (s *TrainServiceServer) PurchaseTicket(ctx context.Context, req *proto.PurchaseTicketRequest) (*proto.TicketPurchaseResponse, error) {
	log.Printf("PurchaseTicket for User=%s initiated", req.User.Email)

	ticket, err := s.dao.SaveTicket(req.User, req.From, req.To)
	if err != nil {
		log.Printf("Error saving ticket for user %s: %v", req.User, err)
		return nil, err
	}

	log.Printf("Ticket purchased successfully for user %s", req.User)
	return &proto.TicketPurchaseResponse{
		Ticket:  ticket,
		Message: "Ticket purchased successfully",
	}, nil
}

func (s *TrainServiceServer) GetReceipt(ctx context.Context, req *proto.GetReceiptRequest) (*proto.GetReceiptResponse, error) {
	log.Printf("GetReceipt: UserEmail=%s", req.UserEmail)

	ticket, err := s.dao.GetTicket(req.UserEmail)
	if err != nil {
		log.Printf("Error retrieving ticket for user %s: %v", req.UserEmail, err)
		return nil, err
	}

	return &proto.GetReceiptResponse{Ticket: ticket}, nil
}

func (s *TrainServiceServer) GetUsersBySection(ctx context.Context, req *proto.GetUsersBySectionRequest) (*proto.GetUsersBySectionResponse, error) {
	// Validate section (e.g., ensure it is a valid string or within allowed sections)
	if req.Section == "" {
		log.Printf("Invalid section: %s", req.Section)
		return nil, fmt.Errorf("invalid section: %s", req.Section)
	}

	log.Printf("GetUsersBySection: Section=%s", req.Section)

	sectionTickets, err := s.dao.GetUsersBySection(req.Section)
	if err != nil {
		log.Printf("Error fetching users by section %s: %v", req.Section, err)
		return nil, err
	}

	var usersInSection []*proto.UserSeatAllocation
	for _, ticket := range sectionTickets {
		usersInSection = append(usersInSection, &proto.UserSeatAllocation{
			User: ticket.User,
			Seat: ticket.Seat,
		})
	}

	return &proto.GetUsersBySectionResponse{UserSeats: usersInSection}, nil
}

func (s *TrainServiceServer) ModifySeat(ctx context.Context, req *proto.ModifySeatRequest) (*proto.ModifySeatResponse, error) {
	log.Printf("ModifySeat: UserEmail=%s, NewSeat=%s", req.UserEmail, req.NewSeat)

	ticket, err := s.dao.GetTicket(req.UserEmail)
	if err != nil {
		log.Printf("Error retrieving ticket for user %s: %v", req.UserEmail, err)
		return nil, err
	}

	// Deallocate the old seat and allocate the new seat
	if err := s.dao.ModifySeat(ticket.Seat, req.NewSeat, req.UserEmail); err != nil {
		log.Printf("Error modifying seat for user %s: %v", req.UserEmail, err)
		return nil, err
	}

	ticket.Seat = req.NewSeat
	log.Printf("Seat modified successfully for user %s", req.UserEmail)
	return &proto.ModifySeatResponse{NewTicket: ticket, Message: "Seat modified successfully"}, nil
}

func (s *TrainServiceServer) RemoveUser(ctx context.Context, req *proto.RemoveUserRequest) (*proto.RemoveUserResponse, error) {
	log.Printf("RemoveUser: UserEmail=%s", req.UserEmail)

	ticket, err := s.dao.GetTicket(req.UserEmail)
	if err != nil {
		log.Printf("Error retrieving ticket for user %s: %v", req.UserEmail, err)
		return nil, err
	}

	deletedTicket, err := s.dao.DeleteTicket(ticket)
	if err != nil {
		log.Printf("Error deleting ticket for user %s: %v", req.UserEmail, err)
		return nil, err
	}

	log.Printf("User removed successfully: %s", req.UserEmail)
	return &proto.RemoveUserResponse{User: deletedTicket.User, Message: "User removed successfully"}, nil
}

func main() {
	server := grpc.NewServer()
	proto.RegisterTrainServiceServer(server, NewTrainServiceServer())
	reflection.Register(server)
	listener, err := net.Listen("tcp", ":7001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("Server is listening on port 7001...")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
