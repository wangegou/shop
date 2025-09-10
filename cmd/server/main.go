package main

import (
	"auto-vending-system/internal/database"
	"auto-vending-system/internal/handlers"
	"auto-vending-system/internal/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func main() {
	database.Connect()
	createDefaultAdmin()

	r := gin.Default()
	r.LoadHTMLGlob("web/templates/*.html")
	store := cookie.NewStore([]byte("a-very-secret-key"))
	r.Use(sessions.Sessions("mysession", store))

	// --- Storefront HTML Pages ---
	r.GET("/", handlers.ShowStorefront)
	r.GET("/order/:id", handlers.ShowOrderPage)
	r.GET("/status/:order_no", handlers.ShowStatusPage)

	// --- Storefront API ---
	r.GET("/api/products", handlers.GetPublicProducts) // Renamed for clarity
	r.POST("/api/orders", handlers.CreateOrder)
	r.GET("/api/orders/:order_no", handlers.GetOrderStatus)
	r.GET("/api/coupons/verify", handlers.VerifyCoupon)


	// --- Admin HTML Pages ---
	adminPages := r.Group("/admin")
	{
		adminPages.GET("/login", handlers.ShowLoginPage)
		adminPages.GET("", func(c *gin.Context) {
			session := sessions.Default(c)
			if session.Get("userID") == nil {
				c.Redirect(http.StatusFound, "/admin/login")
			} else {
				c.Redirect(http.StatusFound, "/admin/dashboard")
			}
		})

		authPages := adminPages.Group("")
		authPages.Use(handlers.AuthRequired)
		{
			authPages.GET("/dashboard", handlers.ShowDashboard)
			authPages.GET("/categories", handlers.ShowCategoriesPage)
			authPages.GET("/categories/new", handlers.ShowNewCategoryForm)
			authPages.GET("/categories/edit/:id", handlers.ShowEditCategoryForm)
			authPages.GET("/products", handlers.ShowProductsPage)
			authPages.GET("/products/new", handlers.ShowNewProductForm)
			authPages.GET("/products/edit/:id", handlers.ShowEditProductForm)
			authPages.GET("/products/stock/:id", handlers.ShowStockPage)
			authPages.GET("/orders", handlers.ShowOrdersPage)
			authPages.GET("/coupons", handlers.ShowCouponsPage)
			authPages.GET("/coupons/new", handlers.ShowNewCouponForm)
			authPages.GET("/coupons/edit/:id", handlers.ShowEditCouponForm)
			authPages.GET("/settings", handlers.ShowSettingsPage)
		}
	}

	// --- Admin API endpoints ---
	adminApi := r.Group("/admin/api")
	{
		adminApi.POST("/login", handlers.Login)

		authApi := adminApi.Group("")
		authApi.Use(handlers.AuthRequired)
		{
			authApi.POST("/logout", handlers.Logout)
			authApi.POST("/categories", handlers.CreateCategory)
			authApi.GET("/categories", handlers.GetCategories)
			authApi.GET("/categories/:id", handlers.GetCategory)
			authApi.PUT("/categories/:id", handlers.UpdateCategory)
			authApi.DELETE("/categories/:id", handlers.DeleteCategory)
			authApi.GET("/products", handlers.GetProducts)
			authApi.POST("/products", handlers.CreateProduct)
			authApi.GET("/products/:id", handlers.GetProduct)
			authApi.PUT("/products/:id", handlers.UpdateProduct)
			authApi.DELETE("/products/:id", handlers.DeleteProduct)
			authApi.GET("/products/:id/stock", handlers.GetStockForProduct)
			authApi.POST("/products/:id/stock", handlers.AddStock)
			authApi.GET("/orders", handlers.GetOrders)
			authApi.GET("/orders/:id", handlers.GetOrder)
			authApi.GET("/coupons", handlers.GetCoupons)
			authApi.POST("/coupons", handlers.CreateCoupon)
			authApi.GET("/coupons/:id", handlers.GetCoupon)
			authApi.PUT("/coupons/:id", handlers.UpdateCoupon)
			authApi.DELETE("/coupons/:id", handlers.DeleteCoupon)
			authApi.GET("/settings", handlers.GetSettings)
			authApi.POST("/settings", handlers.UpdateSettings)
			authApi.POST("/email/test", handlers.EmailTest)
		}
	}

	r.Run(":8080")
}

func createDefaultAdmin() {
	var userCount int64
	database.DB.Model(&models.User{}).Count(&userCount)
	if userCount == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		if err != nil {
			panic("Failed to hash password")
		}
		admin := models.User{Username: "admin", Password: string(hashedPassword)}
		database.DB.Create(&admin)
	}
}
