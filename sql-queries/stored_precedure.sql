-----------------------------------------------------------
--- 1. گرفتن کل بلیط های یک کاربر با استفاده از ایمیل او
-----------------------------------------------------------
CREATE OR REPLACE PROCEDURE get_tickets_of_users(
	IN user_email VARCHAR(100))
LANGUAGE plpgsql
AS $$ 
DECLARE 
	tickets refcursor;
BEGIN
	tickets := 'tickets';
	
	OPEN tickets FOR 
		SELECT t.* FROM ticket t JOIN reservation r ON t.ticketid = r.ticketid
			   WHERE r.userid IN (SELECT userid FROM "User" WHERE email = user_email);
			   
END;
$$;

CALL get_tickets_of_users('maryam.jafari@example.com');
FETCH ALL FROM tickets;

-----------------------------------------------------------
--- 2. نشان دادن نام کاربرانی که حداقل یک بار بلیط ان ها کنسل شده با استفاده از ایمیل پشتیبان
-----------------------------------------------------------
CREATE OR REPLACE PROCEDURE get_users_with_cancelled_reservation_using_support_email(
	IN support_email VARCHAR(100)
)
LANGUAGE plpgsql AS $$
DECLARE
	users refcursor;
BEGIN
	users := 'users';
	
	OPEN users FOR 
		SELECT u.firstname, u.lastname FROM "User" u JOIN reservation r 
			ON u.userid = r.userid WHERE r.reservationstatus = 'Cancelled';
END;
$$;

CALL get_users_with_cancelled_reservation_using_support_email('ehsan.shirazi@example.com');
FETCH ALL FROM users;

-----------------------------------------------------------
--- 3. لیست بلیط های خریداری شده یک شهر با دریافت نام آن.
-----------------------------------------------------------
CREATE OR REPLACE PROCEDURE get_bought_tickets_of_a_city(
	IN city VARCHAR(100)
)
LANGUAGE plpgsql AS $$
DECLARE
	tickets refcursor;
BEGIN 
	tickets := 'tickets';
	
	OPEN tickets FOR 
		SELECT t.* FROM ticket t 
		JOIN reservation r ON t.ticketid = r.ticketid
		JOIN payment p ON p.reservationid = r.reservationid
		WHERE p.paymentstatus = 'Successful' AND t.destination = city;
END;
$$;

CALL get_bought_tickets_of_a_city('Mashhad');
FETCH ALL FROM tickets;

-----------------------------------------------------------
-- عبارتی را از ورودی گرفته و بلیطهایی را که آن عبارت در نام مسافر، مسیر سفر یا کلاس بلیط آمده باشد را برگردانید.4
-----------------------------------------------------------

CREATE OR REPLACE PROCEDURE search_tickets_proc(
    IN search_term TEXT
)
LANGUAGE plpgsql
AS $$
DECLARE
    search_tickets_cursor refcursor;  -- Declare a refcursor variable
BEGIN
  -- Assign a fixed cursor name to the variable.
  search_tickets_cursor := 'search_tickets_cursor';
  
  -- Open the cursor with the desired query.
  OPEN search_tickets_cursor FOR
    SELECT t.TicketID, t.VehicleType, t.Origin, t.Destination,
           t.DepartureTime, t.ArrivalTime, t.TicketPrice,
           t.RemainingCapacity, t.TravelClass, u.FirstName,
		   u.LastName, t.TravelClass
    FROM Ticket t
    JOIN Reservation r ON t.TicketID = r.TicketID
    JOIN "User" u ON r.UserID = u.UserID
    WHERE u.FirstName ILIKE '%' || search_term || '%'
       OR u.LastName ILIKE '%' || search_term || '%'
       OR t.Origin ILIKE '%' || search_term || '%'
       OR t.Destination ILIKE '%' || search_term || '%'
       OR t.TravelClass ILIKE '%' || search_term || '%';
END;
$$;

CALL search_tickets_proc('Tehran');
FETCH ALL FROM search_tickets_cursor;


-----------------------------------------------------------
--5. شماره تلفن یا ایمیل کاربر را دریافت کرده و اطلاعات سایر کاربران همشهری او را نمایش دهید.
-----------------------------------------------------------
CREATE OR REPLACE PROCEDURE get_citizens_proc(
    IN user_contact TEXT
)
LANGUAGE plpgsql
AS $$
DECLARE
    get_citizens_cursor refcursor;  -- Declare a refcursor variable
    city_name TEXT;
BEGIN
  -- Retrieve the city of the user with the given contact information.
  SELECT City INTO city_name
  FROM "User"
  WHERE Email = user_contact OR PhoneNumber = user_contact
  LIMIT 1;
  
  IF city_name IS NULL THEN
    RAISE NOTICE 'No user found with contact: %', user_contact;
    RETURN;
  END IF;
  
  -- Assign a fixed cursor name to the variable.
  get_citizens_cursor := 'get_citizens_cursor';
  
  -- Open the cursor with the query to get other users from the same city.
  OPEN get_citizens_cursor FOR
    SELECT UserID, FirstName, LastName, Email, PhoneNumber, City
    FROM "User"
    WHERE City = city_name
      AND NOT (Email = user_contact OR PhoneNumber = user_contact);
END;
$$;


CALL get_citizens_proc('09112345678');
FETCH ALL FROM get_citizens_cursor;

-----------------------------------------------------------
--6.  کاربری که از آن تاریخ به بعد بیشترین خرید بلیط را داشتند نمایش دهید n تاریخ و تعداد را به عنوان ورودی دریافت کرده و لیست.
-----------------------------------------------------------

CREATE OR REPLACE PROCEDURE GetTopUsersByTickets_Proc(
    IN p_start_date TIMESTAMP, 
    IN p_limit INT
)
LANGUAGE plpgsql AS $$
DECLARE
    user_cursor refcursor;
BEGIN
    user_cursor := 'user_cursor';
    
    OPEN user_cursor FOR 
        SELECT 
            u.UserID, 
            u.FirstName, 
            u.LastName, 
            COUNT(r.ReservationID) AS TicketCount
        FROM "User" u
        JOIN Reservation r ON u.UserID = r.UserID
        WHERE r.ReservationTime >= p_start_date
        GROUP BY u.UserID, u.FirstName, u.LastName
        ORDER BY TicketCount DESC
        LIMIT p_limit;
END;
$$;


CALL GetTopUsersByTickets_Proc('2025-01-01', 5);
FETCH ALL FROM user_cursor;

-----------------------------------------------------------
--7.با دریافت نوع وسیله نقلیه، لیست بلیط های کنسل شده مربوط به آن را به ترتیب تاریخ نمایش دهید.
-----------------------------------------------------------

CREATE OR REPLACE PROCEDURE GetCancelledTicketsByVehicle_Proc(
    IN p_vehicle_type VARCHAR
)
LANGUAGE plpgsql AS $$
DECLARE
    ticket_cursor refcursor;
BEGIN
    ticket_cursor := 'ticket_cursor';
    
    OPEN ticket_cursor FOR 
        SELECT 
            t.TicketID, 
            t.Origin, 
            t.Destination, 
            t.DepartureTime
        FROM Ticket t
        JOIN Reservation r ON t.TicketID = r.TicketID
        WHERE r.ReservationStatus = 'Cancelled' 
          AND t.VehicleType = p_vehicle_type
        ORDER BY t.DepartureTime;
END;
$$;


CALL GetCancelledTicketsByVehicle_Proc('Airplane');
FETCH ALL FROM ticket_cursor;

-----------------------------------------------------------
--8.با دریافت موضوع گزارش، لیست کاربرانی که بیشترین گزارش در آن موضوع دارند را نمایش دهید.
-----------------------------------------------------------

CREATE OR REPLACE PROCEDURE GetTopReportersByCategory_Proc(
    IN p_report_category VARCHAR
)
LANGUAGE plpgsql AS $$
DECLARE
    report_cursor refcursor;
BEGIN
    report_cursor := 'report_cursor';
    
    OPEN report_cursor FOR 
        SELECT 
            u.UserID, 
            u.FirstName, 
            u.LastName, 
            COUNT(r.ReportID) AS ReportCount
        FROM "User" u
        JOIN Reports r ON u.UserID = r.UserID
        WHERE r.ReportCategory = p_report_category
        GROUP BY u.UserID, u.FirstName, u.LastName
        ORDER BY ReportCount DESC;
END;
$$;


CALL GetTopReportersByCategory_Proc('PaymentIssue');
FETCH ALL FROM report_cursor;
