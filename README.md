# gym_consistency
WIP!

## To-do
* Create ping endpoint to check if the application is Healty
* Create user and authentication
* A lot more... 

## Stack
* Postgres
* Golang
* Docker (optional)

## How to execute

1. You need to execute the init_script file on sql folder.
2. On main.go you can see that "conn" needs an environment variable so to initializate the program you need to pass this variable. 
The string requires something like this: 
> DBCONN="host=myhost port=myport user=myuser password=mypassword dbname=postgres"
With this informations you can initializate the program. 

## Enviroment Variables
Some enviroment variables are needed to run the API succefully.
* DBCONN -> Database string connection
  * Ex: DBCONN="host=1 port=2 user=3 password=4 dbname=5"
  * 1 = Gateway, IP host or url where your postgresql is hosted
  * 2 = The port number postgresql is available
  * 3 = Your postgresql user
  * 4 = The password defined for your postgresql user
  * 5 = The database name defined in  postgresl
* SECRET_KEY_JWT -> Secret key responsible for encoding e decode JWT
  * Ex: SECRET_KEY_JWT="abcde12345"
* EMAIL_GYM -> Valid e-mail address responsible for send e-mails for users
  * Ex: EMAIL_GYM="gym@mail.com"
* PASSWORD_EMAIL_GYM -> E-mail address password
  * Ex: PASSWORD_EMAIL_GYM="mypassword"
## Execute using docker

1. Create a volume for postgres
> docker volume create volume_name

2. Run container
> docker run --name postgresql -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=my_password -p 5432:5432 -v "volume_name:/var/lib/pgsql/data" -d postgres

3. Open the instance of databse and execute the init_script

4. Inspect postgres container to get the gateway
> docker inspect postgresql | grep Gateway

5. Build image of golang
> docker build --tag docker-gym-consistency .

6. Run container 
> docker run --name back_gym_consistency -p 8080:8080 -e DBCONN="host=gateway_of_step_4 port=myport user=myuser password=mypassword dbname=postgres" -e SECRET_KEY_JWT="abcde12345" -e EMAIL_GYM="gym@mail.com" -e PASSWORD_EMAIL_GYM="mypassword" -d docker-gym-consistency

7. Check is it's working
> curl http://localhost:8080/getCurrentMonth

If the output was {"presentDays":0,"missedDays":0} the it's working. 

## Useful articles 
* [Monitoring golang with prometheus and grafana](https://gabrieltanner.org/blog/collecting-prometheus-metrics-in-golang/)
