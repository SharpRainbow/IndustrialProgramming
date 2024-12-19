package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"onlineShop/internal/auth"
	"onlineShop/internal/db"
	"onlineShop/internal/models"
	"os"
	"strconv"
	"time"
)

var config = db.Config{DbHost: "localhost", DbPort: "5432", DbName: "online_shop", DbUser: "postgres", DbPassword: "1234"}

func Run() {
	var host = os.Getenv("HOSTNAME")
	var port = os.Getenv("PORT")
	var dbName = os.Getenv("DATABASE_NAME")
	var dbUser = os.Getenv("DATABASE_USER")
	var dbPassword = os.Getenv("DATABASE_PASSWORD")

	if host != "" && port != "" && dbName != "" && dbUser != "" && dbPassword != "" {
		config = db.Config{DbHost: host, DbPort: port, DbName: dbName, DbUser: dbUser, DbPassword: dbPassword}
	}
	log.Printf("Using config: %s", config)
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow requests from all origins or specific ones
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	version1 := router.Group("/api/v1")

	version1.POST("/login", login)
	version1.POST("/register", registerUser)

	protected := version1.Group("/")
	protected.Use(authMiddleware())
	{
		protected.GET("/manufacturers", getManufacturer)
		protected.GET("/manufacturers/:id", getManufacturer)
		protected.GET("/manufacturers/:id/categories", getCategoriesOfManufacturer)

		protected.GET("/products", getProduct)
		protected.GET("/products/:id", getProduct)
		protected.GET("/products/:id/params", getParametersOfProduct)

		protected.GET("/receipt_types", getReceiptTypes)
		protected.GET("/categories", getCategory)
		protected.GET("/categories/:id", getCategory)
		protected.GET("/categories/:id/params", getParameterOfCategory)

		protected.GET("/user_info", getUserInfo)
		protected.PUT("/user_info", updateUser)
		protected.GET("/orders", getClientOrders)
		protected.GET("/orders/:id", getClientOrders)
		protected.PUT("/orders/:id", updateOrder)
		protected.GET("/orders/:id/products", getProductsFromOrder)
		protected.PUT("/orders/:id/products", updateProductInOrder)
		protected.POST("/orders", addClientOrder)

		moderation := protected.Group("/")
		moderation.Use(roleMiddleware("admin"))
		{
			moderation.GET("/users", getUsers)
			moderation.GET("/users/:id", getUsersById)
			moderation.DELETE("/users/:id", removeUser)
		}
	}
	router.Run(":8080")
	db.CloseMapper()
}

func handleError(c *gin.Context, statusCode int, err error, message string) {
	log.Println(err)
	c.JSON(statusCode, gin.H{"error": message})
}

// @Summary Login a user
// @Description Authenticates a user using email and password, and returns a JWT token upon successful authentication.
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.UserLogin true "User credentials (Only `Email` and `Password` fields are used)"
// @Success 200 {object} map[string]string "JWT token returned upon successful login"
// @Failure 400 {object} models.ErrorResponse "Invalid credentials provided"
// @Failure 404 {object} models.ErrorResponse "User not found or incorrect password"
// @Failure 500 {object} models.ErrorResponse "Server malfunction, try later"
// @Router /login [post]
func login(c *gin.Context) {
	var creds models.User
	if err := c.ShouldBind(&creds); err != nil {
		handleError(c, http.StatusBadRequest, err, "invalid credentials")
		return
	}
	conn, err := db.GetMapper(&config)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, "server malfunction, try later")
		return
	}
	toutCtx, cancel := context.WithTimeout(c.Request.Context(), time.Second*2)
	defer cancel()
	conn = conn.WithContext(toutCtx)
	user, err := models.ReadUserRole(conn, creds.Email)
	if err != nil || user.Password != creds.Password {
		handleError(c, http.StatusNotFound, err, "user not found")
		return
	}
	token, err := auth.GenerateToken(creds.Email, user.Role)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, "server malfunction, try later")
		return
	}
	//c.SetCookie("token", token, 3600, "/", "localhost", true, true)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// @Summary Register a new user
// @Description Creates a new user account with the provided credentials (email and password).
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.User true "User credentials (Only `Name`, `Phone`, `Email`, `Password` fields are used)"
// @Success 200 {string} string "User successfully created"
// @Failure 400 {object} models.ErrorResponse "Invalid credentials provided"
// @Failure 500 {object} models.ErrorResponse "User creation failed or server malfunction"
// @Router /register [post]
func registerUser(c *gin.Context) {
	var creds models.User
	if err := c.ShouldBind(&creds); err != nil {
		handleError(c, http.StatusBadRequest, err, "invalid credentials")
		return
	}
	conn, err := db.GetMapper(&config)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, "server malfunction, try later")
		return
	}
	toutCtx, cancel := context.WithTimeout(c.Request.Context(), time.Second*2)
	defer cancel()
	conn = conn.WithContext(toutCtx)
	userId, err := models.AddUser(conn, &creds)
	if err != nil || userId != 1 {
		handleError(c, http.StatusInternalServerError, err, "user creation failed")
		return
	}
	c.JSON(http.StatusOK, "user successfully created")
}

// @Summary Retrieve specified users
// @Description Fetch user from the database by ID.
// @Tags users
// @Produce json
// @Param Authorization header string true "Admin token"
// @Param id path int true "User ID"
// @Success 200 {object} models.User "Requested user info"
// @Failure 400 {object} models.ErrorResponse "Invalid id provided"
// @Failure 500 {object} models.ErrorResponse "Server malfunction or users selection failed"
// @Router /users/{id} [get]
func getUsersById(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Param("id"))
	if uid <= 0 {
		handleError(c, http.StatusBadRequest, errors.New("invalid user id"), "invalid user id")
		return
	}
	conn, err := db.GetMapper(&config)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, "server malfunction, try later")
		return
	}
	toutCtx, cancel := context.WithTimeout(c.Request.Context(), time.Second*2)
	defer cancel()
	conn = conn.WithContext(toutCtx)
	rows, err := models.ReadUserById(conn, uid)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, "users selection failed")
		return
	}
	c.JSON(http.StatusOK, rows)
}

// @Summary Retrieve all users
// @Description Fetch a list of all users from the database.
// @Tags users
// @Produce json
// @Param Authorization header string true "Admin token"
// @Param page query int false "Page number for pagination" default(1)
// @Param limit query int false "Items per page for pagination" default(5)
// @Param sort query string false "Sort field" default(id)
// @Param email query string false "Search by user email"
// @Success 200 {array} models.User "List of all users"
// @Failure 400 {object} models.ErrorResponse "Wrong sort parameter"
// @Failure 500 {object} models.ErrorResponse "Server malfunction or users selection failed"
// @Router /users [get]
func getUsers(c *gin.Context) {
	conn, err := db.GetMapper(&config)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, "server malfunction, try later")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
	sort := c.DefaultQuery("sort", "id")
	if !models.UsersSortFieldValid(sort) {
		handleError(c, http.StatusBadRequest, errors.New("wrong sort param"), "unknown sort field")
		return
	}
	email := c.DefaultQuery("email", "")
	offset := (page - 1) * limit
	toutCtx, cancel := context.WithTimeout(c.Request.Context(), time.Second*2)
	defer cancel()
	conn = conn.WithContext(toutCtx)
	rows, err := models.ReadAllUsers(conn, limit, offset, email, sort)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, "users selection failed")
		return
	}
	c.JSON(http.StatusOK, rows)
}

// @Summary Remove a user
// @Description Deletes a user from the database by their ID.
// @Tags users
// @Param Authorization header string true "Admin token"
// @Param id path int true "User ID to be deleted"
// @Success 200 {string} string "User successfully deleted"
// @Failure 400 {object} models.ErrorResponse "Invalid User ID format"
// @Failure 500 {object} models.ErrorResponse "Server malfunction or user deletion failed"
// @Router /users/{id} [delete]
func removeUser(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("id"))
	conn, err := db.GetMapper(&config)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, "server malfunction, try later")
		return
	}
	toutCtx, cancel := context.WithTimeout(c.Request.Context(), time.Second*2)
	defer cancel()
	conn = conn.WithContext(toutCtx)
	rows, err := models.RemoveUser(conn, userId)
	if err != nil || rows != 1 {
		handleError(c, http.StatusInternalServerError, err, "user deletion failed")
		return
	}
	c.JSON(http.StatusOK, "user successfully deleted")
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		claims := &jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("my_secret_key"), nil
		})
		if err != nil || !token.Valid {
			handleError(c, http.StatusUnauthorized, err, "token validation error")
			c.Abort()
			return
		}
		c.Set("role", claims.Audience)
		c.Next()
	}
}

func roleMiddleware(expectedRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != expectedRole {
			handleError(c, http.StatusForbidden, errors.New("invalid role"), "access denied")
			c.Abort()
			return
		}
		c.Next()
	}
}

// @Summary Get manufacturer
// @Description Get manufacturer info by ID
// @Tags manufacturers
// @Param Authorization header string true "User token"
// @Param id path int true "Manufacturer ID"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Manufacturer
// @Router /manufacturers/{id} [get]
func getManufacturerById() {

}

// @Summary Get list of manufacturers
// @Description Get all available manufacturers
// @Tags manufacturers
// @Param Authorization header string true "User token"
// @Param page query int false "Page number for pagination" default(1)
// @Param limit query int false "Items per page for pagination" default(5)
// @Param sort query string false "Sort field" default(id)
// @Param name query string false "Search by manufacturer name"
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Manufacturer
// @Router /manufacturers [get]
func getManufacturer(c *gin.Context) {
	paramId := c.Param("id")
	var err error
	conn, err := db.GetMapper(&config)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, "server malfunction, try later")
		return
	}
	toutCtx, cancel := context.WithTimeout(c.Request.Context(), time.Second*2)
	defer cancel()
	conn = conn.WithContext(toutCtx)
	var manufacturers *[]models.Manufacturer
	if paramId == "" {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
		sort := c.DefaultQuery("sort", "id")
		if !models.ManufacturerSortFieldValid(sort) {
			handleError(c, http.StatusBadRequest, errors.New("wrong sort param"), "unknown sort field")
			return
		}
		name := c.DefaultQuery("name", "")
		offset := (page - 1) * limit
		manufacturers, err = models.ReadManufacturers(conn, limit, offset, name, sort)
	} else {
		id, perr := strconv.Atoi(c.Param("id"))
		if perr != nil {
			handleError(c, http.StatusBadRequest, perr, "wrong id")
			return
		}
		manufacturers, err = models.ReadManufacturerById(conn, id)
	}
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, "server malfunction, try later")
		return
	}
	c.JSON(http.StatusOK, manufacturers)
}

// @Summary Get categories of a manufacturer by ID
// @Description Fetches a list of categories associated with a specific manufacturer ID.
// @Tags manufacturers
// @Param Authorization header string true "User token"
// @Param id path int true "Manufacturer ID"
// @Success 200 {array} models.Category "List of categories for the manufacturer"
// @Failure 400 {object} models.ErrorResponse "Invalid manufacturer ID format"
// @Failure 500 {object} models.ErrorResponse "Server malfunction or unable to fetch categories"
// @Router /manufacturers/{id}/categories [get]
func getCategoriesOfManufacturer(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handleError(c, http.StatusBadRequest, err, "wrong id")
		return
	}
	conn, err := db.GetMapper(&config)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, "server malfunction, try later")
		return
	}
	toutCtx, cancel := context.WithTimeout(c.Request.Context(), time.Second*2)
	defer cancel()
	conn = conn.WithContext(toutCtx)
	var categories *[]models.Category
	categories, err = models.ReadCategoryOfManufacturer(conn, id)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, "server malfunction, try later")
		return
	}
	c.JSON(http.StatusOK, categories)
}

//func updateManufacturers(c *gin.Context) {
//	id, _ := strconv.Atoi(c.Param("id"))
//	//conn, err := db.GetMapper(&config)
//	//if err != nil {
//	//	handleError(c, http.StatusInternalServerError, err, "server malfunction, try later")
//	//	return
//	//}
//	var manufacturer models.Manufacturer
//	if err := c.BindJSON(&manufacturer); err != nil {
//		handleError(c, http.StatusBadRequest, err, "wrong user data format")
//		return
//	}
//	if id <= 0 {
//
//	}
//	c.JSON(http.StatusOK, manufacturer)
//}

// @Summary Get product indo
// @Description Get product by id
// @Tags products
// @Param Authorization header string true "User token"
// @Param id path int false "Product ID"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Product
// @Router /products/{id} [get]
func getProductById() {

}

// @Summary Get list of products
// @Description Get all available products
// @Tags products
// @Param Authorization header string true "User token"
// @Param page query int false "Page number for pagination" default(1)
// @Param limit query int false "Items per page for pagination" default(5)
// @Param sort query string false "Sort field" default(id)
// @Param name query string false "Search by product name"
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Product
// @Router /products [get]
func getProduct(c *gin.Context) {
	paramId := c.Param("id")
	conn, err := db.GetMapper(&config)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, "server malfunction, try later")
		return
	}
	toutCtx, cancel := context.WithTimeout(c.Request.Context(), time.Second*2)
	defer cancel()
	conn = conn.WithContext(toutCtx)
	var products *[]models.Product
	if paramId == "" {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
		name := c.DefaultQuery("name", "")
		sort := c.DefaultQuery("sort", "id")
		if !models.ProductSortFieldValid(sort) {
			handleError(c, http.StatusBadRequest, errors.New("wrong sort param"), "unknown sort field")
			return
		}
		offset := (page - 1) * limit
		products, err = models.ReadProducts(conn, limit, offset, name, sort)
	} else {
		id, perr := strconv.Atoi(c.Param("id"))
		if perr != nil {
			handleError(c, http.StatusBadRequest, perr, "wrong id")
			return
		}
		products, err = models.ReadProductById(conn, id)
	}
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, "server malfunction, try later")
		return
	}
	c.JSON(http.StatusOK, products)
}

// @Summary Get parameters of a product by ID
// @Description Fetches a list of parameters associated with a specific product ID.
// @Tags products
// @Param Authorization header string true "User token"
// @Param id path int true "Product ID"
// @Success 200 {array} models.Parameter "List of parameters for the product"
// @Failure 400 {object} models.ErrorResponse "Invalid product ID format"
// @Failure 500 {object} models.ErrorResponse "Server malfunction or unable to fetch parameters"
// @Router /products/{id}/params [get]
func getParametersOfProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handleError(c, http.StatusBadRequest, err, "wrong id")
		return
	}
	conn, err := db.GetMapper(&config)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, "server malfunction, try later")
		return
	}
	toutCtx, cancel := context.WithTimeout(c.Request.Context(), time.Second*2)
	defer cancel()
	conn = conn.WithContext(toutCtx)
	var params *[]models.Parameter
	params, err = models.ReadParametersOfProduct(conn, id)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, "server malfunction, try later")
		return
	}
	c.JSON(http.StatusOK, params)
}

// @Summary Get category details
// @Description Fetches a single category by ID
// @Tags categories
// @Param Authorization header string true "User token"
// @Param id path int optional "Category ID"
// @Success 200 {object} models.Category "List of categories or single category"
// @Failure 400 {object} models.ErrorResponse "Invalid ID format or query parameters"
// @Failure 500 {object} models.ErrorResponse "Server malfunction or unable to fetch category"
// @Router /categories/{id} [get]
func getCategoryById() {

}

// @Summary Get category details or a list of categories
// @Description Fetches a single category by ID or paginated list of categories based on query parameters.
// @Tags categories
// @Param Authorization header string true "User token"
// @Param page query int false "Page number for pagination" default(1)
// @Param limit query int false "Items per page for pagination" default(5)
// @Success 200 {array} models.Category "List of categories or single category"
// @Failure 400 {object} models.ErrorResponse "Invalid ID format or query parameters"
// @Failure 500 {object} models.ErrorResponse "Server malfunction or unable to fetch category"
// @Router /categories [get]
func getCategory(c *gin.Context) {
	paramId := c.Param("id")
	var err error
	conn, err := db.GetMapper(&config)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, "server malfunction, try later")
		return
	}
	toutCtx, cancel := context.WithTimeout(c.Request.Context(), time.Second*2)
	defer cancel()
	conn = conn.WithContext(toutCtx)
	var category *[]models.Category
	if paramId == "" {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
		offset := (page - 1) * limit
		category, err = models.ReadCategories(conn, limit, offset)
	} else {
		id, perr := strconv.Atoi(c.Param("id"))
		if perr != nil {
			handleError(c, http.StatusBadRequest, perr, "wrong id")
			return
		}
		category, err = models.ReadCategoryById(conn, id)
	}
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, "server malfunction, try later")
		return
	}
	c.JSON(http.StatusOK, category)
}

// @Summary Get parameters of a category by ID
// @Description Fetches a list of parameters associated with a specific category ID.
// @Tags categories
// @Param Authorization header string true "User token"
// @Param id path int true "Category ID"
// @Success 200 {array} models.Parameter "List of parameters for the category"
// @Failure 400 {object} models.ErrorResponse "Invalid category ID format"
// @Failure 500 {object} models.ErrorResponse "Server malfunction or unable to fetch parameters"
// @Router /categories/{id}/params [get]
func getParameterOfCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handleError(c, http.StatusBadRequest, err, "wrong id")
		return
	}
	conn, err := db.GetMapper(&config)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, "server malfunction, try later")
		return
	}
	toutCtx, cancel := context.WithTimeout(c.Request.Context(), time.Second*2)
	defer cancel()
	conn = conn.WithContext(toutCtx)
	var parameters *[]models.Parameter
	parameters, err = models.ReadParametersOfCategory(conn, id)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, "server malfunction, try later")
		return
	}
	c.JSON(http.StatusOK, parameters)
}

// @Summary Get user information from the provided authorization token
// @Description Fetches user details using the provided Authorization header token.
// @Tags users
// @Param Authorization header string true "User token"
// @Success 200 {object} models.User "User information"
// @Failure 401 {object} models.ErrorResponse "Unauthorized: Invalid or missing token"
// @Failure 500 {object} models.ErrorResponse "Server malfunction, unable to retrieve user info"
// @Router /user_info [get]
func getUserInfo(c *gin.Context) {
	//cookie, err := c.Cookie("token")
	//if err != nil {
	//	log.Print(err)
	//	c.JSON(http.StatusUnauthorized, gin.H{"status": "Log in to proceed"})
	//	return
	//}
	tokenString := c.GetHeader("Authorization")
	user, err := getUser(tokenString, c.Request.Context())
	if err != nil {
		handleError(c, http.StatusUnauthorized, err, "check authorization")
		return
	}
	c.JSON(http.StatusOK, user)
}

func getUser(tokenString string, ctx context.Context) (*models.User, error) {
	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return auth.GetKey(), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	conn, err := db.GetMapper(&config)
	if err != nil {
		return nil, err
	}
	toutCtx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	conn = conn.WithContext(toutCtx)
	return models.ReadUser(conn, claims.Subject)
}

// @Summary Update user information
// @Description Updates the user information using the provided authorization token and new user data.
// @Tags users
// @Param Authorization header string true "User token"
// @Param request body models.User true "Updated user information (Only `Name`, `Phone`, `Password` fields are used)"
// @Success 200 {string} string "Updated {number} rows"
// @Failure 400 {object} models.ErrorResponse "Invalid user data format"
// @Failure 401 {object} models.ErrorResponse "Unauthorized: Invalid or missing token"
// @Failure 500 {object} models.ErrorResponse "Server malfunction, unable to update user"
// @Router /user_info [put]
func updateUser(c *gin.Context) {
	var newData models.User
	if err := c.BindJSON(&newData); err != nil {
		handleError(c, http.StatusBadRequest, err, "wrong user data format")
		return
	}
	tokenString := c.GetHeader("Authorization")
	user, err := getUser(tokenString, c.Request.Context())
	if err != nil {
		handleError(c, http.StatusUnauthorized, err, "check authorization")
		return
	}
	conn, err := db.GetMapper(&config)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, "server malfunction, try later")
		return
	}
	toutCtx, cancel := context.WithTimeout(c.Request.Context(), time.Second*2)
	defer cancel()
	conn = conn.WithContext(toutCtx)
	rows, err := models.UpdateUser(conn, &newData, user.Id)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, "server malfunction, try later")
		return
	}
	c.JSON(http.StatusOK, fmt.Sprintf("updated %d rows", rows))
}

// @Summary Get all receipt types
// @Description Fetches a list of all receipt types from the database.
// @Tags receipts
// @Param Authorization header string true "User token"
// @Success 200 {array} models.Receipt "List of receipt types"
// @Failure 500 {object} models.ErrorResponse "Server malfunction, unable to fetch receipt types"
// @Router /receipt_types [get]
func getReceiptTypes(c *gin.Context) {
	conn, err := db.GetMapper(&config)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, "server malfunction, try later")
		return
	}
	toutCtx, cancel := context.WithTimeout(c.Request.Context(), time.Second*2)
	defer cancel()
	conn = conn.WithContext(toutCtx)
	var types *[]models.Receipt
	types, err = models.ReadReceiptTypes(conn)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, "server malfunction, try later")
		return
	}
	c.JSON(http.StatusOK, types)
}

// @Summary Get specific order
// @Description Get order by ID
// @Tags orders
// @Accept  json
// @Produce  json
// @Param Authorization header string true "User token"
// @Param id path int true "Order ID"
// @Success 200 {object} models.ClientOrder
// @Router /orders/{id} [get]
func getOrderById() {

}

// @Summary Get list of client orders
// @Description Get all orders
// @Tags orders
// @Accept  json
// @Produce  json
// @Param Authorization header string true "User token"
// @Param page query int false "Page number for pagination" default(1)
// @Param limit query int false "Items per page for pagination" default(5)
// @Param sort query string false "Sort field" default(id)
// @Success 200 {array} models.ClientOrder
// @Failure 400 {object} models.ErrorResponse "Invalid sort param"
// @Failure 500 {object} models.ErrorResponse "Server malfunction or unable to fetch categories"
// @Router /orders [get]
func getClientOrders(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	user, err := getUser(tokenString, c.Request.Context())
	if err != nil {
		handleError(c, http.StatusUnauthorized, err, "check authorization")
		return
	}
	oid, _ := strconv.Atoi(c.Param("id"))
	conn, err := db.GetMapper(&config)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, "server malfunction, try later")
		return
	}
	toutCtx, cancel := context.WithTimeout(c.Request.Context(), time.Second*2)
	defer cancel()
	conn = conn.WithContext(toutCtx)
	var clientOrders *[]models.ClientOrder
	if oid > 0 {
		clientOrders, err = models.ReadClientOrderById(conn, user.Id, oid)
	} else {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
		sort := c.DefaultQuery("sort", "id")
		if !models.OrderSortFieldValid(sort) {
			handleError(c, http.StatusBadRequest, errors.New("wrong sort param"), "unknown sort field")
			return
		}
		offset := (page - 1) * limit
		clientOrders, err = models.ReadClientOrders(conn, user.Id, limit, offset, sort)
	}
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, "server malfunction, try later")
		return
	}
	c.JSON(http.StatusOK, clientOrders)
}

// @Summary Update an existing client order
// @Description Updates the details of an existing client order identified by the order ID. Requires the user to be authenticated via an Authorization token.
// @Tags orders
// @Param Authorization header string true "User token"
// @Param id path int true "Order ID"
// @Param request body models.ClientOrderUpdate true "Updated client order details"
// @Success 200 {string} string "Updated {number} rows"
// @Failure 400 {object} models.ErrorResponse "Wrong ID format or wrong order data format"
// @Failure 401 {object} models.ErrorResponse "Unauthorized: Invalid or missing token"
// @Failure 500 {object} models.ErrorResponse "Server malfunction, unable to update order"
// @Router /orders/{id} [put]
func updateOrder(c *gin.Context) {
	oid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handleError(c, http.StatusBadRequest, err, "wrong id")
	}
	var order models.ClientOrderUpdate
	if err := c.BindJSON(&order); err != nil {
		handleError(c, http.StatusBadRequest, err, "wrong order data format")
		return
	}
	tokenString := c.GetHeader("Authorization")
	user, err := getUser(tokenString, c.Request.Context())
	if err != nil {
		handleError(c, http.StatusNotFound, err, "check authorization")
		return
	}
	conn, err := db.GetMapper(&config)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, "server malfunction, try later")
		return
	}
	toutCtx, cancel := context.WithTimeout(c.Request.Context(), time.Second*2)
	defer cancel()
	conn = conn.WithContext(toutCtx)
	rows, err := models.UpdateClientOrder(conn, oid, user.Id, &order)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, "server malfunction, try later")
		return
	}
	c.JSON(http.StatusOK, fmt.Sprintf("updated %d rows", rows))
}

// @Summary Get products from a specific order
// @Description Retrieves a list of products associated with a specific order, requiring the user to be authenticated via an Authorization token.
// @Tags orders
// @Param Authorization header string true "User token"
// @Param id path int true "Order ID"
// @Success 200 {array} models.ProductInOrder "List of products in the order"
// @Failure 400 {object} models.ErrorResponse "Wrong ID format"
// @Failure 401 {object} models.ErrorResponse "Unauthorized: Invalid or missing token"
// @Failure 500 {object} models.ErrorResponse "Server malfunction, unable to fetch products from order"
// @Router /orders/{id}/products [get]
func getProductsFromOrder(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	user, err := getUser(tokenString, c.Request.Context())
	if err != nil {
		handleError(c, http.StatusUnauthorized, err, "check authorization")
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handleError(c, http.StatusBadRequest, err, "wrong id")
		return
	}
	conn, err := db.GetMapper(&config)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, "server malfunction, try later")
		return
	}
	toutCtx, cancel := context.WithTimeout(c.Request.Context(), time.Second*2)
	defer cancel()
	conn = conn.WithContext(toutCtx)
	var products *[]models.ProductInOrder
	products, err = models.ReadProductsFromOrder(conn, user.Id, id)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, "server malfunction, try later")
		return
	}
	c.JSON(http.StatusOK, products)
}

// @Summary Update a product in an order
// @Description Updates the quantity of a product in a specific order, requiring the user to be authenticated via an Authorization token.
// @Tags orders
// @Param Authorization header string true "User token"
// @Param id path int true "Order ID"
// @Param product body models.ProductInOrderUpdate true "ProductInOrder object"
// @Success 200 {string} string "Updated rows count message"
// @Failure 400 {object} models.ErrorResponse "Invalid product data format or wrong ID"
// @Failure 401 {object} models.ErrorResponse "Unauthorized: Invalid or missing token"
// @Failure 500 {object} models.ErrorResponse "Server malfunction, unable to update product in order"
// @Router /orders/{id}/products [put]
func updateProductInOrder(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	user, err := getUser(tokenString, c.Request.Context())
	if err != nil {
		handleError(c, http.StatusUnauthorized, err, "check authorization")
		return
	}
	oid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handleError(c, http.StatusBadRequest, err, "wrong id")
		return
	}
	var product models.ProductInOrderUpdate
	if err := c.BindJSON(&product); err != nil {
		handleError(c, http.StatusBadRequest, err, "wrong product data format")
		return
	}
	conn, err := db.GetMapper(&config)
	log.Println(product)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, "server malfunction, try later")
		return
	}
	toutCtx, cancel := context.WithTimeout(c.Request.Context(), time.Second*2)
	defer cancel()
	conn = conn.WithContext(toutCtx)
	rows, err := models.UpdateProductInOrder(conn, product.Count, product.ProductId, oid, user.Id)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, "server malfunction, try later")
		return
	}
	c.JSON(http.StatusOK, fmt.Sprintf("updated %d rows", rows))
}

// @Summary Add a new client order
// @Description Create a new order for the authenticated user with a list of products.
// @Tags orders
// @Accept json
// @Produce json
// @Param Authorization header string true "User token"
// @Param order body models.ClientOrderUpdate true "Order data"
// @Success 201 {string} string "order successfully created"
// @Failure 400 {object} models.ErrorResponse "Bad Request: wrong order data format or no products provided"
// @Failure 401 {object} models.ErrorResponse "Unauthorized: check authorization"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error: server malfunction, try later"
// @Router /orders [post]
func addClientOrder(c *gin.Context) {
	var order models.ClientOrderUpdate
	if err := c.BindJSON(&order); err != nil {
		handleError(c, http.StatusBadRequest, err, "wrong order data format")
		return
	}
	if len(order.Products) <= 0 {
		handleError(c, http.StatusBadRequest, errors.New("no products list"), "no products list provided")
		return
	}
	tokenString := c.GetHeader("Authorization")
	user, err := getUser(tokenString, c.Request.Context())
	if err != nil {
		handleError(c, http.StatusUnauthorized, err, "check authorization")
		return
	}
	var products []int
	var counts []int
	for _, product := range order.Products {
		products = append(products, product.ProductId)
		counts = append(counts, product.Count)
	}
	conn, err := db.GetMapper(&config)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, "server malfunction, try later")
		return
	}
	toutCtx, cancel := context.WithTimeout(c.Request.Context(), time.Second*2)
	defer cancel()
	conn = conn.WithContext(toutCtx)
	err = models.InsertOrder(conn, user.Id, order.ReceiptId, products, counts)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, "server malfunction, try later")
		return
	}
	c.JSON(http.StatusCreated, "order successfully created")
}
