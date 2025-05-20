# go_microservice_template

App uses GIN, Gorm, Middleware(Auth and Logrus), Viper.

Set header "Authorization": "apiKey" for Authorization.

Methods:
1) POST - http://localhost:8080/seat
{
    "Name":"Epic",
    "Email": "a@a.kz"
}

2) GET - http://localhost:8080/seat