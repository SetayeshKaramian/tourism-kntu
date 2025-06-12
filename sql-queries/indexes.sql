-- 1. “User” filters/join columns
CREATE INDEX IF NOT EXISTS idx_user_city
  ON "User"(City);

CREATE INDEX IF NOT EXISTS idx_user_usertype
  ON "User"(UserType);


-- 2. “Reservation” joins / filters / ordering
CREATE INDEX IF NOT EXISTS idx_reservation_userid
  ON Reservation(UserID);

CREATE INDEX IF NOT EXISTS idx_reservation_ticketid
  ON Reservation(TicketID);

CREATE INDEX IF NOT EXISTS idx_reservation_status
  ON Reservation(ReservationStatus);

CREATE INDEX IF NOT EXISTS idx_reservation_time
  ON Reservation(ReservationTime);


-- 3. “Payment” joins / filters / grouping / sorting
CREATE INDEX IF NOT EXISTS idx_payment_userid
  ON Payment(UserID);

CREATE INDEX IF NOT EXISTS idx_payment_reservationid
  ON Payment(ReservationID);

-- used for “WHERE PaymentStatus='Successful' AND PaymentTime >= …”,
-- and retrieving the latest or monthly sums
CREATE INDEX IF NOT EXISTS idx_payment_status_time
  ON Payment(PaymentStatus, PaymentTime);


-- 4. “Ticket” joins / filters / grouping
CREATE INDEX IF NOT EXISTS idx_ticket_vehicletype
  ON Ticket(VehicleType);

CREATE INDEX IF NOT EXISTS idx_ticket_destination
  ON Ticket(Destination);

CREATE INDEX IF NOT EXISTS idx_ticket_travelclass
  ON Ticket(TravelClass);


-- 5. “FlightDetails” filter by airline
CREATE INDEX IF NOT EXISTS idx_flightdetails_airlinename
  ON FlightDetails(AirlineName);


-- 6. “Reports” grouping / filtering
CREATE INDEX IF NOT EXISTS idx_reports_ticketid
  ON Reports(TicketID);

CREATE INDEX IF NOT EXISTS idx_reports_category
  ON Reports(ReportCategory);
