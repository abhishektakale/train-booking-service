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

    Run client to test:
    ```bash
        make run-client
    ```

---

## Usage

1. **Purchase a ticket:**
   Use the `PurchaseTicket API` to buy a ticket for a user. You'll need to provide their first name, last name and email.

2. **View receipt:**
   After purchasing the ticket, you can retrieve the receipt details through the `GetReceipt API` using the user's email.

3. **View users and seats by section:**
   To see a list of users and their allocated seats in a given section (A or B), use the `GetUsersBySection API`.

4. **Remove a user:**
   If you need to remove a user from the train, use the `RemoveUser API` by specifying the user's email.

5. **Modify a user's seat:**
   To change a user's seat, use the `ModifySeat API` and provide the desired seat.

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