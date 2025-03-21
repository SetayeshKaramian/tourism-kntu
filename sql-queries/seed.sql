-- Insert Users
INSERT INTO "User" (FirstName, LastName, Email, PhoneNumber, City, HashedPassword, AccountStatus, UserType)
VALUES
('Ali', 'Ahmadi', 'ali.ahmadi@example.com', '09123456789', 'Tehran', 'hashed_pass_1', 'Active', 'Passenger'),
('Sara', 'Mohammadi', 'sara.mohammadi@example.com', '09234567890', 'Mashhad', 'hashed_pass_2', 'Active', 'Supporter'),
('Reza', 'Hosseini', 'reza.hosseini@example.com', '09345678901', 'Isfahan', 'hashed_pass_3', 'Inactive', 'Passenger'),
('Fatemeh', 'Karimi', 'fatemeh.karimi@example.com', '09456789012', 'Shiraz', 'hashed_pass_4', 'Active', 'Passenger'),
('Hossein', 'Ranjbar', 'hossein.ranjbar@example.com', '09567890123', 'Tabriz', 'hashed_pass_5', 'Active', 'Supporter'),
('Mehdi', 'Zarei', 'mehdi.zarei@example.com', '09678901234', 'Ahvaz', 'hashed_pass_6', 'Inactive', 'Passenger'),
('Maryam', 'Jafari', 'maryam.jafari@example.com', '09789012345', 'Qom', 'hashed_pass_7', 'Active', 'Passenger'),
('Ehsan', 'Shirazi', 'ehsan.shirazi@example.com', '09890123456', 'Kerman', 'hashed_pass_8', 'Active', 'Supporter'),
('Zahra', 'Mansouri', 'zahra.mansouri@example.com', '09901234567', 'Rasht', 'hashed_pass_9', 'Inactive', 'Passenger'),
('Hassan', 'Gholami', 'hassan.gholami@example.com', '09112345678', 'Urmia', 'hashed_pass_10', 'Active', 'Passenger');

-- Insert Tickets
INSERT INTO Ticket (VehicleType, Origin, Destination, DepartureTime, ArrivalTime, TicketPrice, RemainingCapacity, TravelClass)
VALUES
('Airplane', 'Tehran', 'Mashhad', '2025-04-01 08:00', '2025-04-01 09:30', 2500000, 50, 'Economy'),
('Train', 'Isfahan', 'Tehran', '2025-04-02 10:00', '2025-04-02 16:00', 1200000, 100, 'Business'),
('Bus', 'Shiraz', 'Tabriz', '2025-04-03 09:00', '2025-04-04 09:00', 800000, 40, 'VIP'),
('Airplane', 'Ahvaz', 'Kerman', '2025-04-05 07:00', '2025-04-05 08:00', 1800000, 60, 'Business'),
('Train', 'Qom', 'Mashhad', '2025-04-06 06:00', '2025-04-06 14:00', 900000, 70, 'Economy'),
('Bus', 'Tehran', 'Ahvaz', '2025-04-07 20:00', '2025-04-08 08:00', 1500000, 35, 'Economy'),
('Airplane', 'Rasht', 'Shiraz', '2025-04-08 15:00', '2025-04-08 16:30', 2000000, 45, 'Business'),
('Train', 'Tabriz', 'Tehran', '2025-04-09 11:00', '2025-04-09 17:00', 1100000, 80, 'Business'),
('Bus', 'Kerman', 'Isfahan', '2025-04-10 08:00', '2025-04-10 20:00', 750000, 30, 'Economy'),
('Airplane', 'Mashhad', 'Tehran', '2025-04-11 14:00', '2025-04-11 15:30', 3000000, 55, 'Economy');

-- Insert Reservations
INSERT INTO Reservation (UserID, TicketID, ReservationStatus)
SELECT u.UserID, t.TicketID, 'Reserved'
FROM "User" u, Ticket t
ORDER BY random()
LIMIT 10;

-- Insert Payments
INSERT INTO Payment (UserID, ReservationID, PaymentAmount, PaymentMethod, PaymentStatus)
SELECT r.UserID, r.ReservationID, 1000000, 'CreditCard', 'Successful'
FROM Reservation r
ORDER BY random()
LIMIT 10;

-- Insert Reports
INSERT INTO Reports (UserID, TicketID, ReportCategory, ReportText, ReportStatus)
SELECT u.UserID, t.TicketID, 'Delay', 'پرواز من ۲ ساعت تاخیر داشت.', 'Pending'
FROM "User" u, Ticket t
ORDER BY random()
LIMIT 5;

INSERT INTO Reports (UserID, PaymentID, ReportCategory, ReportText, ReportStatus)
SELECT u.UserID, p.PaymentID, 'PaymentIssue', 'پول از حساب کسر شد ولی بلیط دریافت نکردم.', 'Reviewed'
FROM "User" u, Payment p
ORDER BY random()
LIMIT 5;

-- Insert TrainDetails
INSERT INTO TrainDetails (TicketID, StarRating, Amenities, IsCoupeAvailable)
SELECT TicketID, '4', 'WiFi, Reclining Seats, Snack Bar', true
FROM Ticket WHERE VehicleType = 'Train'
LIMIT 10;

-- Insert FlightDetails
INSERT INTO FlightDetails (TicketID, AirlineName, FlightClass, NumberOfStops, FlightNumber, OriginAirport, DestinationAirport, Amenities)
SELECT TicketID, 'Iran Air', 'Economy', 1, 'IA123', 'IKA', 'MHD', 'In-flight entertainment, Meals'
FROM Ticket WHERE VehicleType = 'Airplane'
LIMIT 10;

-- Insert BusDetails
INSERT INTO BusDetails (TicketID, BusCompanyName, BusType, SeatsPerRow, Amenities)
SELECT TicketID, 'Seir-o-Safar', 'VIP', '2+1', 'WiFi, Charging Ports, Comfortable Seats'
FROM Ticket WHERE VehicleType = 'Bus'
LIMIT 10;
