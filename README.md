# Train Booking Service

## Description

This project implements a **Train Booking Service** that allows users to purchase train tickets. All APIs are gRPC-based. The service manages seat allocation, ticket purchasing, and provides functionality to modify seat or remove users from the train.

---

## Installation Instructions

1. **Clone the repository**
   ```bash
    git clone https://github.com/abhishektakale/train-booking-service.git
    cd train-booking-service
   ```

2. **Install dependencies**
    
    Ensure you have `Go` installed, then install necessary packages:
    ```bash
        go mod tidy
    ```

3. **Run the service**

    Start the gRPC server that exposes all the APIs:
    ```bash
        make run-sever
    ```

---

## Usage

1. **Purchase a ticket:**
   ```bash
   $ go run cmd/client/main.go -Operation="PurchaseTicket" -Data='{"from": "London", "to": "France", "user": {"first_name": "John", "last_name": "Doe", "email": "johndoe@example.com"}}'
   ```
   Output:
   ```
   Ticket purchased: Ticket purchased successfully, Seat: A1
   ```
2. **View receipt:**
   ```bash
   $ go run cmd/client/main.go -Operation="GetReceipt" -Data='{"user_email": "johndoe@example.com"}'
   ```
   Output:
   ```
   Receipt: London to France, Seat: A1, Price Paid: 20.00
   ```
3. **View users and seats by section:**
   ```bash
   go run cmd/client/main.go -Operation="GetUsersBySection" -Data='{"section": "A"}'
   ```
   Output:
   ```
   User List:

   User: John Doe (johndoe@example.com)
   Seat: A1
   ```
4. **Modify a user's seat:**
   ```bash
   go run cmd/client/main.go -Operation="ModifySeat" -Data='{"user_email": "johndoe@example.com", "new_seat": "A2"}'
   ```
   Output:
   ```
   New Receipt: London to France, Seat: A2, Price Paid: 20.00
   ```
5. **Remove a user:**
   ```bash
   go run cmd/client/main.go -Operation="RemoveUser" -Data='{"user_email": "johndoe@example.com"}'
   ```
   Output:
   ```
   Removed User: Name: John Doe, Email: johndoe@example.com
   ```

---

## APIs

### 1. **PurchaseTicket API**

**Description:** Submit a purchase request for a train ticket to be booked.

**Fields:**

- `From`: Source Station
- `To`: Destination Station
- `User`: User Details (First Name, Last Name, Email)
    - `First Name`: User's first name
    - `Last Name`: User's last name
    - `Email`: User's email address

**Response:**

- Confirms ticket purchase
- Allocates a seat in the train (either Section A or Section B)

- **Details on the Receipt:**
    ```
    From: Source
    To: Destination
    User: First Name, Last Name, Email Address
    Price Paid: Price
    Seat: Seat Number
    ```

---

### 2. **GetReceipt API**

**Description:** Fetches the details of a receipt for the user.  
 **Fields:**

- `Email`: Email address of the user

**Response:**

- Displays the details of the ticket including user name, price paid, and allocated seat.

---

### 3. **GetUsersbySection API**

**Description:** Allows the user to view the seat allocation of users in a specified train section (Section A or Section B).  
 **Fields:**

- `Section`: The section (A or B) to retrieve seat allocations

**Response:**

- List of users and the seat they are allocated in the specified section.

---

### 4. **RemoveUser API**

**Description:** Removes a user from the train system.  
**Fields:**

- `Email`: The email address of the user to be removed

**Response:**

- Confirms that the user has been removed from the system and no longer has a seat allocation.

---

### 5. **ModifySeat API**

**Description:** Allows modification of a user's seat assignment.  
 **Fields:**

- `Email`: The email address of the user whose seat is to be modified
- `New Seat`: The new seat that the user has requested

**Response:**

- Confirms the seat modification and provides the updated seat information.

---

## Ticket Receipt Sample

```bash
    - From: London
    - To: France
    - Ticket Price: 20.00
    - Seat: A1
    - User: 
        - First Name: John
        - Last Name: Doe
        - Email: johndoe@example.com
```