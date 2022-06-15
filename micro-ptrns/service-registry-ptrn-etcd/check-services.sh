#!/bin/bash

printf "Test registration /services/user-service\n"
curl http://localhost:2379/v3/kv/range -X POST -d '{"key":"L3NlcnZpY2VzL3VzZXItc2VydmljZQo="}' #/services/user-service


printf "\nTest registration /services/order-service\n"
curl http://localhost:2379/v3/kv/range -X POST -d '{"key":"L3NlcnZpY2VzL29yZGVyLXNlcnZpY2UK"}' #/services/order-service


printf "\n\nWorking test /services/users-service\n"
curl -X POST http://localhost:8000/users -H "Content-Type: application/json" -d '{"id": "1", "name": "Alice"}'
curl -X GET http://localhost:8000/users/1

printf "\nWorking test /services/order-service\n"
curl -X POST http://localhost:8001/orders -H "Content-Type: application/json" -d '{"id": "101", "user_id": "1", "amount": 250}'
curl -X GET http://localhost:8001/orders/101


