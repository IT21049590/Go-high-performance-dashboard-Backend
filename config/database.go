package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var RowCount int

func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	DB = database
}

func Init() {
	createMaterializedViews()
	createIndexes()
	RowCount, _ = GetRowCount()

	go scheduleMaterializedViewUpdates()
}
func GetRowCount() (int, error) {
	var count int
	err := DB.Raw(`SELECT COUNT(*) FROM mv_country_product_revenue`).Scan(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func createMaterializedViews() {
	query := `
		CREATE MATERIALIZED VIEW IF NOT EXISTS mv_country_product_revenue AS
		SELECT
			ROW_NUMBER() OVER () AS row_id,
			Country,
			Product_Name,
			SUM(Total_Price) AS Total_Revenue,
			COUNT(Transaction_Id) AS Number_Of_Transactions
		FROM
			sales
		GROUP BY
			Country,
			Product_Name;
	`
	if err := DB.Exec(query).Error; err != nil {
		log.Printf("Error creating materialized view: %v", err)
	}

	query = `
		CREATE MATERIALIZED VIEW IF NOT EXISTS mv_top_product AS
				WITH latest_stock AS (
		SELECT product_id, stock_quantity
		FROM (
			SELECT 
				product_id,
				stock_quantity,
				ROW_NUMBER() OVER (PARTITION BY product_id ORDER BY added_date DESC) as rn
			FROM sales
		) AS ranked
		WHERE rn = 1
		)

		SELECT 
			s.product_id, 
			SUM(s.quantity) AS purchase_count,
			ls.stock_quantity
		FROM sales s
		JOIN latest_stock ls ON s.product_id = ls.product_id
		GROUP BY s.product_id, ls.stock_quantity
		ORDER BY purchase_count DESC
		LIMIT 20;
	`
	if err := DB.Exec(query).Error; err != nil {
		log.Printf("Error creating materialized view: %v", err)
	}
}

func refreshMaterializedViews() {
	refreshQuery := `
		REFRESH MATERIALIZED VIEW mv_country_product_revenue;
		REFRESH MATERIALIZED VIEW mv_top_product;
	`
	if err := DB.Exec(refreshQuery).Error; err != nil {
		log.Printf("Error refreshing materialized views: %v", err)
	}
}

func createIndexes() {
	indexQuery := `
		CREATE INDEX IF NOT EXISTS idx_mv_country_product_revenue_country_product 
		ON mv_country_product_revenue (row_id);
	`
	if err := DB.Exec(indexQuery).Error; err != nil {
		log.Printf("Error creating index: %v", err)
	}

	indexQuery = `
		CREATE INDEX IF NOT EXISTS idx_mv_country_product_country_revenue
		ON mv_country_product_revenue (Country, Total_Revenue DESC);
	`
	if err := DB.Exec(indexQuery).Error; err != nil {
		log.Printf("Error creating index: %v", err)
	}

	indexQuery = `
		CREATE INDEX IF NOT EXISTS idx_mv_top_product
		ON mv_top_product (product_id);
	`
	if err := DB.Exec(indexQuery).Error; err != nil {
		log.Printf("Error creating index: %v", err)
	}

	indexQuery = `
		CREATE  INDEX IF NOT EXISTS idx_country_product
		ON sales (country, product_name);
	`
	if err := DB.Exec(indexQuery).Error; err != nil {
		log.Printf("Error creating index: %v", err)
	}
}

func scheduleMaterializedViewUpdates() {
	createIndexes()

	// Run initial refresh
	refreshMaterializedViews()

	// Refresh every 6 hours
	ticker := time.NewTicker(time.Duration(AppSettings.RefreshTime) * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("Refreshing materialized views...")
			refreshMaterializedViews()
		}
	}
}
