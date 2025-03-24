package config

// import (
// 	"database/sql"
// 	"fmt"

// 	"github.com/google/uuid"
// 	"golang.org/x/crypto/bcrypt"
// )

// func init() {
// 	Connect()
// 	db = GetDB()
// 	query := `
// 	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

// 	CREATE TABLE "User" (
// 		UserID UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
// 		FirstName VARCHAR(50) NOT NULL,
// 		LastName VARCHAR(50) NOT NULL,
// 		Email VARCHAR(100) NOT NULL UNIQUE,
// 		PhoneNumber VARCHAR(20) NOT NULL UNIQUE,
// 		City VARCHAR(50),
// 		HashedPassword VARCHAR(255) NOT NULL,
// 		RegistrationDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
// 		AccountStatus VARCHAR(10) NOT NULL CHECK (AccountStatus IN ('Active', 'Inactive')),
// 		UserType VARCHAR(20) NOT NULL CHECK (UserType IN ('Passenger', 'Supporter')),
// 		CONSTRAINT email_format CHECK (Email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
// 		CONSTRAINT phone_number_format CHECK (PhoneNumber ~ '^09\d{9}$')
// 	);

// 	CREATE TABLE IF NOT EXISTS Ticket (
// 	    TicketID UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
// 	    VehicleType VARCHAR(20) NOT NULL CHECK (VehicleType IN ('Airplane', 'Train', 'Bus')),
// 	    Origin VARCHAR(100) NOT NULL,
// 	    Destination VARCHAR(100) NOT NULL,
// 	    DepartureTime TIMESTAMP NOT NULL,
// 	    ArrivalTime TIMESTAMP NOT NULL,
// 	    TicketPrice DECIMAL(10,2) NOT NULL CHECK (TicketPrice >= 0),
// 	    RemainingCapacity INT NOT NULL CHECK (RemainingCapacity >= 0),
// 	    CompanyID UUID,
// 	    TravelClass VARCHAR(20) NOT NULL CHECK (TravelClass IN ('Economy', 'Business', 'VIP')),
// 	    CONSTRAINT check_arrival_after_departure CHECK (ArrivalTime > DepartureTime)
// 	);
// 	CREATE TABLE Reservation (
// 		ReservationID UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
// 		UserID UUID NOT NULL,
// 		TicketID UUID NOT NULL,
// 		ReservationStatus VARCHAR(20) NOT NULL CHECK (ReservationStatus IN ('Reserved', 'Paid', 'Cancelled')),
// 		ReservationTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
// 		ReservationExpirationTime TIMESTAMP,
// 		FOREIGN KEY (UserID) REFERENCES "User"(UserID) ON DELETE CASCADE,
// 		FOREIGN KEY (TicketID) REFERENCES Ticket(TicketID) ON DELETE CASCADE
// 	);
// 	CREATE TABLE Payment (
// 		PaymentID UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
// 		UserID UUID NOT NULL,
// 		ReservationID UUID NOT NULL,
// 		PaymentAmount DECIMAL(10,2) NOT NULL CHECK (PaymentAmount >= 0),
// 		PaymentMethod VARCHAR(20) NOT NULL CHECK (PaymentMethod IN ('CreditCard', 'Wallet', 'Crypto')),
// 		PaymentStatus VARCHAR(20) NOT NULL CHECK (PaymentStatus IN ('Successful', 'Failed', 'Pending')),
// 		PaymentTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
// 		FOREIGN KEY (UserID) REFERENCES "User"(UserID) ON DELETE CASCADE,
// 		FOREIGN KEY (ReservationID) REFERENCES Reservation(ReservationID) ON DELETE CASCADE
// 	);
// 	CREATE TABLE Reports (
// 		ReportID UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
// 		UserID UUID NOT NULL,
// 		TicketID UUID,  -- Optional: if the report relates to a specific ticket
// 		PaymentID  UUID,  -- Optional: if the report relates to a specific payment
// 		ReportCategory VARCHAR(20) NOT NULL CHECK (ReportCategory IN ('PaymentIssue', 'Delay', 'Cancellation', 'Other')),
// 		ReportText TEXT,
// 		ReportStatus VARCHAR(20) DEFAULT 'Pending' CHECK (ReportStatus IN ('Reviewed', 'Pending')),
// 		ReportTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
// 		FOREIGN KEY (UserID) REFERENCES "User"(UserID) ON DELETE CASCADE,
// 		FOREIGN KEY (TicketID) REFERENCES Ticket(TicketID) ON DELETE SET NULL,
// 		FOREIGN KEY (PaymentID) REFERENCES Payment(PaymentID) ON DELETE SET NULL
// 	);
// 	CREATE TABLE TrainDetails (
// 		TicketID UUID PRIMARY KEY,
// 		StarRating VARCHAR(1) NOT NULL CHECK (StarRating IN ('3', '4', '5')),
// 		Amenities VARCHAR(255),  -- List features as a comma-separated string or JSON if supported
// 		IsCoupeAvailable BOOLEAN NOT NULL,
// 		FOREIGN KEY (TicketID) REFERENCES Ticket(TicketID) ON DELETE CASCADE
// 	);
// 	CREATE TABLE FlightDetails (
// 		TicketID UUID PRIMARY KEY,
// 		AirlineName VARCHAR(100) NOT NULL,
// 		FlightClass VARCHAR(20) NOT NULL CHECK (FlightClass IN ('Economy', 'Business', 'First')),
// 		NumberOfStops INT DEFAULT 0 CHECK (NumberOfStops >= 0),
// 		FlightNumber VARCHAR(50),
// 		OriginAirport VARCHAR(100),
// 		DestinationAirport VARCHAR(100),
// 		Amenities VARCHAR(255),
// 		FOREIGN KEY (TicketID) REFERENCES Ticket(TicketID) ON DELETE CASCADE
// 	);
// 	CREATE TABLE BusDetails (
// 		TicketID UUID PRIMARY KEY,
// 		BusCompanyName VARCHAR(100) NOT NULL,
// 		BusType VARCHAR(20) NOT NULL CHECK (BusType IN ('VIP', 'Regular', 'Sleeper')),
// 		SeatsPerRow VARCHAR(10),  -- e.g., '1+2' or '2+2'
// 		Amenities VARCHAR(255),
// 		FOREIGN KEY (TicketID) REFERENCES Ticket(TicketID) ON DELETE CASCADE
// 	);

// 	-- Insert Users
// 	INSERT INTO "User" (FirstName, LastName, Email, PhoneNumber, City, HashedPassword, AccountStatus, UserType)
// 	VALUES
// 	('Ali', 'Ahmadi', 'ali.ahmadi@example.com', '09123456789', 'Tehran', 'hashed_pass_1', 'Active', 'Passenger'),
// 	('Sara', 'Mohammadi', 'sara.mohammadi@example.com', '09234567890', 'Mashhad', 'hashed_pass_2', 'Active', 'Supporter'),
// 	('Reza', 'Hosseini', 'reza.hosseini@example.com', '09345678901', 'Isfahan', 'hashed_pass_3', 'Inactive', 'Passenger'),
// 	('Fatemeh', 'Karimi', 'fatemeh.karimi@example.com', '09456789012', 'Shiraz', 'hashed_pass_4', 'Active', 'Passenger'),
// 	('Hossein', 'Ranjbar', 'hossein.ranjbar@example.com', '09567890123', 'Tabriz', 'hashed_pass_5', 'Active', 'Supporter'),
// 	('Mehdi', 'Zarei', 'mehdi.zarei@example.com', '09678901234', 'Ahvaz', 'hashed_pass_6', 'Inactive', 'Passenger'),
// 	('Maryam', 'Jafari', 'maryam.jafari@example.com', '09789012345', 'Qom', 'hashed_pass_7', 'Active', 'Passenger'),
// 	('Ehsan', 'Shirazi', 'ehsan.shirazi@example.com', '09890123456', 'Kerman', 'hashed_pass_8', 'Active', 'Supporter'),
// 	('Zahra', 'Mansouri', 'zahra.mansouri@example.com', '09901234567', 'Rasht', 'hashed_pass_9', 'Inactive', 'Passenger'),
// 	('Hassan', 'Gholami', 'hassan.gholami@example.com', '09112345678', 'Urmia', 'hashed_pass_10', 'Active', 'Passenger');
// 	-- Insert Tickets
// 	INSERT INTO Ticket (VehicleType, Origin, Destination, DepartureTime, ArrivalTime, TicketPrice, RemainingCapacity, TravelClass)
// 	VALUES
// 	('Airplane', 'Tehran', 'Mashhad', '2025-04-01 08:00', '2025-04-01 09:30', 2500000, 50, 'Economy'),
// 	('Train', 'Isfahan', 'Tehran', '2025-04-02 10:00', '2025-04-02 16:00', 1200000, 100, 'Business'),
// 	('Bus', 'Shiraz', 'Tabriz', '2025-04-03 09:00', '2025-04-04 09:00', 800000, 40, 'VIP'),
// 	('Airplane', 'Ahvaz', 'Kerman', '2025-04-05 07:00', '2025-04-05 08:00', 1800000, 60, 'Business'),
// 	('Train', 'Qom', 'Mashhad', '2025-04-06 06:00', '2025-04-06 14:00', 900000, 70, 'Economy'),
// 	('Bus', 'Tehran', 'Ahvaz', '2025-04-07 20:00', '2025-04-08 08:00', 1500000, 35, 'Economy'),
// 	('Airplane', 'Rasht', 'Shiraz', '2025-04-08 15:00', '2025-04-08 16:30', 2000000, 45, 'Business'),
// 	('Train', 'Tabriz', 'Tehran', '2025-04-09 11:00', '2025-04-09 17:00', 1100000, 80, 'Business'),
// 	('Bus', 'Kerman', 'Isfahan', '2025-04-10 08:00', '2025-04-10 20:00', 750000, 30, 'Economy'),
// 	('Airplane', 'Mashhad', 'Tehran', '2025-04-11 14:00', '2025-04-11 15:30', 3000000, 55, 'Economy');

// 	-- Insert Reservations
// 	INSERT INTO Reservation (UserID, TicketID, ReservationStatus)
// 	SELECT u.UserID, t.TicketID, 'Reserved'
// 	FROM "User" u, Ticket t
// 	ORDER BY random()
// 	LIMIT 10;

// 	-- Insert Payments
// 	INSERT INTO Payment (UserID, ReservationID, PaymentAmount, PaymentMethod, PaymentStatus)
// 	SELECT r.UserID, r.ReservationID, 1000000, 'CreditCard', 'Successful'
// 	FROM Reservation r
// 	ORDER BY random()
// 	LIMIT 10;

// 	-- Insert Reports
// 	INSERT INTO Reports (UserID, TicketID, ReportCategory, ReportText, ReportStatus)
// 	SELECT u.UserID, t.TicketID, 'Delay', 'پرواز من ۲ ساعت تاخیر داشت.', 'Pending'
// 	FROM "User" u, Ticket t
// 	ORDER BY random()
// 	LIMIT 5;

// 	INSERT INTO Reports (UserID, PaymentID, ReportCategory, ReportText, ReportStatus)
// 	SELECT u.UserID, p.PaymentID, 'PaymentIssue', 'پول از حساب کسر شد ولی بلیط دریافت نکردم.', 'Reviewed'
// 	FROM "User" u, Payment p
// 	ORDER BY random()
// 	LIMIT 5;

// 	-- Insert TrainDetails
// 	INSERT INTO TrainDetails (TicketID, StarRating, Amenities, IsCoupeAvailable)
// 	SELECT TicketID, '4', 'WiFi, Reclining Seats, Snack Bar', true
// 	FROM Ticket WHERE VehicleType = 'Train'
// 	LIMIT 10;

// 	-- Insert FlightDetails
// 	INSERT INTO FlightDetails (TicketID, AirlineName, FlightClass, NumberOfStops, FlightNumber, OriginAirport, DestinationAirport, Amenities)
// 	SELECT TicketID, 'Iran Air', 'Economy', 1, 'IA123', 'IKA', 'MHD', 'In-flight entertainment, Meals'
// 	FROM Ticket WHERE VehicleType = 'Airplane'
// 	LIMIT 10;

// 	-- Insert BusDetails
// 	INSERT INTO BusDetails (TicketID, BusCompanyName, BusType, SeatsPerRow, Amenities)
// 	SELECT TicketID, 'Seir-o-Safar', 'VIP', '2+1', 'WiFi, Charging Ports, Comfortable Seats'
// 	FROM Ticket WHERE VehicleType = 'Bus'
// 	LIMIT 10;
// 	`
// 	db.Exec(query)
// 	sqlDb, _ := db.DB()
// 	UpdateUserPasswords(sqlDb)
	
// 	fmt.Println("Passwords updated successfully!")
// }

// func UpdateUserPasswords(db *sql.DB) error {
// 	rows, err := db.Query(`SELECT userid, hashedpassword FROM "User"`)
// 	if err != nil {
// 		return err
// 	}
// 	defer rows.Close()

// 	type User struct {
// 		ID             uuid.UUID
// 		HashedPassword string
// 	}

// 	var users []User
// 	for rows.Next() {
// 		var user User
// 		if err := rows.Scan(&user.ID, &user.HashedPassword); err != nil {
// 			return err
// 		}
// 		users = append(users, user)
// 	}

// 	for _, user := range users {
// 		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.HashedPassword), bcrypt.DefaultCost)
// 		if err != nil {
// 			return err
// 		}

// 		_, err = db.Exec(`UPDATE "User" SET hashedpassword = $1 WHERE userid = $2`, string(hashedPassword), user.ID)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }