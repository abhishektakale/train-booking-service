package dao

import (
	"errors"
	"fmt"
	"sync"
	"train-booking-service/proto"
)

type TrainDAO struct {
	users       map[string]*proto.TicketReceipt
	sections    map[string]map[string]bool // Maps sections to seat availability
	seatCounter map[string]int
	mu          sync.Mutex
}

func NewTrainDAO() *TrainDAO {
	return &TrainDAO{
		users: make(map[string]*proto.TicketReceipt),
		sections: map[string]map[string]bool{
			"A": make(map[string]bool),
			"B": make(map[string]bool),
		},
		seatCounter: make(map[string]int),
	}
}

// Assigns the next available seat in a section.
func (dao *TrainDAO) AssignSeat() (string, error) {
	dao.mu.Lock()
	defer dao.mu.Unlock()

	for _, section := range []string{"A", "B"} {
		dao.seatCounter[section]++
		seat := fmt.Sprintf("%s%s", section ,string(dao.seatCounter[section]))
		if !dao.sections[section][seat] {
			dao.sections[section][seat] = true
			return seat, nil
		}
	}
	return "", errors.New("no available seats")
}

// Checks if a seat is available.
func (dao *TrainDAO) IsSeatAvailable(seat string) bool {
	dao.mu.Lock()
	defer dao.mu.Unlock()

	section := string(seat[0])
	_, exists := dao.sections[section][seat]
	return !exists
}

// Allocates a specific seat to a user.
func (dao *TrainDAO) AllocateSeat(seat string) error {
	dao.mu.Lock()
	defer dao.mu.Unlock()

	section := string(seat[0])
	if dao.sections[section][seat] {
		return errors.New("seat already allocated")
	}
	dao.sections[section][seat] = true
	return nil
}

// Deallocates a seat when a user is removed.
func (dao *TrainDAO) DeallocateSeat(seat string) {
	dao.mu.Lock()
	defer dao.mu.Unlock()

	section := string(seat[0])
	delete(dao.sections[section], seat)
}

// Stores ticket purchase information.
func (dao *TrainDAO) SaveTicket(email string, ticket *proto.TicketReceipt) {
	dao.mu.Lock()
	defer dao.mu.Unlock()

	dao.users[email] = ticket
}

// Deletes a user's ticket.
func (dao *TrainDAO) DeleteTicket(email string) (*proto.TicketReceipt, error) {
	dao.mu.Lock()
	defer dao.mu.Unlock()

	ticket, exists := dao.users[email]
	if !exists {
		return nil, errors.New("user not found")
	}
	delete(dao.users, email)
	dao.DeallocateSeat(ticket.Seat)
	return ticket, nil
}

// Retrieves a user's ticket.
func (dao *TrainDAO) GetTicket(email string) (*proto.TicketReceipt, error) {
	dao.mu.Lock()
	defer dao.mu.Unlock()

	ticket, exists := dao.users[email]
	if !exists {
		return nil, errors.New("user not found")
	}
	return ticket, nil
}

// Gets all users in a section.
func (dao *TrainDAO) GetUsersBySection(section string) []*proto.TicketReceipt {
	dao.mu.Lock()
	defer dao.mu.Unlock()

	users := []*proto.TicketReceipt{}
	for seat := range dao.sections[section] {
		for _, ticket := range dao.users {
			if ticket.Seat == seat {
				users = append(users, ticket)
			}
		}
	}
	return users
}
