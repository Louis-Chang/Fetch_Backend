# Fetch_Backend
Fetch Backend Take-Home Exercise

## Prerequisites

- Go 1.21 or higher
- Docker (optional)

## Running the Application

### Using Go directly

1. Clone the repository:
```Bash
git clone https://github.com/Louis-Chang/Fetch_Backend.git
cd Fetch_Backend
```
2. Run the application
```Bash
go run main.go
```
The server will start on `http://localhost:8080`
### Using Docker
1. Build the Docker image
```Bash
docker build -t receipt-processor .
```
2. Run the container
```Bash
docker run -p 8080:8080 receipt-processor
```
## API Endpoints

### 1. Process Receipt
- **POST** `/receipts/process`
- Processes a receipt and returns a unique ID
- Request body example:
```json
{
    "retailer": "Target",
    "purchaseDate": "2022-01-02",
    "purchaseTime": "13:13",
    "total": "1.25",
    "items": [
        {
            "shortDescription": "Pepsi - 12-oz",
            "price": "1.25"
        }
    ]
}
```
- Response example:
```json
{
    "id": "c3b9e3d2-5d0f-424f-9fb6-a1996b49c86e"
}
```

### 2. Get Points
- **GET** `/receipts/{id}/points`
- Returns points awarded for the receipt
- Response example:
```json
{
    "points": 31
}
```

## Points Calculation Rules

Points are awarded based on the following rules:
1. One point for every alphanumeric character in the retailer name
2. 50 points if the total is a round dollar amount with no cents
3. 25 points if the total is a multiple of 0.25
4. 5 points for every two items on the receipt
5. If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up
6. 5 points if the total is greater than 10.00 (LLM rule)
7. 6 points if the day in the purchase date is odd
8. 10 points if the time of purchase is after 2:00pm and before 4:00pm

## Implementation Notes

- Data is stored in memory and does not persist after application restart
- All dates are in YYYY-MM-DD format
- All times are in 24-hour format (HH:MM)
- All monetary amounts have exactly two decimal places
- The API follows the OpenAPI 3.0.3 specification as defined in api.yml

## Project Structure
```
receipt-processor/
├── handlers/ # HTTP request handlers
├── models/ # Data models
├── services/ # Business logic
├── main.go # Application entry point
├── go.mod # Go module file
├── go.sum # Go module checksum
└── Dockerfile # Docker configuration
```