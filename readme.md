[![CircleCI](https://circleci.com/gh/jcorry/morellis/tree/master.svg?style=svg&circle-token=c66443d46cc348481a050ce58e1fb2c41a8803b9)](https://circleci.com/gh/jcorry/morellis/tree/master)
# Morellis API

The Morellis API is a simple marketing app for use by Morellis Ice Cream stores. It was born from a desire on my part
to be notified when my favorite flavors appear in the cooler so I never have to miss an opportunity to get a 
coconut jalapeno cone.

## What Does it Do?
The idea is that as the store personnel retires empty barrels and replaces them with
new flavors, they can update the current flavors using a mobile web app. Each of the store locations keeps 12 flavors current at a time. The company has ~40 total flavors. Each store's active/current flavors are selected from the company's ~40 total flavors. Company flavors change infrequently, current flavors at a store change frequently.

### Morellis Admin can:
- Manage store location data
- Manage company available flavors

### Morellis Staff Can:
- Manage store active flavors

### Customers can 
- Save their favorite flavors (or flavor keywords) using a very low friction SMS interface.
- Be notified when a flavor matching their saved keywords is made active at a store

## Client Growth Opportunity
I expect this to drive sales even as an MVP. Moving forward, the client might be approached
about the inclusion of additional features, such as:

- Allowing customers to tailor the frequency of their notifications
- Providing a web based interface for customer flavor keyword management
- Capturing customer visits in response to flavor update SMS messages by interfacing with the
Square POS API to cross reference customer phone numbers to SMS messages sent
- A reporting dashboard
- Broader inventory management for the store capturing data such as per flavor rate of consumption,
overall frequency of flavor change, average time on premises and real time inventory data

# How it works

## Start the API
- docker-compose up db
- go run cmd/api/\*.go~*_test.go

## Run the Tests
- docker-compose up db-test
- TEST_DSN="morellistest:testpass@tcp(127.0.0.1:33062)/morellistest?parseTime=true&multiStatements=true" go test -v -short ./...
(`-short` flag will skip integration tests)


