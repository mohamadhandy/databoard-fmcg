package handlers

import (
	"context"
	"fmt"
	"io"
	"klikdaily-databoard/helper"
	"klikdaily-databoard/models"
	"klikdaily-databoard/usecases"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

type ProductHandlerInterface interface {
	CreateProduct(c *gin.Context)
	GetProductById(c *gin.Context)
	GetProducts(c *gin.Context)
	UpdateProduct(c *gin.Context)
	UploadImageProduct(c *gin.Context)
}

type productHandler struct {
	productUseCase usecases.ProductUseCaseInterface
}

func InitProductHandler(u usecases.ProductUseCaseInterface) ProductHandlerInterface {
	return &productHandler{
		productUseCase: u,
	}
}

func (h *productHandler) UploadImageProduct(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Error: %s", err.Error()))
		return
	}
	// Generate a unique filename for the uploaded image
	filename := fmt.Sprintf("%s%s", uuid.New().String(), filepath.Ext(file.Filename))

	// Create a new Google Cloud Storage client
	ctx := context.Background()
	storageClient, err := storage.NewClient(ctx, option.WithCredentialsFile(os.Getenv("FIREBASE_SERVICE_JSON")))
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to create Google Cloud Storage client: %s", err.Error()))
		return
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to open file: %s", err.Error()))
		return
	}
	defer src.Close()

	// Create a new bucket handle
	bucketName := os.Getenv("FIREBASE_BUCKET")
	bucket := storageClient.Bucket(bucketName)

	// Create an object writer using the bucket and filename
	obj := bucket.Object(filename)
	w := obj.NewWriter(ctx)

	// Copy the file data to the object writer
	if _, err := io.Copy(w, src); err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to upload file to Firebase Storage: %s", err.Error()))
		return
	}

	// Close the object writer to ensure the data is flushed to Firebase Storage
	if err := w.Close(); err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to close object writer: %s", err.Error()))
		return
	}

	// Generate the download URL for the uploaded image
	downloadURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, filename)
	imageUrl := helper.ImageURL{
		ImageUrl:   downloadURL,
		Error:      nil,
		Message:    "Image uploaded successfully!",
		StatusCode: http.StatusOK,
	}
	c.JSON(http.StatusOK, imageUrl)
}

func (h *productHandler) UpdateProduct(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	id := c.Param("id")
	productRequest := models.ProductRequest{
		ID: id,
	}
	c.BindJSON(&productRequest)
	result := h.productUseCase.UpdateProduct(authHeader, productRequest)
	c.JSON(result.StatusCode, result)
}

func (h *productHandler) GetProducts(c *gin.Context) {
	searchKeyword, _ := c.GetQuery("search")
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)
	productReq := models.ProductRequest{
		Page:  pageInt,
		Limit: limitInt,
	}
	result := h.productUseCase.GetProducts(productReq, searchKeyword)
	c.JSON(result.StatusCode, result)
}

func (h *productHandler) GetProductById(c *gin.Context) {
	id := c.Param("id")
	product := h.productUseCase.GetProductById(id)
	c.JSON(product.StatusCode, product)
}

func (h *productHandler) CreateProduct(c *gin.Context) {
	productRequest := models.ProductRequest{}
	authHeader := c.GetHeader("Authorization")
	c.BindJSON(&productRequest)
	productResult := h.productUseCase.CreateProduct(authHeader, productRequest)
	c.JSON(productResult.StatusCode, productResult)
}
