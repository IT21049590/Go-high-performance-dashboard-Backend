# High Performance Go Dashboard Backend

This project is a high-performance backend built with **Go**, using the **Fiber** framework. It efficiently ingests and processes large CSV datasets using **goroutines and batch processing**, stores the data in **PostgreSQL**, and exposes fast, optimized APIs via **materialized views**.

---

## 🚀 Features

- ✅ CSV ingestion with batch processing and goroutines for optimal speed
- ✅ REST API built with Fiber (minimal and fast Go web framework)
- ✅ PostgreSQL as the database engine
- ✅ `MATERIALIZED VIEW`s for optimized and faster data queries
- ✅ Periodic refresh of materialized views based on a configurable schedule
- ✅ All API responses are optimized to return within **5 seconds**

---

## 📦 Data Model

Data extracted from the CSV includes the following fields:

| Column           | Type      | Description                        |
|------------------|-----------|------------------------------------|
| TransactionId    | `string`  | Primary key                        |
| TransactionDate  | `date`    | Date of transaction                |
| UserId           | `string`  | Unique user ID                     |
| Country          | `string`  | Country of user                    |
| Region           | `string`  | Region of user                     |
| ProductId        | `string`  | Product identifier                 |
| ProductName      | `string`  | Name of the product                |
| Category         | `string`  | Product category                   |
| Price            | `float64` | Price per unit                     |
| Quantity         | `int`     | Quantity purchased                 |
| TotalPrice       | `float64` | Price × Quantity                   |
| StockQuantity    | `int`     | Available stock quantity           |
| AddedDate        | `date`    | Date stock was added               |

---

## 🧠 Materialized Views

Two PostgreSQL materialized views are used to accelerate data retrieval:

### 1. `mv_country_product_revenue`
Aggregates total revenue and transaction count per country and product.

### 2. `mv_top_product`
Calculates the top 20 most purchased products along with their latest stock quantity.

> These views are **periodically refreshed** using a scheduled Go routine based on a configurable refresh interval (e.g., every 6 hours).

---

## ⚙️ Configuration

Create a `.env` file at the root:

