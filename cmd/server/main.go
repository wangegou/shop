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
	seedDatabase() // Add seed data

	r := gin.Default()
	r.LoadHTMLGlob("web/templates/*.html")
	store := cookie.NewStore([]byte("a-very-secret-key"))
	r.Use(sessions.Sessions("mysession", store))

	// --- Storefront HTML Pages ---
	r.GET("/", handlers.ShowStorefront)
	r.GET("/order/:id", handlers.ShowOrderPage)
	r.GET("/status/:order_no", handlers.ShowStatusPage)

	// --- Storefront API ---
	r.GET("/api/products", handlers.GetPublicProducts)
	r.POST("/api/orders", handlers.CreateOrder)
	r.GET("/api/orders/:order_no", handlers.GetOrderStatus)
	r.GET("/api/coupons/verify", handlers.VerifyCoupon)

	// --- Payment Flow ---
	r.GET("/payment/redirect/:order_no", handlers.RedirectToPayment)
	r.GET("/payment/notify", handlers.NotifyHandler)

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

func seedDatabase() {
	var productCount int64
	database.DB.Model(&models.Product{}).Count(&productCount)
	if productCount > 0 {
		return // Database already seeded
	}

	// Seed Categories
	cat1 := models.Category{Name: "Software Licenses"}
	cat2 := models.Category{Name: "Game Keys"}
	database.DB.Create(&cat1)
	database.DB.Create(&cat2)

	// Seed Products
	prod1 := models.Product{Name: "IDE Pro License", Price: 99.99, CategoryID: cat1.ID, Description: "A professional IDE for developers."}
	prod2 := models.Product{Name: "OS License Key", Price: 129.50, CategoryID: cat1.ID, Description: "The latest version of the OS."}
	prod3 := models.Product{Name: "Indie Adventure Game", Price: 19.99, CategoryID: cat2.ID, Description: "An exciting adventure game."}
	database.DB.Create(&prod1)
	database.DB.Create(&prod2)
	database.DB.Create(&prod3)

	// Seed Stock
	stock1 := []models.ProductStock{
		{ProductID: prod1.ID, Key: "IDE-PRO-KEY-1111-AAAA"},
		{ProductID: prod1.ID, Key: "IDE-PRO-KEY-2222-BBBB"},
	}
	stock2 := []models.ProductStock{
		{ProductID: prod2.ID, Key: "OS-LICENSE-KEY-3333-CCCC"},
	}
	stock3 := []models.ProductStock{
		{ProductID: prod3.ID, Key: "INDIE-GAME-KEY-4444-DDDD"},
		{ProductID: prod3.ID, Key: "INDIE-GAME-KEY-5555-EEEE"},
		{ProductID: prod3.ID, Key: "INDIE-GAME-KEY-6666-FFFF"},
	}
	database.DB.Create(&stock1)
	database.DB.Create(&stock2)
	database.DB.Create(&stock3)
}
