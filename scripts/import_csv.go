package scripts

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"hiruna.com/project/config"
)

type Sale struct {
	TransactionId   string    `gorm:"primaryKey"`
	TransactionDate time.Time `gorm:"type:date"`
	UserId          string
	Country         string
	Region          string
	ProductId       string
	ProductName     string
	Category        string
	Price           float64
	Quantity        int
	TotalPrice      float64
	StockQuantity   int
	AddedDate       time.Time `gorm:"type:date"`
}

const batchSize = 1000
const workerCount = 10

func ImportSalesCSV(filePath string) {
	config.DB.AutoMigrate(&Sale{})

	csvFile, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to open CSV:", err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse CSV:", err)
	}

	var wg sync.WaitGroup
	batchChan := make(chan []Sale, workerCount)

	// Start workers
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for batch := range batchChan {
				if err := config.DB.CreateInBatches(batch, batchSize).Error; err != nil {
					log.Println("Batch insert error:", err)
				}
			}
		}()
	}

	// Prepare and send batches
	currentBatch := make([]Sale, 0, batchSize)
	for i, record := range records {
		if i == 0 {
			continue
		}

		qty, _ := strconv.Atoi(record[9])
		stockQty, _ := strconv.Atoi(record[11])
		price, _ := strconv.ParseFloat(record[8], 64)
		totPrice, _ := strconv.ParseFloat(record[10], 64)
		transactionDate, _ := time.Parse("2006-01-02", record[1])
		addedDate, _ := time.Parse("2006-01-02", record[12])

		sale := Sale{
			TransactionId:   record[0],
			TransactionDate: transactionDate,
			UserId:          record[2],
			Country:         record[3],
			Region:          record[4],
			ProductId:       record[5],
			ProductName:     record[6],
			Category:        record[7],
			Price:           price,
			Quantity:        qty,
			TotalPrice:      totPrice,
			StockQuantity:   stockQty,
			AddedDate:       addedDate,
		}

		currentBatch = append(currentBatch, sale)
		if len(currentBatch) >= batchSize {
			tmp := make([]Sale, len(currentBatch))
			copy(tmp, currentBatch)
			batchChan <- tmp
			currentBatch = currentBatch[:0]
		}
	}

	if len(currentBatch) > 0 {
		batchChan <- currentBatch
	}

	close(batchChan)
	wg.Wait()

	log.Println("CSV data imported successfully with goroutines.")
}
