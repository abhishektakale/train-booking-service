package dao

import (
	"fmt"
	"slices"
	"sync"
	"train-booking-service/proto"
)

// Define section constants for type safety
const (
	SectionA   = "A"
	SectionB   = "B"
	SectionCap = 25 // Maximum capacity for Section A and B
)

// TrainDAO is the data access object for managing train seat reservations and user tickets.
type TrainDAO struct {
	users          map[string]*proto.User
	sections       map[string]map[string]*proto.TicketReceipt
	availableSeats map[string][]string
	mu             sync.Mutex
}

// NewTrainDAO initializes a new TrainDAO instance.
func NewTrainDAO() *TrainDAO {
	availableSeats := map[string][]string{
		SectionA: make([]string, SectionCap),
		SectionB: make([]string, SectionCap),
	}
	for i := 0; i < SectionCap; i++ {
		availableSeats[SectionA][i] = fmt.Sprintf("%s%v", SectionA, i+1)
		availableSeats[SectionB][i] = fmt.Sprintf("%s%v", SectionB, i+1)
	}
	return &TrainDAO{
		users: make(map[string]*proto.User),
		sections: map[string]map[string]*proto.TicketReceipt{
			SectionA: make(map[string]*proto.TicketReceipt),
			SectionB: make(map[string]*proto.TicketReceipt),
		},
		availableSeats: availableSeats,
	}
}

func (dao *TrainDAO) newUser(firstName, lastName, email string) *proto.User {
	dao.users[email] = &proto.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}

	return dao.users[email]
}

func (dao *TrainDAO) newTicket(from, to string, user *proto.User, seat string) *proto.TicketReceipt {
	section := string(seat[0])

	ticket := &proto.TicketReceipt{
		From:      from,
		To:        to,
		User:      user,
		PricePaid: 20.0,
		Seat:      seat,
	}
	dao.sections[section][seat] = ticket
	dao.availableSeats[section] = slices.DeleteFunc(dao.availableSeats[section], func(element string) bool {
		return element == seat
	})
	return ticket
}

// AssignSeat assigns the next available seat in the least occupied section.
func (dao *TrainDAO) assignSeat() (string, error) {
	// // Ensure sections exist
	// if (dao.sections[SectionA]) || dao.sectionLimits[SectionB] == 0 {
	// 	return "", fmt.Errorf("invalid section capacities")
	// }

	// Find the section with the least allocated seats
	section := SectionA
	if len(dao.sections[SectionA]) > len(dao.sections[SectionB]) {
		section = SectionB
	}

	// Check if the section is full
	if len(dao.sections[section]) == 25 {
		return "", fmt.Errorf("no available seats in section %s", section)
	}

	seat := dao.availableSeats[section][0]
	dao.availableSeats[section] = slices.Delete(dao.availableSeats[section], 0, 1)

	// Mark the seat as allocated
	dao.sections[section][seat] = nil

	return seat, nil
}

// IsSeatAvailable checks if a seat is available in a section.
func (dao *TrainDAO) isSeatAvailable(seat string) bool {
	section := string(seat[0]) // Extract section from the seat ID (first character)
	exists := slices.Contains(dao.availableSeats[section], seat)
	return exists // If seat does not exist, it is available
}

// AllocateSeat allocates a specific seat to a user if it's available.
func (dao *TrainDAO) ModifySeat(oldSeat, newSeat string, email string) error {
	dao.mu.Lock()
	defer dao.mu.Unlock()

	if available := dao.isSeatAvailable(newSeat); !available {
		return fmt.Errorf("seat %s already booked", newSeat)
	}

	deletedTicket := dao.deallocateSeat(oldSeat)

	dao.newTicket(deletedTicket.From, deletedTicket.To, deletedTicket.User, newSeat)

	return nil
}

// DeallocateSeat deallocates a seat when a user is removed.
func (dao *TrainDAO) deallocateSeat(seat string) *proto.TicketReceipt {
	section := string(seat[0])
	deletedTicket := dao.sections[section][seat]
	delete(dao.sections[section], seat)
	dao.availableSeats[section] = append(dao.availableSeats[section], seat)
	slices.Sort(dao.availableSeats[section])
	return deletedTicket
}

// SaveTicket stores ticket purchase information for a user.
func (dao *TrainDAO) SaveTicket(userDetails *proto.User, from, to string) (*proto.TicketReceipt, error) {
	dao.mu.Lock()
	defer dao.mu.Unlock()

	if user, exists := dao.users[userDetails.Email]; exists {
		return nil, fmt.Errorf("user %s has already booked a ticket", user.Email)
	}

	dao.newUser(userDetails.FirstName, userDetails.LastName, userDetails.Email)

	seat, err := dao.assignSeat()
	if err != nil {
		return nil, err
	}

	ticket := dao.newTicket(from, to, userDetails, seat)
	return ticket, nil
}

// DeleteTicket deletes a user's ticket and deallocates their seat.
func (dao *TrainDAO) DeleteTicket(ticket *proto.TicketReceipt) (*proto.TicketReceipt, error) {
	dao.mu.Lock()
	defer dao.mu.Unlock()

	deletedTicket := dao.deallocateSeat(ticket.Seat)

	return deletedTicket, nil
}

// GetTicket retrieves a user's ticket by their email.
func (dao *TrainDAO) GetTicket(email string) (*proto.TicketReceipt, error) {
	dao.mu.Lock()
	defer dao.mu.Unlock()

	for _, section := range dao.sections {
		for _, ticket := range section {
			if ticket.User.Email == email {
				return ticket, nil
			}
		}
	}
	return nil, fmt.Errorf("ticket for user with email %s not found", email)
}

// GetUsersBySection retrieves all users assigned to seats in a given section.
func (dao *TrainDAO) GetUsersBySection(section string) ([]*proto.TicketReceipt, error) {
	dao.mu.Lock()
	defer dao.mu.Unlock()

	// Validate section
	if section != SectionA && section != SectionB {
		return nil, fmt.Errorf("invalid section: %s", section)
	}

	tickets := []*proto.TicketReceipt{}
	for _, ticket := range dao.sections[section] {
		// Only include users who have a ticket
		tickets = append(tickets, ticket)
	}
	return tickets, nil
}
