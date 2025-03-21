-- Enable UUID generation in PostgreSQL
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-----------------------------------------------------------
-- 1. User Table
-----------------------------------------------------------
CREATE TABLE "User" (
    UserID UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    FirstName VARCHAR(50) NOT NULL,
    LastName VARCHAR(50) NOT NULL,
    Email VARCHAR(100) NOT NULL UNIQUE,
    PhoneNumber VARCHAR(20) NOT NULL UNIQUE,
    City VARCHAR(50),
    HashedPassword VARCHAR(255) NOT NULL,
    RegistrationDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    AccountStatus VARCHAR(10) NOT NULL CHECK (AccountStatus IN ('Active', 'Inactive')),
    UserType VARCHAR(20) NOT NULL CHECK (UserType IN ('Passenger', 'Supporter')),
	CONSTRAINT email_format CHECK (Email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
    CONSTRAINT phone_number_format CHECK (PhoneNumber ~ '^09\d{9}$')
);

-----------------------------------------------------------
-- 2. Ticket Table
-----------------------------------------------------------
CREATE TABLE Ticket (
    TicketID UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    VehicleType VARCHAR(20) NOT NULL CHECK (VehicleType IN ('Airplane', 'Train', 'Bus')),
    Origin VARCHAR(100) NOT NULL,
    Destination VARCHAR(100) NOT NULL,
    DepartureTime TIMESTAMP NOT NULL,
    ArrivalTime TIMESTAMP NOT NULL,
    TicketPrice DECIMAL(10,2) NOT NULL CHECK (TicketPrice >= 0),
    RemainingCapacity INT NOT NULL CHECK (RemainingCapacity >= 0),
    CompanyID UUID,  -- Optional: if referencing a Company table
    TravelClass VARCHAR(20) NOT NULL CHECK (TravelClass IN ('Economy', 'Business', 'VIP')),
	CONSTRAINT check_arrival_after_departure CHECK (ArrivalTime > DepartureTime)
);

-----------------------------------------------------------
-- 3. Reservation Table
-----------------------------------------------------------
CREATE TABLE Reservation (
    ReservationID UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    UserID UUID NOT NULL,
    TicketID UUID NOT NULL,
    ReservationStatus VARCHAR(20) NOT NULL CHECK (ReservationStatus IN ('Reserved', 'Paid', 'Cancelled')),
    ReservationTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ReservationExpirationTime TIMESTAMP,
    FOREIGN KEY (UserID) REFERENCES "User"(UserID) ON DELETE CASCADE,
    FOREIGN KEY (TicketID) REFERENCES Ticket(TicketID) ON DELETE CASCADE
);

-----------------------------------------------------------
-- 4. Payment Table
-----------------------------------------------------------
CREATE TABLE Payment (
    PaymentID UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    UserID UUID NOT NULL,
    ReservationID UUID NOT NULL,
    PaymentAmount DECIMAL(10,2) NOT NULL CHECK (PaymentAmount >= 0),
    PaymentMethod VARCHAR(20) NOT NULL CHECK (PaymentMethod IN ('CreditCard', 'Wallet', 'Crypto')),
    PaymentStatus VARCHAR(20) NOT NULL CHECK (PaymentStatus IN ('Successful', 'Failed', 'Pending')),
    PaymentTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (UserID) REFERENCES "User"(UserID) ON DELETE CASCADE,
    FOREIGN KEY (ReservationID) REFERENCES Reservation(ReservationID) ON DELETE CASCADE
);

-----------------------------------------------------------
-- 5. Reports Table
-----------------------------------------------------------
CREATE TABLE Reports (
    ReportID UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    UserID UUID NOT NULL,
    TicketID UUID,  -- Optional: if the report relates to a specific ticket
    PaymentID  UUID,  -- Optional: if the report relates to a specific payment
    ReportCategory VARCHAR(20) NOT NULL CHECK (ReportCategory IN ('PaymentIssue', 'Delay', 'Cancellation', 'Other')),
    ReportText TEXT,
    ReportStatus VARCHAR(20) DEFAULT 'Pending' CHECK (ReportStatus IN ('Reviewed', 'Pending')),
    ReportTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (UserID) REFERENCES "User"(UserID) ON DELETE CASCADE,
    FOREIGN KEY (TicketID) REFERENCES Ticket(TicketID) ON DELETE SET NULL,
    FOREIGN KEY (PaymentID) REFERENCES Payment(PaymentID) ON DELETE SET NULL
-----------------------------------------------------------
-- 6. TrainDetails Table (Specialized Entity for Train Tickets)
-----------------------------------------------------------
CREATE TABLE TrainDetails (
    TicketID UUID PRIMARY KEY,
    StarRating VARCHAR(1) NOT NULL CHECK (StarRating IN ('3', '4', '5')),
    Amenities VARCHAR(255),  -- List features as a comma-separated string or JSON if supported
    IsCoupeAvailable BOOLEAN NOT NULL,
    FOREIGN KEY (TicketID) REFERENCES Ticket(TicketID) ON DELETE CASCADE
);

-----------------------------------------------------------
-- 7. FlightDetails Table (Specialized Entity for Flight Tickets)
-----------------------------------------------------------
CREATE TABLE FlightDetails (
    TicketID UUID PRIMARY KEY,
    AirlineName VARCHAR(100) NOT NULL,
    FlightClass VARCHAR(20) NOT NULL CHECK (FlightClass IN ('Economy', 'Business', 'First')),
    NumberOfStops INT DEFAULT 0 CHECK (NumberOfStops >= 0),
    FlightNumber VARCHAR(50),
    OriginAirport VARCHAR(100),
    DestinationAirport VARCHAR(100),
    Amenities VARCHAR(255),
    FOREIGN KEY (TicketID) REFERENCES Ticket(TicketID) ON DELETE CASCADE
);

-----------------------------------------------------------
-- 8. BusDetails Table (Specialized Entity for Bus Tickets)
-----------------------------------------------------------
CREATE TABLE BusDetails (
    TicketID UUID PRIMARY KEY,
    BusCompanyName VARCHAR(100) NOT NULL,
    BusType VARCHAR(20) NOT NULL CHECK (BusType IN ('VIP', 'Regular', 'Sleeper')),
    SeatsPerRow VARCHAR(10),  -- e.g., '1+2' or '2+2'
    Amenities VARCHAR(255),
    FOREIGN KEY (TicketID) REFERENCES Ticket(TicketID) ON DELETE CASCADE
);