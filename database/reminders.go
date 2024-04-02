package database

import (
	"database/sql"
	"time"

	"github.com/aussiebroadwan/tony/pkg/reminders"
	"github.com/bwmarrin/discordgo"

	log "github.com/sirupsen/logrus"

	_ "github.com/mattn/go-sqlite3"
)

func SetupRemindersDB(db *sql.DB, session *discordgo.Session) {
	// Create the reminders table
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS reminders (
		id INTEGER PRIMARY KEY,
		created_by TEXT NOT NULL,
		channel_id TEXT NOT NULL,
		trigger_time TEXT NOT NULL,
		message TEXT NOT NULL
		reminded BOOLEAN DEFAULT FALSE
	)`)
	if err != nil {
		panic(err)
	}

	// Load all reminders from the database
	rows, err := db.Query(`SELECT id, created_by, channel_id, trigger_time, message, reminded FROM reminders`)
	if err != nil {
		panic(err)
	}

	var loadReminders = make(map[int64]reminders.Reminder)

	// Iterate over each reminder
	for rows.Next() {
		var id int64
		var createdBy, channelId, triggerTime, message string
		var reminded bool

		err := rows.Scan(&id, &createdBy, &channelId, &triggerTime, &message, &reminded)
		if err != nil {
			log.WithField("src", "database").WithError(err).Error("Failed to scan reminder row")
			continue
		}

		// Parse the trigger time
		t, err := time.Parse(time.RFC3339, triggerTime)
		if err != nil {
			log.WithField("src", "database").WithError(err).Error("Failed to parse trigger time")
			continue
		}

		// If the reminder has already been reminded, skip it
		if reminded {
			loadReminders[id] = reminders.Reminder{
				ID:          id,
				CreatedBy:   createdBy,
				TriggerTime: t,

				Action: func(id int64) { /* Do nothing */ },
			}
			continue
		}

		// Add the reminder to the reminders package
		loadReminders[id] = reminders.Reminder{
			ID:          id,
			CreatedBy:   createdBy,
			TriggerTime: t,

			Action: func(id int64) {
				// Send the reminder message
				session.ChannelMessageSend(channelId, message)

				// Set the reminder as reminded
				_, err := db.Exec(`UPDATE reminders SET reminded = TRUE WHERE id = ?`, id)
				if err != nil {
					log.WithField("src", "database").WithError(err).Errorf("Failed to mark reminder %d as reminded", id)
				}
			},
		}
	}

	// Load the reminders into the reminders package
	reminders.Load(loadReminders)
}

func AddReminder(db *sql.DB, createdBy string, triggerTime time.Time, session *discordgo.Session, channelId string, message string) error {
	id := reminders.Add(triggerTime, createdBy, func(id int64) {
		// Send the reminder message
		session.ChannelMessageSend(channelId, message)

		// Set the reminder as reminded
		_, err := db.Exec(`UPDATE reminders SET reminded = TRUE WHERE id = ?`, id)
		if err != nil {
			log.WithField("src", "database").WithError(err).Errorf("Failed to mark reminder %d as reminded", id)
		}
	})

	_, err := db.Exec(`INSERT INTO reminders (id, created_by, trigger_time, message) VALUES (?, ?, ?, ?)`, id, createdBy, triggerTime.Format(time.RFC3339), message)
	return err
}

func DeleteReminder(db *sql.DB, id int64) error {
	err := reminders.Delete(id)
	if err != nil {
		return err
	}

	_, err = db.Exec(`UPDATE reminders SET reminded = TRUE WHERE id = ?`, id)
	return err
}
