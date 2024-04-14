package event

import "time-manager/internal/entity"

func (r *Repository) Update(event entity.Event) error {

	err := r.store.Update(event)
	if err != nil {
		return err
	}
	return nil
}