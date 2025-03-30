-----------------------------------------------------------
-- 1. نام و نام خانوادگی کسانی که هیچ بلیطی رزرو نکرده اند.
-----------------------------------------------------------
SELECT firstname, lastname 
FROM "User"
WHERE userid NOT IN (SELECT userid FROM reservation);

-----------------------------------------------------------
-- 2. نام و نام خانوادگی کسانی که حداقل یک بلیط رزرو کرده اند.
-----------------------------------------------------------
SELECT firstname, lastname 
FROM "User"
WHERE userid IN (SELECT userid FROM reservation);

-----------------------------------------------------------
-- 3. مجموع پرداختی هر کاربر در هر ماه.
-----------------------------------------------------------
SELECT userid, DATE_TRUNC('month', paymenttime) AS month, SUM(paymentamount)
FROM payment
GROUP BY userid, month

-----------------------------------------------------------
-- 4. کاربران هر شهر که فقط یک رزرو انجام داده اند.
-----------------------------------------------------------
SELECT city, userid 
FROM "User"
WHERE userid IN (
    SELECT userid 
    FROM payment 
	  WHERE paymentstatus = 'Successful'
    GROUP BY userid 
    HAVING COUNT(userid) = 1
)


-----------------------------------------------------------
-- 5. کاربری که جدید ترین بلیط را خریداری کرده است.
-----------------------------------------------------------
SELECT * FROM "User"
WHERE userid IN (
	SELECT userid FROM payment
	WHERE paymentstatus = 'Successful' 
	ORDER BY paymenttime DESC
	LIMIT 1
)

-----------------------------------------------------------
-- 6. ایمیل کاربرانی که مجموع پرداختی آن ها از میانگین بیشتر است.
-----------------------------------------------------------
SELECT u.email
FROM "User" u
JOIN payment p ON u.userid = p.userid
GROUP BY u.Email
HAVING SUM(p.paymentamount) > (SELECT AVG(paymentamount) FROM payment)

-----------------------------------------------------------
-- 7. تعداد بلیط های فروخته شده به تفکیک وسیله نقلیه.
-----------------------------------------------------------
SELECT t.vehicletype, COUNT(*) as numbers FROM ticket t 
JOIN reservation r ON t.ticketid = r.ticketid 
JOIN payment p ON p.reservationid = r.reservationid
WHERE p.paymentstatus = 'Successful'
GROUP BY t.vehicletype

-----------------------------------------------------------
-- 8. ام 3 کاربر با بیشترین خرید بلیط در هفته اخیر را برگردانید.
-- we used inner join to insure just including reservations that have a corresponding successful payment.
-----------------------------------------------------------
SELECT u.FirstName, u.LastName, COUNT(r.ReservationID) AS TicketCount
FROM Reservation r
JOIN "User" u ON r.UserID = u.UserID
JOIN Payment p ON r.ReservationID = p.ReservationID
WHERE p.PaymentStatus = 'Successful' 
  AND p.PaymentTime >= NOW() - INTERVAL '7 days'
GROUP BY u.UserID, u.FirstName, u.LastName
ORDER BY TicketCount DESC
LIMIT 3;


-----------------------------------------------------------
-- 9.تعداد بلیطهای فروختهشده در استان تهران را به تفکیک شهر نمایش دهید.
-----------------------------------------------------------

SELECT u.City, COUNT(r.ReservationID) AS TicketsSold
FROM Reservation r
JOIN "User" u ON r.UserID = u.UserID
JOIN Payment p ON r.ReservationID = p.ReservationID
WHERE p.PaymentStatus = 'Successful'
  AND u.City IN ('Tehran', 'Karaj', 'Eslamshahr', 'Varamin', 'Shahriar', 'Pakdasht', 'Robat Karim') -- لیست شهرهای استان تهران
GROUP BY u.City
ORDER BY TicketsSold DESC;

-----------------------------------------------------------
-- 10. نام شهرهایی که قدیمیترین کاربر ثبتنام شده در سیستم از آنجا خرید داشته است را لیست کنید.
-----------------------------------------------------------
SELECT u.City, COUNT(r.ReservationID) AS TicketsSold
FROM Reservation r
JOIN "User" u ON r.UserID = u.UserID
JOIN Payment p ON r.ReservationID = p.ReservationID
WHERE p.PaymentStatus = 'Successful'
  AND u.City IN ('Tehran', 'Karaj', 'Eslamshahr', 'Varamin', 'Shahriar', 'Pakdasht', 'Robat Karim') -- لیست شهرهای استان تهران
GROUP BY u.City
ORDER BY TicketsSold DESC;

-----------------------------------------------------------
-- 11. نام پشتیبانهای سایت را لیست کنید.
-----------------------------------------------------------
select firstname, lastname from "User" where usertype = 'Supporter';


-----------------------------------------------------------
-- 12. نام کاربرانی که حداقل 2 بلیط در سیستم خریداری کردهاند را برگردانید.
-----------------------------------------------------------
SELECT u.FirstName, u.LastName
FROM "User" u
JOIN Reservation r ON u.UserID = r.UserID
JOIN Payment p ON r.ReservationID = p.ReservationID
WHERE p.PaymentStatus = 'Successful'
GROUP BY u.UserID, u.FirstName, u.LastName
HAVING COUNT(p.PaymentID) >= 2;

-----------------------------------------------------------
-- 13. نام کاربرانی را لیست کنید که حداکثر 2 بلیط از یک وسیله نقلیه خاص (مثلاً قطار) خریدهاند.
-----------------------------------------------------------
SELECT u.FirstName, u.LastName
FROM "User" u
WHERE NOT EXISTS (
    SELECT 1
    FROM Reservation r
    JOIN Payment p ON r.ReservationID = p.ReservationID
    JOIN Ticket t ON r.TicketID = t.TicketID
    WHERE r.UserID = u.UserID
      AND p.PaymentStatus = 'Successful'
    GROUP BY t.VehicleType
    HAVING COUNT(*) > 2
);

--------------------------------------------------------
--14.ایمیل یا شماره تلفن کاربرانی که از تمام وسایل نقلیه (هواپیما، قطار و اتوبوس) حداقل یک بار بلیط خریده اند را برگردانید.
--------------------------------------------------------

SELECT DISTINCT u.Email, u.PhoneNumber
FROM "User" u
JOIN Reservation r ON u.UserID = r.UserID
JOIN Ticket t ON r.TicketID = t.TicketID
WHERE 'Airplane' IN (SELECT VehicleType FROM Ticket WHERE TicketID = r.TicketID)
  AND 'Train' IN (SELECT VehicleType FROM Ticket WHERE TicketID = r.TicketID)
  AND 'Bus' IN (SELECT VehicleType FROM Ticket WHERE TicketID = r.TicketID);

--------------------------------------------------------
--15.اطلاعات بلیط های خریداری شده امروز را با ترتیب ساعت خرید لیست کنید..
--------------------------------------------------------

SELECT t.*, r.ReservationTime
FROM Ticket t
JOIN Reservation r ON t.TicketID = r.TicketID
WHERE DATE(r.ReservationTime) = CURRENT_DATE
ORDER BY r.ReservationTime;

--------------------------------------------------------
--16.دومین بلیط پرفروش در بین کل بلیط ها را نمایش دهید.
--------------------------------------------------------

SELECT TicketID, COUNT(*) AS SalesCount
FROM Reservation
GROUP BY TicketID
ORDER BY SalesCount DESC
LIMIT 1 OFFSET 1;

--------------------------------------------------------
--17. نام پشتیبان با بیشترین تعداد لغو رزرو بلیط، همراه با درصد لغوها را برگردانید.
--------------------------------------------------------

SELECT u.FirstName, u.LastName,
       COUNT(r.ReservationID) AS CanceledCount,
       ROUND((COUNT(r.ReservationID) * 100.0) / (SELECT COUNT(*) FROM Reservation), 2) AS CancelPercentage
FROM "User" u
JOIN Reservation r ON u.UserID = r.UserID
WHERE u.UserType = 'Supporter' AND r.ReservationStatus = 'Cancelled'
GROUP BY u.UserID
ORDER BY CanceledCount DESC
LIMIT 1;

--------------------------------------------------------
--18.نام خانوادگی کاربری که بیشترین تعداد بلیط کنسل شده دارد را به "ردینگتون" تغییر دهید.
--------------------------------------------------------

UPDATE "User"
SET LastName = 'ردینگتون'
WHERE UserID = (
    SELECT r.UserID
    FROM Reservation r
    WHERE r.ReservationStatus = 'Cancelled'
    GROUP BY r.UserID
    ORDER BY COUNT(*) DESC
    LIMIT 1
);

--------------------------------------------------------
--19.تمام بلیط های کنسل شده کاربر ردینگتون را حذف کنید.
--------------------------------------------------------

DELETE FROM Reservation
WHERE UserID = (SELECT UserID FROM "User" WHERE LastName = 'ردینگتون')
  AND ReservationStatus = 'Cancelled';

--------------------------------------------------------
--20.تمام بلیط های کنسل شده در سیستم را پاک کنید.
--------------------------------------------------------

DELETE FROM Reservation
WHERE ReservationStatus = 'Cancelled';

--------------------------------------------------------
--21.قیمت بلیط هایی که دیروز توسط شرکت هواپیمایی ماهان فروخته شده اند را ٪1۰ کاهش دهید.
--------------------------------------------------------

UPDATE Ticket t
SET TicketPrice = TicketPrice * 0.9
WHERE t.TicketID IN (
    SELECT f.TicketID 
    FROM FlightDetails f
    JOIN Ticket t ON f.TicketID = t.TicketID
    JOIN Reservation r ON t.TicketID = r.TicketID
    JOIN Payment p ON p.ReservationID = r.ReservationID
    WHERE f.AirlineName = 'Mahan' 
    AND p.PaymentStatus = 'Successful'
    AND DATE(p.PaymentTime) = CURRENT_DATE - INTERVAL '1 day'
);

--------------------------------------------------------
--22. موضوع و تعداد گزارش ها را برای بلیط با بیشترین تعداد گزارش، نمایش دهید.
--------------------------------------------------------

SELECT r.ReportCategory, COUNT(*) AS ReportCount
FROM Reports r
WHERE r.TicketID = (
    SELECT TicketID
    FROM Reports
    GROUP BY TicketID
    ORDER BY COUNT(*) DESC
    LIMIT 1
)
GROUP BY r.ReportCategory;

