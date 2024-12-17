# Microservice_Recomendtaions
Microservice Recomendations for users

Некоторые заметки для построения архитектуры проекта
создание юзера, обновление происходят в обычном режиме
сервис юзеров и product service : структура юзера, структура корзины для пользователя, и структура обеъктов этой кор : 

 ## type User struct {
   Username      string        // Логин пользователя
   Password_hash string        // Пароль пользователя который будет хэшироваться
   Email         string        // Электронная потча юзера
   CreatedAt time.Time         // дата создания 
   UpdatedAt time.Time         // дата обновления 
 }

// Cart представляет корзину пользователя
## type Cart struct {
###	UserID     string        // Идентификатор пользователя
	Items      []CartItem    // Список продуктов в корзине
	TotalPrice float64       // Общая стоимость корзины
	CreatedAt  time.Time     // Дата и время создания корзины
	UpdatedAt  time.Time     // Дата и время последнего обновления корзины
## }

// Объект корзины, то есть сам продукт
## type CartItem struct {
###	ProductID   string      // Идентификатор продукта
	ProductName string      // Название продукта
	Quantity    int         // Количество продукта
	Price       float64     // Цена продукта
## }
