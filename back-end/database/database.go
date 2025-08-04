package database

import (
	"appseclabsplataform/config"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	Conn *gorm.DB
}

func NewDatabase(config *config.Config) *Database {
	var db *gorm.DB
	var err error

	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(config.DatabaseConfig.ConString), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err == nil {
			break
		}

		log.Printf("Tentativa %d/%d: Erro ao conectar ao PostgreSQL: %v", i+1, maxRetries, err)
		if i < maxRetries-1 {
			log.Println("Tentando novamente em 2 segundos...")
			time.Sleep(2 * time.Second)
		}
	}

	if err != nil {
		log.Fatalf("Erro ao conectar ao PostgreSQL após %d tentativas: %v", maxRetries, err)
	}

	log.Println("Conectado ao PostgreSQL com sucesso via GORM.")

	return &Database{
		Conn: db,
	}
}
