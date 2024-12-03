package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"train-booking-service/dao"
	"train-booking-service/proto"

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
	seat, err := s.dao.AssignSeat()
	if err != nil {
		return nil, err
	}

	ticket := &proto.TicketReceipt{
		From:      req.From,
		To:        req.To,
		User:      req.User,
		PricePaid: 20.0,
		Seat:      seat,
	}
	s.dao.SaveTicket(req.User.Email, ticket)

	return &proto.TicketPurchaseResponse{
		Ticket:  ticket,
		Message: "Ticket purchased successfully",
	}, nil
}

func (s *TrainServiceServer) GetReceipt(ctx context.Context, req *proto.GetReceiptRequest) (*proto.GetReceiptResponse, error) {

	ticket, err := s.dao.GetTicket(req.UserEmail)
	if err != nil {
		return nil, err
	}

	return &proto.GetReceiptResponse{Ticket: ticket}, nil
}

func (s *TrainServiceServer) GetUsersBySection(ctx context.Context, req *proto.GetUsersBySectionRequest) (*proto.GetUsersBySectionResponse, error) {

	sectionTickets := s.dao.GetUsersBySection(req.Section)
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
	ticket, err := s.dao.GetTicket(req.UserEmail)
	if err != nil {
		return nil, err
	}

	if !s.dao.IsSeatAvailable(req.NewSeat) {
		return nil, fmt.Errorf("seat %s is already taken", req.NewSeat)
	}

	// Deallocate the old seat and allocate the new seat
	s.dao.DeallocateSeat(ticket.Seat)
	if err := s.dao.AllocateSeat(req.NewSeat); err != nil {
		return nil, err
	}

	ticket.Seat = req.NewSeat
	return &proto.ModifySeatResponse{Message: "Seat modified successfully"}, nil
}

func (s *TrainServiceServer) RemoveUser(ctx context.Context, req *proto.RemoveUserRequest) (*proto.RemoveUserResponse, error) {
	_, err := s.dao.DeleteTicket(req.UserEmail)
	if err != nil {
		return nil, err
	}
	return &proto.RemoveUserResponse{Message: "User removed successfully"}, nil
}

func main() {
	server := grpc.NewServer()
	proto.RegisterTrainServiceServer(server, NewTrainServiceServer())

	listener, err := net.Listen("tcp", ":7001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Println("Server is listening on port 50051...")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
