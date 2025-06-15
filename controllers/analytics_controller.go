package controllers

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"hiruna.com/project/config"
)

type ProductCount struct {
	ProductId     string
	PurchaseCount int
	StockQuantity int
}

type MonthlySales struct {
	Month string
	Count int
}

type CountryRevenue struct {
	Country      string
	ProductName  string
	TotalRevenue float64
	Transactions int
}

type RegionRevenue struct {
	Region       string
	TotalRevenue float64
	ItemsSold    int
}

func GetTopProducts(c *fiber.Ctx) error {
	db := config.DB
	var products []ProductCount
	db.Raw(`
    SELECT *
    FROM mv_top_product;
	`).Scan(&products)
	return c.JSON(products)
}

func GetMonthlySales(c *fiber.Ctx) error {
	db := config.DB
	var sales []MonthlySales
	db.Raw(`
		SELECT 
    TO_CHAR(transaction_date, 'Month') AS month,
    COUNT(*) AS count
FROM 
    sales
GROUP BY 
    TO_CHAR(transaction_date, 'Month'), EXTRACT(MONTH FROM transaction_date)
ORDER BY 
    count DESC
LIMIT 12;


	`).Scan(&sales)
	return c.JSON(sales)
}

func GetTopRegions(c *fiber.Ctx) error {
	db := config.DB
	var regions []RegionRevenue
	db.Raw(`
		SELECT 
    region,
    SUM(total_price) AS total_revenue,
    SUM(quantity) AS items_sold
FROM 
    sales
GROUP BY 
    region
ORDER BY 
    total_revenue DESC
LIMIT 30;
	`).Scan(&regions)
	return c.JSON(regions)
}

func GetChunkedViewData(c *fiber.Ctx) error {
	fmt.Println("Start time:", time.Now())
	db := config.DB

	chunkSize := 1000000
	var wg sync.WaitGroup
	var mu sync.Mutex
	finalResults := make([]CountryRevenue, 0)

	for start := 1; start <= config.RowCount; start += chunkSize {
		end := start + chunkSize - 1
		wg.Add(1)

		go func(start, end int) {
			defer wg.Done()

			var chunk []CountryRevenue
			query := `
				SELECT country, product_name, total_revenue, number_of_transactions AS transactions
				FROM mv_country_product_revenue
				WHERE row_id BETWEEN ? AND ?
				ORDER BY row_id;
			`

			if err := db.Raw(query, start, end).Scan(&chunk).Error; err != nil {
				log.Printf("Error fetching chunk %d-%d: %v", start, end, err)
				return
			}

			mu.Lock()
			finalResults = append(finalResults, chunk...)
			mu.Unlock()
		}(start, end)
	}

	wg.Wait()
	fmt.Println("end time:", time.Now())
	log.Println("Total rows fetched:", len(finalResults))
	return c.JSON(finalResults)
}
