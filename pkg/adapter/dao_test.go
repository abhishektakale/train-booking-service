package adapter

import (
	"fmt"
	"testing"
	"train-booking-service/proto"

	"github.com/stretchr/testify/assert"
)

func TestNewTrainDAO(t *testing.T) {
	dao := NewTrainDAO()

	assert.NotNil(t, dao)
	assert.Len(t, dao.AvailableSeats(SectionA), SectionCap)
	assert.Len(t, dao.AvailableSeats(SectionB), SectionCap)
}

func TestPurchaseTicketSuccessful(t *testing.T) {
	dao := NewTrainDAO()
	ticket, err := dao.SaveTicket(&proto.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@example.com",
	}, "London", "France")
	assert.NoError(t, err)
	assert.NotEmpty(t, ticket)
	assert.Equal(t, "John", ticket.User.FirstName)
	assert.Equal(t, "Doe", ticket.User.LastName)
	assert.Equal(t, "johndoe@example.com", ticket.User.Email)
	assert.Equal(t, "London", ticket.From)
	assert.Equal(t, "France", ticket.To)
	assert.Equal(t, float32(20), ticket.PricePaid)
	assert.Equal(t, "A1", ticket.Seat)
}

func TestPurchaseTicketFailure_SameUserBookingTwoTickets(t *testing.T) {
	dao := NewTrainDAO()
	ticket1, err := dao.SaveTicket(&proto.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@example.com",
	}, "London", "France")
	assert.NoError(t, err)
	assert.NotEmpty(t, ticket1)
	assert.Equal(t, "John", ticket1.User.FirstName)
	assert.Equal(t, "Doe", ticket1.User.LastName)
	assert.Equal(t, "johndoe@example.com", ticket1.User.Email)
	assert.Equal(t, "London", ticket1.From)
	assert.Equal(t, "France", ticket1.To)
	assert.Equal(t, float32(20), ticket1.PricePaid)
	assert.Equal(t, "A1", ticket1.Seat)
	_, err = dao.SaveTicket(&proto.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@example.com",
	}, "London", "France")
	assert.Error(t, err)
}

func TestPurchaseTicketFailure_SeatsFull(t *testing.T) {
	dao := NewTrainDAO()
	for i := 0; i < 50; i++ {
		_, err := dao.SaveTicket(&proto.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     fmt.Sprintf("johndoe%v@example.com", i),
		}, "London", "France")
		assert.NoError(t, err)
	}
	_, err := dao.SaveTicket(&proto.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@example.com",
	}, "London", "France")
	assert.Error(t, err)
}

func TestModifySeat_SuccessfulAllocation(t *testing.T) {
	dao := NewTrainDAO()
	oldTicket, err := dao.SaveTicket(&proto.User{
		FirstName: "Alice",
		LastName:  "Doe",
		Email:     "alicedoe@example.com",
	}, "London", "France")
	assert.NoError(t, err)
	assert.NotEmpty(t, oldTicket)
	assert.Equal(t, "Alice", oldTicket.User.FirstName)
	assert.Equal(t, "Doe", oldTicket.User.LastName)
	assert.Equal(t, "alicedoe@example.com", oldTicket.User.Email)
	assert.Equal(t, "London", oldTicket.From)
	assert.Equal(t, "France", oldTicket.To)
	assert.Equal(t, float32(20), oldTicket.PricePaid)
	assert.Equal(t, "A1", oldTicket.Seat)

	err = dao.ModifySeat(oldTicket.Seat, "B1", oldTicket.User.Email)
	assert.NoError(t, err)

	newTicket, err := dao.GetTicket(oldTicket.User.Email)
	assert.NoError(t, err)
	assert.Equal(t, "B1", newTicket.Seat)
}

func TestModifySeat_AlreadyAllocated(t *testing.T) {
	dao := NewTrainDAO()
	ticket1, err := dao.SaveTicket(&proto.User{
		FirstName: "Alice",
		LastName:  "Doe",
		Email:     "alicedoe@example.com",
	}, "London", "France")
	assert.NoError(t, err)
	assert.NotEmpty(t, ticket1)
	assert.Equal(t, "Alice", ticket1.User.FirstName)
	assert.Equal(t, "Doe", ticket1.User.LastName)
	assert.Equal(t, "alicedoe@example.com", ticket1.User.Email)
	assert.Equal(t, "London", ticket1.From)
	assert.Equal(t, "France", ticket1.To)
	assert.Equal(t, float32(20), ticket1.PricePaid)
	assert.Equal(t, "A1", ticket1.Seat)

	ticket2, err := dao.SaveTicket(&proto.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@example.com",
	}, "London", "France")
	assert.NoError(t, err)
	assert.NotEmpty(t, ticket2)
	assert.Equal(t, "John", ticket2.User.FirstName)
	assert.Equal(t, "Doe", ticket2.User.LastName)
	assert.Equal(t, "johndoe@example.com", ticket2.User.Email)
	assert.Equal(t, "London", ticket2.From)
	assert.Equal(t, "France", ticket2.To)
	assert.Equal(t, float32(20), ticket2.PricePaid)
	assert.Equal(t, "B1", ticket2.Seat)

	err = dao.ModifySeat(ticket1.Seat, "B1", ticket1.User.Email)
	assert.Error(t, err)
}

func TestDeleteTicket(t *testing.T) {
	dao := NewTrainDAO()
	ticket, err := dao.SaveTicket(&proto.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@example.com",
	}, "London", "France")
	assert.NoError(t, err)
	assert.NotEmpty(t, ticket)
	assert.Equal(t, "John", ticket.User.FirstName)
	assert.Equal(t, "Doe", ticket.User.LastName)
	assert.Equal(t, "johndoe@example.com", ticket.User.Email)
	assert.Equal(t, "London", ticket.From)
	assert.Equal(t, "France", ticket.To)
	assert.Equal(t, float32(20), ticket.PricePaid)
	assert.Equal(t, "A1", ticket.Seat)

	deletedTicket, err := dao.DeleteTicket(ticket)
	assert.NoError(t, err)
	assert.Equal(t, ticket.Seat, deletedTicket.Seat)

	_, err = dao.GetTicket(deletedTicket.User.Email)
	assert.Error(t, err)
}

func TestGetUsersBySection(t *testing.T) {
	dao := NewTrainDAO()
	ticket1, err := dao.SaveTicket(&proto.User{
		FirstName: "Alice",
		LastName:  "Doe",
		Email:     "alicedoe@example.com",
	}, "London", "France")
	assert.NoError(t, err)
	assert.NotEmpty(t, ticket1)
	assert.Equal(t, "Alice", ticket1.User.FirstName)
	assert.Equal(t, "Doe", ticket1.User.LastName)
	assert.Equal(t, "alicedoe@example.com", ticket1.User.Email)
	assert.Equal(t, "London", ticket1.From)
	assert.Equal(t, "France", ticket1.To)
	assert.Equal(t, float32(20), ticket1.PricePaid)
	assert.Equal(t, "A1", ticket1.Seat)

	ticket2, err := dao.SaveTicket(&proto.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@example.com",
	}, "London", "France")
	assert.NoError(t, err)
	assert.NotEmpty(t, ticket2)
	assert.Equal(t, "John", ticket2.User.FirstName)
	assert.Equal(t, "Doe", ticket2.User.LastName)
	assert.Equal(t, "johndoe@example.com", ticket2.User.Email)
	assert.Equal(t, "London", ticket2.From)
	assert.Equal(t, "France", ticket2.To)
	assert.Equal(t, float32(20), ticket2.PricePaid)
	assert.Equal(t, "B1", ticket2.Seat)

	ticket3, err := dao.SaveTicket(&proto.User{
		FirstName: "Bob",
		LastName:  "Doe",
		Email:     "bobdoe@example.com",
	}, "London", "France")
	assert.NoError(t, err)
	assert.NotEmpty(t, ticket3)
	assert.Equal(t, "Bob", ticket3.User.FirstName)
	assert.Equal(t, "Doe", ticket3.User.LastName)
	assert.Equal(t, "bobdoe@example.com", ticket3.User.Email)
	assert.Equal(t, "London", ticket3.From)
	assert.Equal(t, "France", ticket3.To)
	assert.Equal(t, float32(20), ticket3.PricePaid)
	assert.Equal(t, "A2", ticket3.Seat)

	tickets, err := dao.GetUsersBySection(SectionA)
	assert.NoError(t, err)
	assert.Len(t, tickets, 2)

	assert.Equal(t, ticket1.User.Email, tickets[0].User.Email)
	assert.Equal(t, ticket3.User.Email, tickets[1].User.Email)
}

func TestGetUsersBySection_InvalidSection(t *testing.T) {
	dao := NewTrainDAO()
	ticket1, err := dao.SaveTicket(&proto.User{
		FirstName: "Alice",
		LastName:  "Doe",
		Email:     "alicedoe@example.com",
	}, "London", "France")
	assert.NoError(t, err)
	assert.NotEmpty(t, ticket1)
	assert.Equal(t, "Alice", ticket1.User.FirstName)
	assert.Equal(t, "Doe", ticket1.User.LastName)
	assert.Equal(t, "alicedoe@example.com", ticket1.User.Email)
	assert.Equal(t, "London", ticket1.From)
	assert.Equal(t, "France", ticket1.To)
	assert.Equal(t, float32(20), ticket1.PricePaid)
	assert.Equal(t, "A1", ticket1.Seat)

	ticket2, err := dao.SaveTicket(&proto.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@example.com",
	}, "London", "France")
	assert.NoError(t, err)
	assert.NotEmpty(t, ticket2)
	assert.Equal(t, "John", ticket2.User.FirstName)
	assert.Equal(t, "Doe", ticket2.User.LastName)
	assert.Equal(t, "johndoe@example.com", ticket2.User.Email)
	assert.Equal(t, "London", ticket2.From)
	assert.Equal(t, "France", ticket2.To)
	assert.Equal(t, float32(20), ticket2.PricePaid)
	assert.Equal(t, "B1", ticket2.Seat)

	ticket3, err := dao.SaveTicket(&proto.User{
		FirstName: "Bob",
		LastName:  "Doe",
		Email:     "bobdoe@example.com",
	}, "London", "France")
	assert.NoError(t, err)
	assert.NotEmpty(t, ticket3)
	assert.Equal(t, "Bob", ticket3.User.FirstName)
	assert.Equal(t, "Doe", ticket3.User.LastName)
	assert.Equal(t, "bobdoe@example.com", ticket3.User.Email)
	assert.Equal(t, "London", ticket3.From)
	assert.Equal(t, "France", ticket3.To)
	assert.Equal(t, float32(20), ticket3.PricePaid)
	assert.Equal(t, "A2", ticket3.Seat)

	_, err = dao.GetUsersBySection("C")
	assert.Error(t, err)
}
