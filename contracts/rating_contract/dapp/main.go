package main

import (
	//"fmt"
	//"net/http"
    "math"
	"github.com/gin-gonic/gin"
)

func RoundToPrecision(val float64, precision int, tolerance int) float64 {
    factor := math.Pow(10, float64(precision))
    return math.Round(val*factor) / factor
}

func GetRatingsFromChain(c *gin.Context) {
    assetID := c.Query("asset_id")
    if assetID == "" {
        wrapError(c.JSON, "asset_id is required")
        return
    }

    avg, err := GetRatingFromChain(assetID)
    if err != nil {
        wrapError(c.JSON, err.Error())
        return
    }

    //wrapSuccess(c.JSON, fmt.Sprintf("average rating for asset %s: %.2f", assetID, avg))

    //roundedAvg := math.Round(avg*100) / 100
    roundedAvg := RoundToPrecision(avg, 2, 1)
// RoundToPrecision rounds a float64 to the specified number of decimal places.

    wrapSuccessJSON(c.JSON, map[string]interface{}{
        "asset_id": assetID,
        "average_rating": roundedAvg,
    })

}

// Error response wrapper
func wrapError(f func(code int, obj any), msg string) {
    f(404, gin.H{"message": msg})
}

// Success response wrapper
func wrapSuccess(f func(code int, obj any), msg string) {
   f(200, gin.H{"message": msg})
}

func wrapSuccessJSON(f func(code int, obj any), msg map[string]interface{}) {
    f(200, msg)
}

func main() {
    r := gin.Default()
    r.GET("/api/get_rating_by_asset", GetRatingsFromChain)
    r.Run(":8082")
}
