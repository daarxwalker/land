package land

// Databases types
const (
	Postgres string = "postgres"
)

// Query types
const (
	Select      = "SELECT"
	Insert      = "INSERT"
	Update      = "UPDATE"
	Delete      = "DELETE"
	CreateTable = "CREATE TABLE"
	AlterTable  = "ALTER TABLE"
	DropTable   = "DROP TABLE"
)

// Columns names
const (
	Id        string = "id"
	Vectors          = "vectors"
	CreatedAt        = "created_at"
	UpdatedAt        = "updated_at"
)

// Data types
const (
	Varchar           string = "varchar"
	Char                     = "char"
	Text                     = "text"
	Int                      = "int"
	BigInt                   = "bigint"
	Float                    = "float"
	Boolean                  = "boolean"
	Jsonb                    = "jsonb"
	ArrayText                = "text[]"
	ArrayInt                 = "integer[]"
	TsVector                 = "tsvector"
	Timestamp                = "timestamp"
	TimestampWithZone        = "timestampz"
	Serial                   = "serial"
)

// Default values
const (
	DefaultLimit            = 20
	CurrentTimestamp string = "CURRENT_TIMESTAMP"
)
