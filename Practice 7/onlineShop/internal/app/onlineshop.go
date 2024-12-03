package app

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"onlineShop/internal/auth"
	"onlineShop/internal/db"
	"onlineShop/internal/models"
	"os"
	"strconv"
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

	router.POST("/login", login)

	router.GET("/manufacturers", getManufacturer)
	router.GET("/manufacturers/:id", getManufacturer)
	router.GET("/manufacturers/:id/categories", getCategoriesOfManufacturer)

	router.GET("/products", getProduct)
	router.GET("/products/:id", getProduct)
	router.GET("/products/:id/params", getParametersOfProduct)

	router.GET("/receipt_types", getReceiptTypes)
	router.GET("/categories", getCategory)
	router.GET("/categories/:id", getCategory)
	router.GET("/categories/:id/params", getParameterOfCategory)

	protected := router.Group("/")
	protected.Use(authMiddleware())
	{
		protected.GET("/user_info", getUserInfo)
		protected.PUT("/user_info", updateUser)
		protected.GET("/orders", getClientOrders)
		protected.GET("/orders/:id", getClientOrders)
		protected.PUT("/orders/:id", updateOrder)
		protected.GET("/orders/:id/products", getProductsFromOrder)
		protected.PUT("/orders/:id/products", updateProductInOrder)
		protected.POST("/orders", addClientOrder)
	}
	router.Run(":8080")
	db.CloseDB()
}

func login(c *gin.Context) {
	var creds models.User
	if err := c.ShouldBind(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	conn, err := db.GetDB(&config)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	user, err := models.ReadUser(conn, creds.Email) //findUser(creds.Email)
	if err != nil || user.Password != creds.Password {
		log.Print(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Login error"})
	}
	token, err := auth.GenerateToken(creds.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create token"})
		return
	}
	//c.SetCookie("token", token, 3600, "/", "localhost", true, true)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		claims := &jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("my_secret_key"), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

//func findUser(email string) (*models.User, error) {
//	rows, err := executePreparedQuery("SELECT * FROM client WHERE email = $1", email)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//	if rows.Next() {
//		var user models.User
//		err := rows.Scan(&user.Id, &user.Name, &user.Phone, &user.Email, &user.Password)
//		if err != nil {
//			return nil, err
//		}
//		return &user, nil
//	}
//	return nil, sql.ErrNoRows
//}

func getManufacturer(c *gin.Context) {
	paramId := c.Param("id")
	var err error
	conn, err := db.GetDB(&config)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	var manufacturers []*models.Manufacturer
	if paramId == "" {
		manufacturers, err = models.ReadManufacturers(conn)
	} else {
		id, perr := strconv.Atoi(c.Param("id"))
		if perr != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		manufacturers, err = models.ReadManufacturerById(conn, id)
	}
	if err != nil {
		log.Print(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, manufacturers)
}

func getCategoriesOfManufacturer(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	conn, err := db.GetDB(&config)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	var categories []*models.Category
	categories, err = models.ReadCategoryOfManufacturer(conn, id)
	if err != nil {
		log.Print(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, categories)
}

func getProduct(c *gin.Context) {
	paramId := c.Param("id")
	var err error
	conn, err := db.GetDB(&config)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	var products []*models.Product
	if paramId == "" {
		products, err = models.ReadProducts(conn)
	} else {
		id, perr := strconv.Atoi(c.Param("id"))
		if perr != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		products, err = models.ReadProductById(conn, id)
	}
	if err != nil {
		log.Print(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, products)
}

func getParametersOfProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	conn, err := db.GetDB(&config)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	var params []*models.Parameter
	params, err = models.ReadParametersOfProduct(conn, id)
	if err != nil {
		log.Print(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, params)
}

func getCategory(c *gin.Context) {
	paramId := c.Param("id")
	var err error
	conn, err := db.GetDB(&config)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	var category []*models.Category
	if paramId == "" {
		category, err = models.ReadCategories(conn)
	} else {
		id, perr := strconv.Atoi(c.Param("id"))
		if perr != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		category, err = models.ReadCategoryById(conn, id)
	}
	if err != nil {
		log.Print(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, category)
}

func getParameterOfCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	conn, err := db.GetDB(&config)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	var parameters []*models.Parameter
	parameters, err = models.ReadParametersOfCategory(conn, id)
	//executePreparedQuery("SELECT * FROM parameter WHERE id = $1", id)
	if err != nil {
		log.Print(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, parameters)
}

func getUserInfo(c *gin.Context) {
	//cookie, err := c.Cookie("token")
	//if err != nil {
	//	log.Print(err)
	//	c.JSON(http.StatusUnauthorized, gin.H{"status": "Log in to proceed"})
	//	return
	//}
	tokenString := c.GetHeader("Authorization")
	user, err := getUser(tokenString)
	if err != nil {
		log.Print(err)
		c.Status(http.StatusNotFound)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, user)
}

func getUser(tokenString string) (*models.User, error) {
	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return auth.GetKey(), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	conn, err := db.GetDB(&config)
	if err != nil {
		return nil, err
	}
	return models.ReadUser(conn, claims.Subject)
}

func updateUser(c *gin.Context) {
	var newData models.User
	if err := c.BindJSON(&newData); err != nil {
		log.Print(err)
		c.Status(http.StatusBadRequest)
	}
	tokenString := c.GetHeader("Authorization")
	user, err := getUser(tokenString)
	if err != nil {
		log.Print(err)
		c.Status(http.StatusNotFound)
		c.Abort()
		return
	}
	conn, err := db.GetDB(&config)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	rows, err := models.UpdateUser(conn, &newData, user.Id)
	if err != nil {
		log.Print(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, fmt.Sprintf("Updated %d rows", rows))
}

func getReceiptTypes(c *gin.Context) {
	conn, err := db.GetDB(&config)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	var types []*models.Receipt
	types, err = models.ReadReceiptTypes(conn)
	if err != nil {
		log.Print(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, types)
}

func getClientOrders(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	user, err := getUser(tokenString)
	if err != nil {
		log.Print(err)
		c.Status(http.StatusNotFound)
		c.Abort()
		return
	}
	oid, _ := strconv.Atoi(c.Param("id"))
	conn, err := db.GetDB(&config)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	var clientOrders []*models.ClientOrder
	if oid > 0 {
		clientOrders, err = models.ReadClientOrderById(conn, user.Id, oid)
	} else {
		clientOrders, err = models.ReadClientOrders(conn, user.Id)
	}
	if err != nil {
		log.Print(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, clientOrders)
}

func updateOrder(c *gin.Context) {
	oid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Print(err)
		c.Status(http.StatusBadRequest)
	}
	var order models.ClientOrder
	if err := c.BindJSON(&order); err != nil {
		log.Print(err)
		c.Status(http.StatusBadRequest)
		return
	}
	tokenString := c.GetHeader("Authorization")
	user, err := getUser(tokenString)
	if err != nil {
		log.Print(err)
		c.Status(http.StatusNotFound)
		c.Abort()
		return
	}
	conn, err := db.GetDB(&config)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	rows, err := models.UpdateOrder(conn, oid, user.Id, &order)
	if err != nil {
		log.Print(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, fmt.Sprintf("Updated %d rows", rows))
}

func getProductsFromOrder(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	user, err := getUser(tokenString)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		c.Abort()
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Print(err)
		c.Status(http.StatusBadRequest)
	}
	conn, err := db.GetDB(&config)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	var products []*models.ProductInOrder
	products, err = models.ReadProductsFromOrder(conn, user.Id, id)
	if err != nil {
		log.Print(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, products)
}

func updateProductInOrder(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	user, err := getUser(tokenString)
	if err != nil {
		log.Print(err)
		c.Status(http.StatusNotFound)
		c.Abort()
		return
	}
	oid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Print(err)
		c.Status(http.StatusBadRequest)
		return
	}
	var product models.ProductInOrder
	if err := c.BindJSON(&product); err != nil {
		log.Print(err)
		c.Status(http.StatusBadRequest)
		return
	}
	conn, err := db.GetDB(&config)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	rows, err := models.UpdateProductInOrder(conn, product.Count, product.Product.Id, oid, user.Id)
	if err != nil {
		log.Print(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, fmt.Sprintf("Updated %d rows", rows))
}

func addClientOrder(c *gin.Context) {
	var order models.ClientOrder
	if err := c.BindJSON(&order); err != nil {
		log.Print(err)
		c.Status(http.StatusBadRequest)
		return
	}
	tokenString := c.GetHeader("Authorization")
	user, err := getUser(tokenString)
	if err != nil {
		log.Print(err)
		c.Status(http.StatusNotFound)
		c.Abort()
		return
	}
	var products []int
	var counts []int
	for _, product := range order.Products {
		products = append(products, product.Product.Id)
		counts = append(counts, product.Count)
	}
	conn, err := db.GetDB(&config)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	err = models.InsertOrder(conn, user.Id, order.Receipt.Id, products, counts)
	if err != nil {
		log.Print(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusCreated, "Created")
}
