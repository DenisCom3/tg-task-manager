package event

import "time-manager/internal/entity"

func (r *Repository) GetByName(name string) (*entity.Event, error) {
	
	event, err := r.store.GetByName(name)
	if err != nil {
		return nil, err
	}

	return event, nil
}