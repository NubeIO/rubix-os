# Gormigrate [For more details](https://pkg.go.dev/github.com/go-gormigrate/gormigrate/v2#section-documentation)

## Usage (Sqlite)

```go
package main

import (
    "github.com/go-gormigrate/gormigrate/v2",
    "gorm.io/driver/sqlite"
)

func main() {
    db, err := gorm.Open(sqlite.Open("example.db"), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }
    
    // Create migrations
    m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
        ID: "", // Migration identifier
        Migrate: func(tx *gorm.DB) error { // Executes migration		
        },
        Rollback: func(tx *gorm.DB) error { // Executes rollback
        },
    },
    // Executes all migrations
    if err = m.Migrate();
}
```

## Example
 
```
gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
    {
        ID: "201608301400",
        Migrate: func(tx *gorm.DB) error {
            type Alert struct {
                Type     string
                Duration int
            }
            type Message struct {
                ID            uint `gorm:"AUTO_INCREMENT;primary_key;index"`
                ApplicationID uint
                Message       string `gorm:"type:text"`
                Title         string `gorm:"type:text"`
                Priority      int
                Extras        []byte
                Date          time.Time
            }
            return tx.AutoMigrate(
                &Alert{},
                &Message{},
            )
        },
        Rollback: func(tx *gorm.DB) error {
            return tx.Migrator().DropTable(
                "alerts",
                "messages",
            )
        },
    },
    {
        ID: "201608301500",
        Migrate: func(tx *gorm.DB) error {
            nodes := `CREATE TABLE "nodes" (
            "uuid"	varchar(255) UNIQUE,
            "name"	text,
            "node_type"	text,
            "help"	text,
            "in1"	text,
            "in2"	text,
            "out1_value"	text,
            "out2_value"	text,
            "node_settings"	JSON,
            PRIMARY KEY("uuid")
        );`
            return tx.Exec(nodes).Error
        },
        Rollback: func(tx *gorm.DB) error {
            return tx.Migrator().DropTable(
                "nodes",
            )
        },
    },
    {
        ID: "201608301600",
        Migrate: func(tx *gorm.DB) error {
            type Alert struct {
                Status bool
            }
            return tx.AutoMigrate(
                &Alert{},
            )
        },
        Rollback: func(tx *gorm.DB) error {
            return tx.Migrator().DropColumn("alerts", "status")
        },
    },
    {
        ID: "201608301700",
        Migrate: func(tx *gorm.DB) error {
            query := "ALTER TABLE nodes ADD COLUMN status numeric;"
            return tx.Exec(query).Error
        },
        Rollback: func(tx *gorm.DB) error {
            return tx.Migrator().DropColumn("nodes", "status")
        },
    },
})
```

