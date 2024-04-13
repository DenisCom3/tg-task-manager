package event

import "time-manager/internal/entity"

func (r *Repository) GetAllWithStatus(status string) ([]entity.Event, error) {

	events, err := r.store.GetAllWithStatus(status)
	if err != nil {
		return []entity.Event{}, err
	}
	return events, nil
}