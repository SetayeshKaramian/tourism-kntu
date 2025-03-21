-- Insert Users
INSERT INTO "User" (FirstName, LastName, Email, PhoneNumber, City, HashedPassword, AccountStatus, UserType)
VALUES
('John', 'Doe', 'john.doe@example.com', '09123456789', 'New York', 'hashed_pass_1', 'Active', 'Passenger'),
('Jane', 'Smith', 'jane.smith@example.com', '09234567890', 'Los Angeles', 'hashed_pass_2', 'Active', 'Supporter'),
('Alice', 'Brown', 'alice.brown@example.com', '09345678901', 'Chicago', 'hashed_pass_3', 'Inactive', 'Passenger'),
('Bob', 'Johnson', 'bob.johnson@example.com', '09456789012', 'Houston', 'hashed_pass_4', 'Active', 'Passenger'),
('Charlie', 'Williams', 'charlie.williams@example.com', '09567890123', 'Phoenix', 'hashed_pass_5', 'Active', 'Supporter'),
('David', 'Miller', 'david.miller@example.com', '09678901234', 'Philadelphia', 'hashed_pass_6', 'Inactive', 'Passenger'),
('Emma', 'Davis', 'emma.davis@example.com', '09789012345', 'San Antonio', 'hashed_pass_7', 'Active', 'Passenger'),
('Frank', 'Wilson', 'frank.wilson@example.com', '09890123456', 'San Diego', 'hashed_pass_8', 'Active', 'Supporter'),
('Grace', 'Taylor', 'grace.taylor@example.com', '09901234567', 'Dallas', 'hashed_pass_9', 'Inactive', 'Passenger'),
('Henry', 'Anderson', 'henry.anderson@example.com', '09112345678', 'San Jose', 'hashed_pass_10', 'Active', 'Passenger');

-- Insert Tickets
INSERT INTO Ticket (VehicleType, Origin, Destination, DepartureTime, ArrivalTime, TicketPrice, RemainingCapacity, TravelClass)
VALUES
('Airplane', 'New York', 'Los Angeles', '2025-04-01 08:00', '2025-04-01 11:00', 250.00, 50, 'Economy'),
('Train', 'Chicago', 'Houston', '2025-04-02 10:00', '2025-04-02 18:00', 120.00, 100, 'Business'),
('Bus', 'Phoenix', 'Philadelphia', '2025-04-03 09:00', '2025-04-04 09:00', 80.00, 40, 'VIP'),
('Airplane', 'San Antonio', 'San Diego', '2025-04-05 07:00', '2025-04-05 09:00', 180.00, 60, 'Business'),
('Train', 'Dallas', 'San Jose', '2025-04-06 06:00', '2025-04-06 12:00', 90.00, 70, 'Economy'),
('Bus', 'Los Angeles', 'New York', '2025-04-07 20:00', '2025-04-09 08:00', 150.00, 35, 'Regular'),
('Airplane', 'Houston', 'Chicago', '2025-04-08 15:00', '2025-04-08 18:00', 200.00, 45, 'First'),
('Train', 'San Diego', 'Phoenix', '2025-04-09 11:00', '2025-04-09 17:00', 110.00, 80, 'Business'),
('Bus', 'Philadelphia', 'Dallas', '2025-04-10 08:00', '2025-04-11 08:00', 75.00, 30, 'Sleeper'),
('Airplane', 'San Jose', 'New York', '2025-04-11 14:00', '2025-04-11 20:00', 300.00, 55, 'Economy'),
('Train', 'Houston', 'Los Angeles', '2025-04-12 09:00', '2025-04-12 17:00', 130.00, 60, 'Economy'),
('Bus', 'Chicago', 'Dallas', '2025-04-13 18:00', '2025-04-14 06:00', 95.00, 50, 'VIP'),
('Airplane', 'Philadelphia', 'San Antonio', '2025-04-14 10:00', '2025-04-14 14:00', 220.00, 70, 'Business'),
('Train', 'New York', 'Phoenix', '2025-04-15 05:00', '2025-04-15 15:00', 140.00, 75, 'Business'),
('Bus', 'San Diego', 'Houston', '2025-04-16 08:00', '2025-04-17 08:00', 85.00, 30, 'Regular'),
('Airplane', 'Los Angeles', 'Chicago', '2025-04-17 11:00', '2025-04-17 14:00', 270.00, 65, 'First'),
('Train', 'Dallas', 'Philadelphia', '2025-04-18 07:00', '2025-04-18 18:00', 115.00, 90, 'Economy'),
('Bus', 'San Jose', 'New York', '2025-04-19 22:00', '2025-04-21 10:00', 160.00, 40, 'Sleeper'),
('Airplane', 'Phoenix', 'San Antonio', '2025-04-20 09:00', '2025-04-20 12:00', 210.00, 50, 'Business');

-- Insert Reservations
INSERT INTO Reservation (UserID, TicketID, ReservationStatus)
SELECT u.UserID, t.TicketID, 'Reserved'
FROM "User" u, Ticket t
ORDER BY random()
LIMIT 10;

-- Insert Payments
INSERT INTO Payment (UserID, ReservationID, PaymentAmount, PaymentMethod, PaymentStatus)
SELECT r.UserID, r.ReservationID, 100.00, 'CreditCard', 'Successful'
FROM Reservation r
ORDER BY random()
LIMIT 10;

-- Insert Reports
INSERT INTO Reports (UserID, TicketID, ReportCategory, ReportText, ReportStatus)
SELECT u.UserID, t.TicketID, 'Delay', 'My flight was delayed by 2 hours.', 'Pending'
FROM "User" u, Ticket t
ORDER BY random()
LIMIT 10;

-- Insert TrainDetails
INSERT INTO TrainDetails (TicketID, StarRating, Amenities, IsCoupeAvailable)
SELECT TicketID, '4', 'WiFi, Reclining Seats, Snack Bar', true
FROM Ticket WHERE VehicleType = 'Train'
LIMIT 10;

-- Insert FlightDetails
INSERT INTO FlightDetails (TicketID, AirlineName, FlightClass, NumberOfStops, FlightNumber, OriginAirport, DestinationAirport, Amenities)
SELECT TicketID, 'AirlineX', 'Economy', 1, 'AX123', 'JFK', 'LAX', 'In-flight entertainment, Meals'
FROM Ticket WHERE VehicleType = 'Airplane'
LIMIT 10;

-- Insert BusDetails
INSERT INTO BusDetails (TicketID, BusCompanyName, BusType, SeatsPerRow, Amenities)
SELECT TicketID, 'Greyhound', 'VIP', '2+2', 'WiFi, Charging Ports, Comfortable Seats'
FROM Ticket WHERE VehicleType = 'Bus'
LIMIT 10;
