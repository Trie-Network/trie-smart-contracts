package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func initDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./assets.db")
	if err != nil {
		log.Fatal(err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS assets (
		Id TEXT PRIMARY KEY,
		Name TEXT NOT NULL,
		Main_Category TEXT NOT NULL,
		Secondary_Category TEXT NOT NULL,
		Description TEXT NOT NULL,
		Metrics TEXT NOT NULL,
		Asset_Type TEXT,
		Asset_Id TEXT NULL,
		Owner TEXT NOT NULL
	);
	`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func CreateAssetHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var asset Asset
		if err := c.ShouldBindJSON(&asset); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		asset.ID = uuid.New().String()
		stmt, err := db.Prepare(`
			INSERT INTO assets(id, name, main_category, secondary_category, description, metrics, asset_type, owner, asset_id)
			VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)
		`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer stmt.Close()
		_, err = stmt.Exec(asset.ID, asset.Name, asset.MainCategory, asset.SecondaryCategory, asset.Description, asset.Metrics, asset.AssetType, asset.Owner, asset.AssetID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, asset)
	}
}

func GetAssetByIDHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var asset Asset

		row := db.QueryRow(`
			SELECT id, name, main_category, secondary_category, description, metrics, asset_type, owner, asset_id
			FROM assets WHERE asset_id = ?
		`, id)

		err := row.Scan(&asset.ID, &asset.Name, &asset.MainCategory, &asset.SecondaryCategory, &asset.Description, &asset.Metrics, &asset.AssetType, &asset.Owner, &asset.AssetID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
		c.JSON(http.StatusOK, asset)
	}
}

func GetAssetByOwnerHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		owner := c.Param("owner")
		rows, err := db.Query(`
			SELECT id, name, main_category, secondary_category, description, metrics, asset_type, owner, asset_id
			FROM assets WHERE owner = ?
		`, owner)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var assets []Asset
		for rows.Next() {
			var asset Asset
			if err := rows.Scan(&asset.ID, &asset.Name, &asset.MainCategory, &asset.SecondaryCategory, &asset.Description, &asset.Metrics, &asset.AssetType, &asset.Owner, &asset.AssetID); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			assets = append(assets, asset)
		}
		c.JSON(http.StatusOK, assets)
	}
}

func GetAssetByCategoryHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		category := c.Param("category")
		if category == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "category is required"})
			return
		}
		rows, err := db.Query("SELECT id, name, main_category, secondary_category, description, metrics, asset_type, owner, asset_id FROM assets WHERE asset_type = ?", category)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()
		var assets []Asset
		for rows.Next() {
			var asset Asset
			err := rows.Scan(&asset.ID, &asset.Name, &asset.MainCategory, &asset.SecondaryCategory, &asset.Description, &asset.Metrics, &asset.AssetType, &asset.Owner, &asset.AssetID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			assets = append(assets, asset)
		}
		c.JSON(http.StatusOK, assets)
	}
}

func GetAllAssetsHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query(`
			SELECT id, name, main_category, secondary_category, description, metrics, asset_type, owner, asset_id 
			FROM assets
		`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var assets []Asset
		for rows.Next() {
			var asset Asset
			if err := rows.Scan(
				&asset.ID,
				&asset.Name,
				&asset.MainCategory,
				&asset.SecondaryCategory,
				&asset.Description,
				&asset.Metrics,
				&asset.AssetType,
				&asset.Owner,
				&asset.AssetID,
			); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			assets = append(assets, asset)
		}

		c.JSON(http.StatusOK, assets)
	}
}

func main() {
	db := initDB()
	defer db.Close()

	r := gin.Default()
	r.POST("/api/createAsset", CreateAssetHandler(db))
	r.GET("/api/getAsset/:id", GetAssetByIDHandler(db))
	r.GET("/api/getAssetByOwner/:owner", GetAssetByOwnerHandler(db))
	r.GET("/api/getAssetByCategory/:category", GetAssetByCategoryHandler(db))
	r.GET("/api/getAllAssets", GetAllAssetsHandler(db))

	fmt.Println("Server running on http://localhost:8081")
	r.Run(":8081")
}
/*
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/google/uuid"
)

func main() {
	db, err := sql.Open("sqlite3", "./assets.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTable := `
	CREATE TABLE IF NOT EXISTS assets (
		id TEXT PRIMARY KEY,
		name TEXT,
		main_category TEXT,
		secondary_category TEXT,
		description TEXT,
		metrics TEXT,
		asset_type TEXT,
		owner TEXT,
		asset_id TEXT NULL
	);
	`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}
	
	createInferenceTable := `
	CREATE TABLE IF NOT EXISTS inference_record_queue (
		did TEXT NOT NULL,
		asset_id TEXT NOT NULL
	);
	`
	_, err = db.Exec(createInferenceTable)
	if err != nil {
		log.Fatal(err)
	}
	r := gin.Default()
	r.POST("/api/createAsset", CreateAssetHandler(db))
	r.GET("/api/getAsset/:id", GetAssetHandler(db))
	r.GET("/api/getAssetsByCategory/:category", GetAssetsByCategoryHandler(db))
	r.GET("/api/getAssetsByUserDID/:did", GetAssetsByUserDIDHandler(db))
	r.GET("/api/getAssetsByOwner/:owner", GetAssetsByOwnerHandler(db))

	fmt.Println(" Server running on http://localhost:8080")
	r.Run(":8080")
}

func CreateAssetHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var asset Asset
		if err := c.ShouldBindJSON(&asset); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if asset.ID == "" {
			asset.ID = uuid.New().String()
		}
		query := `
			INSERT INTO assets (id, name, main_category, secondary_category, description, metrics, asset_type, owner, asset_id)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`
		_, err := db.Exec(query,
			asset.ID, asset.Name, asset.MainCategory, asset.SecondaryCategory,
			asset.Description, string(asset.Metrics), asset.AssetType, asset.Owner, asset.AssetID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, asset)
	}
}

func GetAssetHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		query := `
			SELECT id, name, main_category, secondary_category, description, metrics, asset_type, owner, asset_id
			FROM assets WHERE id = ?
		`
		row := db.QueryRow(query, id)
		var a Asset
		if err := row.Scan(
			&a.ID, &a.Name, &a.MainCategory, &a.SecondaryCategory,
			&a.Description, &a.Metrics, &a.AssetType, &a.Owner, &a.AssetID,
		); err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
		c.JSON(http.StatusOK, a)
	}
}

func GetAssetsByCategoryHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) { 
		category := c.Param("category")
		rows, err := db.Query(`
			SELECT id, name, main_category, secondary_category, description, metrics, asset_type, owner, asset_id
			FROM assets WHERE asset_type = ?`, category)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var assets []Asset
		for rows.Next() {
			var a Asset
			if err := rows.Scan(
				&a.ID, &a.Name, &a.MainCategory, &a.SecondaryCategory,
				&a.Description, &a.Metrics, &a.AssetType, &a.Owner, &a.AssetID,
			); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			assets = append(assets, a)
		}

		c.JSON(http.StatusOK, assets)
	}
}

func GetAssetsByUserDIDHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userDID := c.Param("did")

		rows, err := db.Query(`SELECT asset_id FROM inference_record_queue WHERE did = ?`, userDID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()
		var assetIDs []string
		for rows.Next() {
			var assetID string
			if err := rows.Scan(&assetID); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			assetIDs = append(assetIDs, assetID)
		}

		if len(assetIDs) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "No assets found for given user DID"})
			return
		}

		placeholders := strings.Repeat("?,", len(assetIDs))
		placeholders = strings.TrimRight(placeholders, ",")

		query := fmt.Sprintf(`
			SELECT id, name, main_category, secondary_category, description, metrics, asset_type, owner, asset_id
			FROM assets WHERE id IN (%s)`, placeholders)

		args := make([]interface{}, len(assetIDs))
		for i, id := range assetIDs {
			args[i] = id
		}

		assetRows, err := db.Query(query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer assetRows.Close()

		var assets []Asset
		for assetRows.Next() {
			var a Asset
			if err := assetRows.Scan(
				&a.ID, &a.Name, &a.MainCategory, &a.SecondaryCategory,
				&a.Description, &a.Metrics, &a.AssetType, &a.Owner, &a.AssetID,
			); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			assets = append(assets, a)
		}

		c.JSON(http.StatusOK, assets)
	}
}

func GetAssetsByOwnerHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		owner := c.Param("owner")

		rows, err := db.Query(`
			SELECT id, name, main_category, secondary_category, description, metrics, asset_type, owner, asset_id
			FROM assets WHERE owner = ?`, owner)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()
		var assets []Asset
		for rows.Next() {
			var a Asset
			if err := rows.Scan(
				&a.ID, &a.Name, &a.MainCategory, &a.SecondaryCategory,
				&a.Description, &a.Metrics, &a.AssetType, &a.Owner, &a.AssetID,
			); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			assets = append(assets, a)
		}
		if len(assets) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "No assets found for this owner"})
			return
		}
		c.JSON(http.StatusOK, assets)
	}
}
*/