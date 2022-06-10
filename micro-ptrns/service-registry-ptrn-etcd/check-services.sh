#!/bin/bash

printf "Test registration /services/user-service\n"
curl http://localhost:2379/v3/kv/range -X POST -d '{"key":"L3NlcnZpY2VzL3VzZXItc2VydmljZQ=="}' #/services/user-service


printf "\nTest registration /services/order-service\n"
curl -L http://localhost:2379/v3/kv/range -X POST -d '{"key":"L3NlcnZpY2VzL29yZGVyLXNlcnZpY2U="}' #/services/order-service

printf "\nTest health /services/order-service\n"
curl -X GET http://localhost:8001/health
printf "\nTest communications /services/order-service\n"
curl -X GET http://localhost:8001/getuserserviceaddress

printf "\nTest health /services/user-service\n"
curl -X GET http://localhost:8000/health
printf "\nTest communications /services/user-service\n"
curl -X GET http://localhost:8000/getorderserviceaddress


printf "\n\nWorking test /services/users-service\n"
curl -X POST http://localhost:8000/users -H "Content-Type: application/json" -d '{"id": "1", "name": "Alice"}'
curl -X GET http://localhost:8000/users/1

printf "\nWorking test /services/order-service\n"
curl -X POST http://localhost:8001/orders -H "Content-Type: application/json" -d '{"id": "101", "user_id": "1", "amount": 250}'
curl -X GET http://localhost:8001/orders/101


