package land

// Query types
const (
	Select      = "SELECT"
	Insert      = "INSERT"
	Update      = "UPDATE"
	Delete      = "DELETE"
	CreateTable = "CREATE TABLE"
	AlterTable  = "ALTER TABLE"
	DropTable   = "DROP TABLE"
	Where       = "WHERE"
	Join        = "JOIN"
	Order       = "ORDER"
	Column      = "COLUMN"
	Columns     = "COLUMNS"
	Group       = "GROUP"
)

// Columns names
const (
	Id        string = "id"
	Name             = "name"
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
	Int4                     = "int4"
	Int8                     = "int8"
	BigInt                   = "bigint"
	Float                    = "float"
	Boolean                  = "boolean"
	Bool                     = "bool"
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
