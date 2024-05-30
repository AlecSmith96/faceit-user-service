package usecases

import "github.com/AlecSmith96/faceit-user-service/internal/entities"

//go:generate mockgen --build_flags=--mod=mod -destination=../../mocks/changelogWriter.go  . "ChangelogWriter"
type ChangelogWriter interface {
	PublishChangelogEntry(entry entities.ChangelogEntry) error
}
