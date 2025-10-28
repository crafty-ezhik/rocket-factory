package migrator

type PostgresMigrator interface {
	Up() error
}
