
# Pagination Library for Go

[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
![Go version](https://img.shields.io/badge/go-1.20-blue)
[![Coverage Status](https://coveralls.io/repos/github/xans-me/gorm-pagination/badge.svg?branch=master)](https://coveralls.io/github/xans-me/gorm-pagination?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/xans-me/gorm-pagination)](https://goreportcard.com/report/github.com/xans-me/gorm-pagination)
[![GoDoc](https://pkg.go.dev/badge/github.com/xans-me/gorm-pagination.svg)](https://pkg.go.dev/github.com/xans-me/gorm-pagination)
![Contributors](https://img.shields.io/github/contributors/xans-me/gorm-pagination)

This library provides a flexible and customizable pagination solution using GORM. It includes support for filtering, sorting, and calculating summaries like sums, minimums, and maximums.

## Features
- Pagination of GORM queries with support for page size and offsets.
- Filtering with `AND` and `OR` conditions.
- Support for ordering by fields.
- Summarization (e.g., sum, min, max) of specific fields.
- Grouping by fields for aggregated queries.
- Dynamic counting for specific field values.

## Installation

To install the library, use the following command:

```bash
go install github.com/xans-me/gorm-pagination
```

## Usage

### Basic Usage

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/xans-me/gorm-pagination"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(postgres.Open("DSN"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Register routes
	r := mux.NewRouter()
	r.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		paginator := pagination.NewPaginator(
			db.Model(&Data{}),
			pagination.WithPage(1),
			pagination.WithPageSize(10),
			pagination.WithSort("trx_date desc"),
		)

		var transactions []Data
		res, err := paginator.Paginate(&transactions)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Results: %+v", res)
	}).Methods("GET")

	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", r)
}
```

### Advanced Filtering

```go
filterManager := pagination.FilterManager{}
filterManager.AddAndFilter(pagination.ComparisonFilter{
    Field:    "trx_amount",
    Operator: ">=",
    Value:    100,
})

query := filterManager.Apply(db.Model(&Transaction{}))
```

### Summary Calculation

```go
paginator := pagination.NewPaginator(
	db.Model(&Transaction{}),
	pagination.WithPage(1),
	pagination.WithPageSize(10),
	pagination.WithSummaryFields("trx_amount:sum", "trx_amount:min", "trx_amount:max"),
)

var transactions []Transaction
res, _ := paginator.Paginate(&transactions)
fmt.Println(res.Summary)
```

### License

This library is licensed under the MIT License.
