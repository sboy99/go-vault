package strategies

type DatabaseEnum string

const (
	POSTGRES DatabaseEnum = "postgres"
	MYSQL    DatabaseEnum = "mysql"
	MONGO    DatabaseEnum = "mongo"
)

type IDatabase interface {
	Connect(name string, host string, port int, username string, password string) error
	Ping() error
	Backup() ([]byte, error)
	Restore(data []byte) error
	Close() error
}

type Database struct {
	dbMap map[DatabaseEnum]IDatabase
}

func NewDatabase() *Database {
	return &Database{
		dbMap: map[DatabaseEnum]IDatabase{
			POSTGRES: NewPostgresDB(),
			MYSQL:    nil,
			MONGO:    nil,
		},
	}
}

func (d *Database) Connect(dbType DatabaseEnum, name string, host string, port int, username string, password string) error {
	db := d.getDatabase(dbType)
	if err := db.Connect(name, host, port, username, password); err != nil {
		return err
	}
	return nil
}

func (d *Database) Close(dbType DatabaseEnum) error {
	db := d.getDatabase(dbType)
	if err := db.Close(); err != nil {
		return err
	}
	return nil
}

func (d *Database) Ping(dbType DatabaseEnum) error {
	db := d.getDatabase(dbType)
	if err := db.Ping(); err != nil {
		return err
	}
	return nil
}

func (d *Database) Backup(dbType DatabaseEnum) ([]byte, error) {
	db := d.getDatabase(dbType)
	data, err := db.Backup()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (d *Database) Restore(dbType DatabaseEnum, data []byte) error {
	db := d.getDatabase(dbType)
	if err := db.Restore(data); err != nil {
		return err
	}
	return nil
}

func (d *Database) getDatabase(dbType DatabaseEnum) IDatabase {
	return d.dbMap[dbType]
}
