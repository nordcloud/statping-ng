package messages

import (
	"github.com/nordcloud/statping-ng/database"
	"github.com/nordcloud/statping-ng/types/errors"
	"github.com/nordcloud/statping-ng/utils"
)

var (
	db  database.Database
	log = utils.Log.WithField("type", "message")
)

func SetDB(database database.Database) {
	db = database.Model(&Message{})
}

func Find(id int64) (*Message, error) {
	var message Message
	q := db.Where("id = ?", id).Find(&message)
	if q.Error() != nil {
		return nil, errors.Missing(message, id)
	}
	return &message, q.Error()
}

func All() []*Message {
	var messages []*Message
	db.Find(&messages)
	return messages
}

func (m *Message) Create() error {
	q := db.Create(m)
	return q.Error()
}

func (m *Message) Update() error {
	q := db.Update(m)
	return q.Error()
}

func (m *Message) Delete() error {
	q := db.Delete(m)
	return q.Error()
}
