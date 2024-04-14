package event

import (
	"time-manager/internal/entity"

)




func (r *Repository) Save(event entity.Event) error {
	
	err := r.store.Save(event.Title, event.Time, event.Status, event.ChatId)
	if err != nil {
		return err
	}

	return nil
}