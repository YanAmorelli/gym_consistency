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
> docker run --name back_gym_consistency -p 8080:8080 -e DBCONN="host=gateway_of_step_4 port=myport user=myuser password=mypassword dbname=postgres" -d docker-gym-consistency

7. Check is it's working
> curl http://localhost:8080/getCurrentMonth

If the output was {"presentDays":0,"missedDays":0} the it's working. 
