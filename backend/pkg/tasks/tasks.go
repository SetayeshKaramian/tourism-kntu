package tasks

import (
	"context"
	"tourism/pkg/config"

	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

const CheckResrvations = "task:check_reservations"
var db *gorm.DB

func init() {
	config.Connect()
	db = config.GetDB()
}

func NewPeriodicTask(ctx context.Context, t *asynq.Task) error {	
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	query1 := `
	WITH calculated_totals AS (
		SELECT u.userid, SUM(t.ticketprice * r.number) AS total
		FROM "User" u
		JOIN reservation r ON u.userid = r.userid
		JOIN ticket t ON t.ticketid = r.ticketid
		WHERE r.reservationstatus = 'Reserved'
		AND r.reservationtime <= NOW() - INTERVAL '10 minutes'
		GROUP BY u.userid
	)
	UPDATE "User" u
	SET credit = credit + calculated_totals.total
	FROM calculated_totals
	WHERE u.userid = calculated_totals.userid;`
	
	if err := tx.Exec(query1).Error; err != nil {
		tx.Rollback() 
		return err
	}

	query2 := `
	WITH calculated_totals AS (
		SELECT t.ticketid, SUM(r.number) AS total_number
		FROM reservation r
		JOIN ticket t ON t.ticketid = r.ticketid
		WHERE r.reservationstatus = 'Reserved'
		AND r.reservationtime <= NOW() - INTERVAL '10 minutes'
		GROUP BY t.ticketid
	)
	UPDATE ticket t
	SET remainingcapacity = remainingcapacity + calculated_totals.total_number
	FROM calculated_totals
	WHERE t.ticketid = calculated_totals.ticketid;`
	
	if err := tx.Exec(query2).Error; err != nil {
		tx.Rollback() 
		return err
	}

	query3 := `
	UPDATE reservation r
	SET reservationstatus = 'Cancelled'
	WHERE r.reservationstatus = 'Reserved'
	AND r.reservationtime <= NOW() - INTERVAL '10 minutes';`
	
	if err := tx.Exec(query3).Error; err != nil {
		tx.Rollback() 
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}