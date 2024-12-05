package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"train-booking-service/proto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock TrainDAOAdapter
type MockTrainDAOAdapter struct {
	mock.Mock
}

func (m *MockTrainDAOAdapter) SaveTicket(user *proto.User, from, to string) (*proto.TicketReceipt, error) {
	args := m.Called(user, from, to)
	if ticket, ok := args.Get(0).(*proto.TicketReceipt); ok {
		return ticket, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTrainDAOAdapter) GetTicket(userEmail string) (*proto.TicketReceipt, error) {
	args := m.Called(userEmail)
	if ticket, ok := args.Get(0).(*proto.TicketReceipt); ok {
		return ticket, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTrainDAOAdapter) ModifySeat(oldSeat, newSeat, userEmail string) error {
	args := m.Called(oldSeat, newSeat, userEmail)
	return args.Error(0)
}

func (m *MockTrainDAOAdapter) GetUsersBySection(section string) ([]*proto.TicketReceipt, error) {
	args := m.Called(section)
	if ticket, ok := args.Get(0).([]*proto.TicketReceipt); ok {
		return ticket, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTrainDAOAdapter) DeleteTicket(ticket *proto.TicketReceipt) (*proto.TicketReceipt, error) {
	args := m.Called(ticket)
	if ticket, ok := args.Get(0).(*proto.TicketReceipt); ok {
		return ticket, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTrainDAOAdapter) AvailableSeats(section string) []string {
	args := m.Called(section)
	return args.Get(0).([]string)
}

// Test cases

func TestPurchaseTicket_Success(t *testing.T) {
	// Arrange
	mockDAO := new(MockTrainDAOAdapter)
	service := NewTrainService(mockDAO)
	user := &proto.User{Email: "user@example.com"}
	ticket := &proto.TicketReceipt{User: user, From: "StationA", To: "StationB", Seat: "A1"}

	// Mock the SaveTicket method to return the ticket and no error
	mockDAO.On("SaveTicket", user, "StationA", "StationB").Return(ticket, nil)

	// Act
	req := &proto.PurchaseTicketRequest{
		User: user,
		From: "StationA",
		To:   "StationB",
	}
	resp, err := service.PurchaseTicket(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, ticket, resp.Ticket)
	assert.Equal(t, "Ticket purchased successfully", resp.Message)

	// Ensure all expectations were met
	mockDAO.AssertExpectations(t)
}

func TestPurchaseTicket_Error(t *testing.T) {
	// Arrange
	mockDAO := new(MockTrainDAOAdapter)
	service := NewTrainService(mockDAO)
	user := &proto.User{Email: "user@example.com"}

	// Mock the SaveTicket method to return an error
	mockDAO.On("SaveTicket", user, "StationA", "StationB").Return(nil, errors.New("db error"))

	// Act
	req := &proto.PurchaseTicketRequest{
		User: user,
		From: "StationA",
		To:   "StationB",
	}
	resp, err := service.PurchaseTicket(context.Background(), req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "could not save ticket")

	// Ensure all expectations were met
	mockDAO.AssertExpectations(t)
}

func TestGetReceipt_Success(t *testing.T) {
	// Arrange
	mockDAO := new(MockTrainDAOAdapter)
	service := NewTrainService(mockDAO)
	userEmail := "user@example.com"
	ticket := &proto.TicketReceipt{User: &proto.User{Email: userEmail}, From: "StationA", To: "StationB", Seat: "A1"}

	// Mock the GetTicket method to return the ticket and no error
	mockDAO.On("GetTicket", userEmail).Return(ticket, nil)

	// Act
	req := &proto.GetReceiptRequest{UserEmail: userEmail}
	resp, err := service.GetReceipt(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, ticket, resp.Ticket)

	// Ensure all expectations were met
	mockDAO.AssertExpectations(t)
}

func TestGetReceipt_Error(t *testing.T) {
	// Arrange
	mockDAO := new(MockTrainDAOAdapter)
	service := NewTrainService(mockDAO)
	userEmail := "user@example.com"

	// Mock the GetTicket method to return an error
	mockDAO.On("GetTicket", userEmail).Return(nil, errors.New("ticket not found"))

	// Act
	req := &proto.GetReceiptRequest{UserEmail: userEmail}
	resp, err := service.GetReceipt(context.Background(), req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "could not retrieve ticket")

	// Ensure all expectations were met
	mockDAO.AssertExpectations(t)
}

func TestModifySeat_Success(t *testing.T) {
	// Arrange
	mockDAO := new(MockTrainDAOAdapter)
	service := NewTrainService(mockDAO)
	userEmail := "user@example.com"
	oldTicket := &proto.TicketReceipt{User: &proto.User{Email: userEmail}, Seat: "A1"}
	newSeat := "B2"

	// Mock the GetTicket and ModifySeat methods
	mockDAO.On("GetTicket", userEmail).Return(oldTicket, nil)
	mockDAO.On("ModifySeat", "A1", newSeat, userEmail).Return(nil)

	// Act
	req := &proto.ModifySeatRequest{UserEmail: userEmail, NewSeat: newSeat}
	resp, err := service.ModifySeat(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, newSeat, resp.NewTicket.Seat)
	assert.Equal(t, "Seat modified successfully", resp.Message)

	// Ensure all expectations were met
	mockDAO.AssertExpectations(t)
}

func TestModifySeat_Error(t *testing.T) {
	// Arrange
	mockDAO := new(MockTrainDAOAdapter)
	service := NewTrainService(mockDAO)
	userEmail := "user@example.com"
	newSeat := "B2"

	// Mock GetTicket method to return an error
	mockDAO.On("GetTicket", userEmail).Return(nil, errors.New("ticket not found"))

	// Act
	req := &proto.ModifySeatRequest{UserEmail: userEmail, NewSeat: newSeat}
	resp, err := service.ModifySeat(context.Background(), req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "could not retrieve ticket")

	// Ensure all expectations were met
	mockDAO.AssertExpectations(t)
}

func TestModifySeat_Error_OwnedByAnotherUser(t *testing.T) {
	// Arrange
	mockDAO := new(MockTrainDAOAdapter)
	service := NewTrainService(mockDAO)
	userEmail2 := "user1@example.com"
	ticket2 := &proto.TicketReceipt{User: &proto.User{Email: userEmail2}, Seat: "A2"}

	// Mock GetTicket for the second user (user1@example.com), it will simulate that the user does not own the seat
	mockDAO.On("GetTicket", userEmail2).Return(ticket2, nil)

	// Mock the ModifySeat method to return an error because the seat is owned by another user
	mockDAO.On("ModifySeat", "A2", "A1", userEmail2).Return(errors.New("seat already owned by another user"))

	// Act
	req := &proto.ModifySeatRequest{UserEmail: userEmail2, NewSeat: "A1"}
	resp, err := service.ModifySeat(context.Background(), req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "seat already owned by another user")

	// Ensure all expectations were met
	mockDAO.AssertExpectations(t)
}

func TestRemoveUser_Success(t *testing.T) {
	// Arrange
	mockDAO := new(MockTrainDAOAdapter)
	service := NewTrainService(mockDAO)
	userEmail := "user@example.com"
	ticket := &proto.TicketReceipt{User: &proto.User{Email: userEmail}}

	// Mock the GetTicket and DeleteTicket methods
	mockDAO.On("GetTicket", userEmail).Return(ticket, nil)
	mockDAO.On("DeleteTicket", ticket).Return(ticket, nil)

	// Act
	req := &proto.RemoveUserRequest{UserEmail: userEmail}
	resp, err := service.RemoveUser(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "User removed successfully", resp.Message)

	// Ensure all expectations were met
	mockDAO.AssertExpectations(t)
}

func TestRemoveUser_Error(t *testing.T) {
	// Arrange
	mockDAO := new(MockTrainDAOAdapter)
	service := NewTrainService(mockDAO)
	userEmail := "user@example.com"

	// Mock GetTicket method to return an error
	mockDAO.On("GetTicket", userEmail).Return(nil, errors.New("ticket not found"))

	// Act
	req := &proto.RemoveUserRequest{UserEmail: userEmail}
	resp, err := service.RemoveUser(context.Background(), req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "could not retrieve ticket")

	// Ensure all expectations were met
	mockDAO.AssertExpectations(t)
}

func TestRemoveUser_Error_OwnedByAnotherUser(t *testing.T) {
	// Arrange
	mockDAO := new(MockTrainDAOAdapter)
	service := NewTrainService(mockDAO)
	userEmail := "user@example.com"
	ticket := &proto.TicketReceipt{User: &proto.User{Email: userEmail}, Seat: "A1"}

	// Mock DeleteTicket method to return an error
	mockDAO.On("GetTicket", userEmail).Return(ticket, nil)
	mockDAO.On("DeleteTicket", ticket).Return(nil, errors.New("seat already owned by another user"))

	// Act
	req := &proto.RemoveUserRequest{UserEmail: userEmail}
	resp, err := service.RemoveUser(context.Background(), req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "seat already owned by another user")

	// Ensure all expectations were met
	mockDAO.AssertExpectations(t)
}

func TestGetUsersBySection_Success(t *testing.T) {
	// Arrange
	mockDAO := new(MockTrainDAOAdapter)
	service := NewTrainService(mockDAO)
	section := "A1"
	ticket1 := &proto.TicketReceipt{User: &proto.User{Email: "user1@example.com"}, Seat: "1A"}
	ticket2 := &proto.TicketReceipt{User: &proto.User{Email: "user2@example.com"}, Seat: "1B"}

	// Mock GetUsersBySection to return the tickets for the given section
	mockDAO.On("GetUsersBySection", section).Return([]*proto.TicketReceipt{ticket1, ticket2}, nil)

	// Act
	req := &proto.GetUsersBySectionRequest{Section: section}
	resp, err := service.GetUsersBySection(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.UserSeats, 2)
	assert.Equal(t, "user1@example.com", resp.UserSeats[0].User.Email)
	assert.Equal(t, "1A", resp.UserSeats[0].Seat)
	assert.Equal(t, "user2@example.com", resp.UserSeats[1].User.Email)
	assert.Equal(t, "1B", resp.UserSeats[1].Seat)

	// Ensure all expectations were met
	mockDAO.AssertExpectations(t)
}

func TestGetUsersBySection_Error_InvalidSection(t *testing.T) {
	// Arrange
	mockDAO := new(MockTrainDAOAdapter)
	service := NewTrainService(mockDAO)
	invalidSection := ""

	// Act
	req := &proto.GetUsersBySectionRequest{Section: invalidSection}
	resp, err := service.GetUsersBySection(context.Background(), req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "invalid section")

	// Ensure all expectations were met
	mockDAO.AssertExpectations(t)
}

func TestGetUsersBySection_Error_FetchingTickets(t *testing.T) {
	// Arrange
	mockDAO := new(MockTrainDAOAdapter)
	service := NewTrainService(mockDAO)
	section := "A1"

	// Mock GetUsersBySection to return an error
	mockDAO.On("GetUsersBySection", section).Return(nil, fmt.Errorf("db error"))

	// Act
	req := &proto.GetUsersBySectionRequest{Section: section}
	resp, err := service.GetUsersBySection(context.Background(), req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "could not fetch users by section")

	// Ensure all expectations were met
	mockDAO.AssertExpectations(t)
}
