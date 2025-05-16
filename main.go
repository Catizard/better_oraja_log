package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type ScoreDataLog struct {
	Sha256     string
	Mode       string
	Clear      int32
	Epg        int32
	Lpg        int32
	Egr        int32
	Lgr        int32
	Egd        int32
	Lgd        int32
	Ebd        int32
	Lbd        int32
	Epr        int32
	Lpr        int32
	Ems        int32
	Lms        int32
	Notes      int32
	Combo      int32
	Minbp      int32
	PlayCount  int32
	ClearCount int32
	Option     int32
	Seed       int64
	Random     int32
	Date       int64
	State      int32
}

func (ScoreDataLog) TableName() string {
	return "scoredatalog"
}

const outputPath = "better_scoredatalog.db"

func main() {
	if _, err := os.Stat(outputPath); err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}
	} else {
		panic(fmt.Errorf("%s existed, please remove it first", outputPath))
	}
	var prevFilePath string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewFilePicker().Title("scoredatalog.db:").Value(&prevFilePath),
		),
	)

	if err := form.Run(); err != nil {
		panic(err)
	}

	// Load everything into memory
	dsn := prevFilePath + "?open_mode=1"
	db, err := gorm.Open(sqlite.Open(dsn))
	if err != nil {
		panic(err)
	}
	var rawLogs []*ScoreDataLog
	if err := db.Find(&rawLogs).Error; err != nil {
		panic(err)
	}

	// Output to another file
	outputDSN := fmt.Sprintf("%s?open_mode=1", outputPath)
	outputDB, err := gorm.Open(sqlite.Open(outputDSN))
	if err := outputDB.Table("scoredatalog").AutoMigrate(&ScoreDataLog{}); err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	if err := outputDB.Model(&ScoreDataLog{}).CreateInBatches(rawLogs, 50).Error; err != nil {
		panic(err)
	}
}
