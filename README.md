# go_go_cinema_app

App uses GIN, Gorm, Middleware(Auth and Logrus), Viper.

Set header "Authorization": "apiKey" for Authorization.

Methods:
1) POST:
reserve - http://localhost:8080/seat
buy - http://localhost:8080/buyseat
{
    "Seat": 13,
    "Status": "sold",
    "Movie": "Matrix",
    "Time": "10:00",
    "User": "b"
}

2) GET - http://localhost:8080/seat

DO: