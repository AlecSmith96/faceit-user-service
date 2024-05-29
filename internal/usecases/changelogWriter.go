package usecases

import "github.com/AlecSmith96/faceit-user-service/internal/entities"

type ChangelogWriter interface {
	PublishChangelogEntry(entry entities.ChangelogEntry) error
}
