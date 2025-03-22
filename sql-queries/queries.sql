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