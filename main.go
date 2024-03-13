package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var (
	db *gorm.DB
)

type Order struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Customer string `json:"customer"`
	Product  string `json:"product"`
	Quantity int    `json:"quantity"`
}

func main() {
	// Koneksi ke database
	var err error
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=yusrigo dbname=db_go password=123321 sslmode=disable")
	if err != nil {
		panic("gagal terhubung ke database")
	}
	defer db.Close()

	// Migrasi skema
	db.AutoMigrate(&Order{})

	// Inisialisasi router Gin
	router := gin.Default()

	// Definisikan rute
	router.POST("/orders", createOrder)
	router.GET("/orders", getOrders)
	router.PUT("/orders/:id", updateOrder)
	router.DELETE("/orders/:id", deleteOrder)

	// Mulai server
	router.Run(":8080")
}

// Buat Order
func createOrder(c *gin.Context) {
	var order Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Create(&order)
	c.JSON(http.StatusCreated, order)
}

// Dapatkan Orders
func getOrders(c *gin.Context) {
	var orders []Order
	db.Find(&orders)
	c.JSON(http.StatusOK, orders)
}

// Perbarui Order
func updateOrder(c *gin.Context) {
	id := c.Param("id")
	var order Order
	if err := db.Where("id = ?", id).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan!"})
		return
	}

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Save(&order)
	c.JSON(http.StatusOK, order)
}

// Hapus Order
func deleteOrder(c *gin.Context) {
	id := c.Param("id")
	var order Order
	if err := db.Where("id = ?", id).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan!"})
		return
	}

	db.Delete(&order)
	c.JSON(http.StatusOK, gin.H{"message": "Order berhasil dihapus"})
}
