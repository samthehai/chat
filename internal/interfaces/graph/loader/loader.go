package loader

import (
	"github.com/graph-gophers/dataloader"
	"github.com/samthehai/chat/internal/domain/entity"
)

func getIDsFromKeys(keys dataloader.Keys) []entity.ID {
	ids := make([]entity.ID, 0, len(keys))

	for _, key := range keys {
		ids = append(ids, key.(entity.ID))
	}

	return ids
}

func fillUpResultsWithError(size int, err error) []*dataloader.Result {
	results := make([]*dataloader.Result, 0, size)

	for i := 0; i < size; i++ {
		results = append(results, &dataloader.Result{
			Error: err,
		})
	}

	return results
}
