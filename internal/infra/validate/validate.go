package validate

import "victo/wynnguardian/internal/domain/entity"

func HasAllCriterias(item *entity.ItemInstance, criteria *entity.ItemCriteria) bool {
	for id := range criteria.Modifiers {
		if _, ok := item.Stats[id]; !ok {
			return ok
		}
	}
	return true
}
