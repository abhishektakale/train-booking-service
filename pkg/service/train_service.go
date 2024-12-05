// service/train_service.go

package service

import (
	"context"
	"fmt"
	"log"
	"train-booking-service/pkg/adapter"
	"train-booking-service/proto"
)

type TrainService struct {
	proto.UnimplementedTrainServiceServer
	dao adapter.TrainDAOAdapter
}

func NewTrainService(dao adapter.TrainDAOAdapter) *TrainService {
	return &TrainService{
		dao: dao,
	}
}

// PurchaseTicket handles the purchase of a new ticket.
func (s *TrainService) PurchaseTicket(ctx context.Context, req *proto.PurchaseTicketRequest) (*proto.TicketPurchaseResponse, error) {
	log.Printf("PurchaseTicket for User=%s initiated", req.User.Email)

	ticket, err := s.dao.SaveTicket(req.User, req.From, req.To)
	if err != nil {
		log.Printf("Error saving ticket for user %s: %v", req.User.Email, err)
		return nil, fmt.Errorf("could not save ticket: %v", err)
	}

	log.Printf("Ticket purchased successfully for user %s", req.User.Email)
	return &proto.TicketPurchaseResponse{
		Ticket:  ticket,
		Message: "Ticket purchased successfully",
	}, nil
}

// GetReceipt retrieves a user's ticket receipt.
func (s *TrainService) GetReceipt(ctx context.Context, req *proto.GetReceiptRequest) (*proto.GetReceiptResponse, error) {
	log.Printf("GetReceipt: UserEmail=%s", req.UserEmail)

	ticket, err := s.dao.GetTicket(req.UserEmail)
	if err != nil {
		log.Printf("Error retrieving ticket for user %s: %v", req.UserEmail, err)
		return nil, fmt.Errorf("could not retrieve ticket: %v", err)
	}

	return &proto.GetReceiptResponse{Ticket: ticket}, nil
}

// GetUsersBySection retrieves users assigned to a specific section.
func (s *TrainService) GetUsersBySection(ctx context.Context, req *proto.GetUsersBySectionRequest) (*proto.GetUsersBySectionResponse, error) {
	// Validate section (could check against a list of allowed sections)
	if req.Section == "" {
		log.Printf("Invalid section: %s", req.Section)
		return nil, fmt.Errorf("invalid section: %s", req.Section)
	}

	log.Printf("GetUsersBySection: Section=%s", req.Section)

	sectionTickets, err := s.dao.GetUsersBySection(req.Section)
	if err != nil {
		log.Printf("Error fetching users by section %s: %v", req.Section, err)
		return nil, fmt.Errorf("could not fetch users by section: %v", err)
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

// ModifySeat updates a user's seat.
func (s *TrainService) ModifySeat(ctx context.Context, req *proto.ModifySeatRequest) (*proto.ModifySeatResponse, error) {
	log.Printf("ModifySeat: UserEmail=%s, NewSeat=%s", req.UserEmail, req.NewSeat)

	ticket, err := s.dao.GetTicket(req.UserEmail)
	if err != nil {
		log.Printf("Error retrieving ticket for user %s: %v", req.UserEmail, err)
		return nil, fmt.Errorf("could not retrieve ticket: %v", err)
	}

	// Deallocate the old seat and allocate the new seat
	if err := s.dao.ModifySeat(ticket.Seat, req.NewSeat, req.UserEmail); err != nil {
		log.Printf("Error modifying seat for user %s: %v", req.UserEmail, err)
		return nil, fmt.Errorf("could not modify seat: %v", err)
	}

	ticket.Seat = req.NewSeat
	log.Printf("Seat modified successfully for user %s", req.UserEmail)
	return &proto.ModifySeatResponse{NewTicket: ticket, Message: "Seat modified successfully"}, nil
}

// RemoveUser removes a user and their ticket.
func (s *TrainService) RemoveUser(ctx context.Context, req *proto.RemoveUserRequest) (*proto.RemoveUserResponse, error) {
	log.Printf("RemoveUser: UserEmail=%s", req.UserEmail)

	ticket, err := s.dao.GetTicket(req.UserEmail)
	if err != nil {
		log.Printf("Error retrieving ticket for user %s: %v", req.UserEmail, err)
		return nil, fmt.Errorf("could not retrieve ticket: %v", err)
	}

	deletedTicket, err := s.dao.DeleteTicket(ticket)
	if err != nil {
		log.Printf("Error deleting ticket for user %s: %v", req.UserEmail, err)
		return nil, fmt.Errorf("could not delete ticket: %v", err)
	}

	log.Printf("User removed successfully: %s", req.UserEmail)
	return &proto.RemoveUserResponse{User: deletedTicket.User, Message: "User removed successfully"}, nil
}
