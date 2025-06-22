# this file contains curl-based commands that can be used for testing of API exposed by the server

# List
curl http://localhost:8080/customers
curl http://localhost:8080/organizations
curl http://localhost:8080/campaigns
curl http://localhost:8080/shorturls

# returns 404
curl http://localhost:8080/wrong

# Retrieve
# returns 400 - invalid input syntax for type uuid
curl http://localhost:8080/customers/xyz
curl http://localhost:8080/organizations/xyz
curl http://localhost:8080/campaigns/xyz

# returns 404
curl http://localhost:8080/customers/f81d4fae-7dec-11d0-a765-00a0c91e6bf6
curl http://localhost:8080/organizations/f81d4fae-7dec-11d0-a765-00a0c91e6bf6
curl http://localhost:8080/campaigns/f81d4fae-7dec-11d0-a765-00a0c91e6bf6

curl http://localhost:8080/shorturls/xyz



# Create
curl -X POST -H 'Content-Type: application/json' -d '{"email":"stranger@example.com","name":"Stranger"'} http://localhost:8080/customers

curl -X POST -H 'Content-Type: application/json' -d '{"name":"Tsofim"'} http://localhost:8080/organizations
# take ID of the organization and create customer that belongs to it. 
curl -X POST -H 'Content-Type: application/json' -d '{"email":"one@tsofim.com","name":"One", "organization_id": "735aef8a-4d24-11f0-9888-002b67d6b1c3"}' http://localhost:8080/customers
# take ID of organization and customer and create campaign
curl -X POST -H 'Content-Type: application/json' -d '{"name":"Tiul", "organization_id": "735aef8a-4d24-11f0-9888-002b67d6b1c3", "customer_id": "02695f62-4d25-11f0-9888-002b67d6b1c3"}' http://localhost:8080/campaigns



curl -X POST -H 'Content-Type: application/json' -d '{"target":"http://www.google.com", "campaign_id": "735aef8a-4d24-11f0-9888-002b67d6b1c3", "customer_id": "02695f62-4d25-11f0-9888-002b67d6b1c3"}' http://localhost:8080/shorturls

# custom short URL
curl -X POST -H 'Content-Type: application/json' -d '{"target":"http://www.google.com", "key": "cnn", "custom": true, "campaign_id": "735aef8a-4d24-11f0-9888-002b67d6b1c3", "customer_id": "02695f62-4d25-11f0-9888-002b67d6b1c3"}' http://localhost:8080/shorturls
