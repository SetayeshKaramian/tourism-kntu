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