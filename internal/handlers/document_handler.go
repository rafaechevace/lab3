package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

const rootDataPath = "db" // Root directory where documents are stored

// DocumentContent represents the expected request body containing the document JSON.
type DocumentContent struct {
	Content json.RawMessage `json:"doc_content" binding:"required"`
}

// CreateOrUpdateDocument creates or updates a document for a user.
func CreateOrUpdateDocument(c *gin.Context) {
	requestUsername := c.Param("username")
	docID := c.Param("doc_id")

	// Ensure the authenticated user can only write in their own namespace.
	authUsername, _ := c.Get("username")
	if authUsername != requestUsername {
		c.JSON(http.StatusForbidden, gin.H{"error": "No tienes permiso para acceder a este recurso"})
		return
	}

	// Bind and validate incoming JSON body.
	var docBody DocumentContent
	if err := c.ShouldBindJSON(&docBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El contenido del documento (doc_content) es requerido y debe ser un JSON válido"})
		return
	}

	// Ensure user directory exists.
	userPath := filepath.Join(rootDataPath, requestUsername)
	if err := os.MkdirAll(userPath, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el directorio del usuario"})
		return
	}

	// Write the document content to a file.
	filePath := filepath.Join(userPath, docID+".json")
	err := ioutil.WriteFile(filePath, docBody.Content, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo escribir el documento en disco"})
		return
	}

	// Return the size of the stored content.
	c.JSON(http.StatusOK, gin.H{"size": len(docBody.Content)})
}

// GetDocument returns the JSON content of a specific document.
func GetDocument(c *gin.Context) {
	requestUsername := c.Param("username")
	docID := c.Param("doc_id")

	// Authorization: ensure the requester is the owner.
	authUsername, _ := c.Get("username")
	if authUsername != requestUsername {
		c.JSON(http.StatusForbidden, gin.H{"error": "No tienes permiso para acceder a este recurso"})
		return
	}

	// Read the file and return its content as application/json.
	filePath := filepath.Join(rootDataPath, requestUsername, docID+".json")
	content, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "El documento no existe"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo leer el documento"})
		}
		return
	}

	// Validate that the file contains valid JSON (best-effort).
	var jsonData json.RawMessage
	json.Unmarshal(content, &jsonData)

	c.Data(http.StatusOK, "application/json", content)
}

// DeleteDocument removes a user's document file.
func DeleteDocument(c *gin.Context) {
	requestUsername := c.Param("username")
	docID := c.Param("doc_id")

	// Authorization: only owner may delete.
	authUsername, _ := c.Get("username")
	if authUsername != requestUsername {
		c.JSON(http.StatusForbidden, gin.H{"error": "No tienes permiso para acceder a este recurso"})
		return
	}

	filePath := filepath.Join(rootDataPath, requestUsername, docID+".json")
	err := os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "El documento no existe"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo borrar el documento"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{}) // Return empty success response
}

// GetAllDocuments lists all JSON documents for a given user.
func GetAllDocuments(c *gin.Context) {
	requestUsername := c.Param("username")

	// Authorization: only owner may list documents.
	authUsername, _ := c.Get("username")
	if authUsername != requestUsername {
		c.JSON(http.StatusForbidden, gin.H{"error": "No tienes permiso para acceder a este recurso"})
		return
	}

	// Read user directory and collect .json files.
	userPath := filepath.Join(rootDataPath, requestUsername)
	files, err := ioutil.ReadDir(userPath)
	if err != nil {
		// If the directory doesn't exist, return an empty object.
		if os.IsNotExist(err) {
			c.JSON(http.StatusOK, gin.H{}) // No documents yet
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los documentos"})
		return
	}

	allDocs := make(map[string]json.RawMessage)
	for _, file := range files {
		// Only include regular .json files.
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			docID := strings.TrimSuffix(file.Name(), ".json")
			content, err := os.ReadFile(filepath.Join(userPath, file.Name()))
			if err == nil {
				var jsonData json.RawMessage
				json.Unmarshal(content, &jsonData) // best-effort validation
				allDocs[docID] = jsonData
			}
		}
	}

	// Return a map of docID -> document JSON.
	c.JSON(http.StatusOK, allDocs)
}
