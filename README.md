# **Rate Birds API** ✨

## Overview
The Rate Birds API is a high-performance Go-based backend system designed for a bird rating application. It serves dynamic HTML content to a frontend, manages bird data, and implements an Elo-like ranking algorithm in a PostgreSQL database, with type-safe query generation via `sqlc`.

## Features
-   **Bird Data Ingestion**: Automatically fetches comprehensive bird information from the Nuthatch API on application startup, ensuring a rich dataset.
-   **Elo Rating System**: Implements a robust Elo rating algorithm to calculate and update bird scores based on head-to-head user ratings, reflecting their perceived "attractiveness" or "preference."
-   **Dynamic Leaderboard**: Provides real-time, sortable rankings of birds based on their Elo scores, delivered as HTML fragments for seamless frontend integration.
-   **Database Management**: Utilizes `sqlc` to generate type-safe Go code for interacting with the PostgreSQL database, enhancing data integrity and developer productivity.
-   **RESTful API & HTMX Support**: Exposes a clean API for bird matches and leaderboard data, designed to integrate effortlessly with HTMX-driven frontend components for a dynamic user experience without extensive JavaScript.

## Getting Started

### Installation
To set up and run the Rate Birds API locally, follow these steps:

1.  **Clone the Repository**:
    ```bash
    git clone https://github.com/neeeb1/rate_birds
    cd rate_birds
    ```

2.  **Install Go**:
    Ensure you have Go version 1.22 or higher installed. You can download it from [golang.org](https://golang.org/dl/).

3.  **Install PostgreSQL**:
    Set up a PostgreSQL database instance. For installation instructions, refer to the [official PostgreSQL documentation](https://www.postgresql.org/download/).

4.  **Install `goose` for Migrations**:
    ```bash
    go install github.com/pressly/goose/v3/cmd/goose@latest
    ```

5.  **Create Database and Apply Migrations**:
    Create a new database for the project (e.g., `rate_birds_db`) and apply the schema migrations:
    ```bash
    # Replace <your_db_url> with your PostgreSQL connection string
    # Example: postgres://user:password@host:port/database?sslmode=disable
    goose -dir ./sql/schema postgres "<your_db_url>" up
    ```

6.  **Install `sqlc`**:
    This project uses `sqlc` to generate Go code from SQL queries.
    ```bash
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
    ```
    Then, generate the database queries:
    ```bash
    sqlc generate
    ```

7.  **Fetch Dependencies**:
    ```bash
    go mod tidy
    ```

### Environment Variables
The following environment variables are required to run the application. Create a `.env` file in the project root with these variables:

-   `NUTHATCH_KEY`: Your API key for the Nuthatch bird data API.
    ```
    NUTHATCH_KEY=your_nuthatch_api_key_here
    ```
-   `DB_URL`: Your PostgreSQL database connection string.
    ```
    DB_URL=postgres://user:password@host:5432/rate_birds_db?sslmode=disable
    ```

## API Documentation

### Base URL
`http://localhost:8080/api`

### Endpoints

#### GET /api/scorematch/
Records the outcome of a bird match, updates the Elo ratings of the involved birds, and returns a new set of two random birds for the next rating interaction.

**Request Query Parameters**:
-   `winner`: `string` ("left" or "right") - Indicates which of the two displayed birds was chosen as the winner.
-   `leftBirdID`: `UUID` - The unique identifier of the bird presented on the left card.
-   `rightBirdID`: `UUID` - The unique identifier of the bird presented on the right card.

**Response**:
HTML snippet representing a new bird-wrapper div containing two dynamically generated bird cards ready for the next rating match.

```html
<div id="bird-wrapper" class="w-screen h-3/4 grid grid-flow-col justify-items-center">
    <div class="shadow-lg rounded-sm w-2/3 p-6 flex flex-col align-items-center bg-zinc-300" id="left-bird">
        <img class="card-image object-cover aspect-square object-contain" src="[Image URL of new left bird]">
        <div class="flex flex-col text-center">
            <p>[Common Name of new left bird]</p>
            <p><em>[Scientific Name of new left bird]</em></p>
            <Button class="[CSS Classes]" hx-get="/api/scorematch/" hx-trigger="click" hx-target="#bird-wrapper" hx-swap="outerHTML" hx-vals='{"winner": "left", "leftBirdID": "[UUID of new left bird]", "rightBirdID": "[UUID of new right bird]"}'>
                This one!
            </Button>
        </div>
    </div>
    <div class="card-separator inline-block self-center">OR</div>
    <div class="shadow-lg rounded-sm w-2/3 p-6 flex flex-col align-items-center bg-zinc-300" id="right-bird">
        <img  class="card-image object-cover aspect-square box-content" src="[Image URL of new right bird]">
        <div class="flex flex-col text-center">
            <p>[Common Name of new right bird]</p>
            <p><em>[Scientific Name of new right bird]</em></p>
            <Button class="[CSS Classes]" hx-get="/api/scorematch/" hx-trigger="click" hx-target="#bird-wrapper" hx-swap="outerHTML" hx-vals='{"winner": "right", "leftBirdID": "[UUID of new left bird]", "rightBirdID": "[UUID of new right bird]"}'>
                This one!
            </Button>
        </div>
    </div>
</div>
```

**Errors**:
-   `400 Bad Request`: Returned if `winner` is not "left" or "right", or if `leftBirdID` or `rightBirdID` are not valid UUIDs.
-   `500 Internal Server Error`: Occurs due to database connection issues, failure to retrieve bird data, or errors during Elo rating calculation.

#### GET /api/leaderboard/
Retrieves a list of top-ranked birds based on their Elo scores. The list length can be specified by a query parameter.

**Request Query Parameters**:
-   `listLength`: `integer` (e.g., `10`) - The maximum number of top-ranked birds to return in the leaderboard.

**Response**:
HTML snippet representing a table of birds, ordered by their Elo rating in descending order.

```html
<table>
    <tr>
        <th>Ranking</th>
        <th>Common Name</th>
        <th>ELO Score</th>
    </tr>
    <tr>
        <td>1.</td>
        <td>American Kestrel</td>
        <td>1550</td>
    </tr>
    <tr>
        <td>2.</td>
        <td>Northern Cardinal</td>
        <td>1520</td>
    </tr>
    <tr>
        <td>3.</td>
        <td>Blue Jay</td>
        <td>1490</td>
    </tr>
    <!-- ... more birds based on listLength ... -->
</table>
```

**Errors**:
-   `400 Bad Request`: Returned if `listLength` is not a valid integer.
-   `500 Internal Server Error`: Occurs due to database connection issues or errors retrieving ratings or bird details.

## Usage
Once the server is running, you can interact with the application:

1.  **Start the Server**:
    From the project root, execute:
    ```bash
    go run main.go
    ```
    The server will start on `http://localhost:8080`.

2.  **Access the Rating Interface**:
    Open your web browser and navigate to `http://localhost:8080`. You will be presented with two bird images. Click "This one!" under your preferred bird to rate it. The page will dynamically update with a new pair of birds.

3.  **View the Leaderboard**:
    Navigate to `http://localhost:8080/leaderboard.html` to see the current rankings of all rated birds. You can adjust the number of items displayed using the dropdown and "Reload" button.

4.  **Direct API Interaction**:
    You can also interact with the API endpoints directly using tools like `curl`:
    -   Get the top 5 birds on the leaderboard:
        ```bash
        curl "http://localhost:8080/api/leaderboard?listLength=5"
        ```

## Technologies Used
| Technology         | Description                                                               | Link                                                    |
| :----------------- | :------------------------------------------------------------------------ | :------------------------------------------------------ |
| **Go**             | Core language for the API backend, ensuring performance and concurrency.  | [golang.org](https://golang.org/)                       |
| **PostgreSQL**     | Robust relational database for storing bird and rating data.              | [postgresql.org](https://www.postgresql.org/)           |
| **`sqlc`**         | Generates type-safe Go code from SQL queries, improving development.      | [sqlc.dev](https://sqlc.dev/)                           |
| **HTMX**           | Enables dynamic, interactive HTML interfaces with minimal JavaScript.     | [htmx.org](https://htmx.org/)                           |
| **`go-dotenv`**    | Simple environment variable loading from `.env` files.                    | [github.com/joho/godotenv](https://github.com/joho/godotenv) |
| **Nuthatch API**   | External REST API providing comprehensive bird species data.              | [nuthatch.lastelm.software](https://nuthatch.lastelm.software/) |

## License
This project is open-source and licensed under the MIT License.


---

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![HTMX](https://img.shields.io/badge/HTMX-3069C7?style=for-the-badge&logo=data:image/svg+xml;base64,PHN2ZyB2aWV3Qm94PSIwIDAgMjQgMjQiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+PHBhdGggZmlsbD0iI2ZmZiIgZD0iTTEwLjI2OSAyLjAyNWwtLjAzLS4wNzUtLjAwNC0uMDA4LS41NzYtMS42NDktLjAwNC0uMDA4LS4wMy0uMDc0LS4xNS0uMzcySDB2MjRsMTkuMDIzLjA1NS4wMDMtLjAwNS4wNzgtLjIwNy4wMjctLjA3NC4wMDQtLjAwOEw4Ljg4NCAyNC4wMDlsLS4wMDUtLjAwOHMtMy42LS4zMzUtNC43NzItLjQ0N2ExLjE5LjExOSAwIDAgMS0uNjEyLS43MDZjLS4xNjQtLjU2NC4xNDQtMS4yMDggLjcwOC0xLjM3Mi41NjQtLjE2NC43MDgtLjI1Mi44Mi0uNjc2LjE0NC0uNjE2LTIuMDA0LS44NzItMi43NTItLjk2LTEuNTU2LS4xODgtMS43MTItLjI5Ny0xLjk2OC0uMzkyLS44Mi4wMzYtMS42NC4wMzYtMi41Ni4zNzYtLjUwOC4yMDQtMS4wMy40ODctMS41NzYuODk5LS4xMy4xMDgtLjI1Mi4yMjgtLjM4OC4zNjQtLjIwOC4xODgtLjQyNC4zNzItLjY3Ni41NDQtLjYxMi40NDQtMS40OTYuNjI4LTIuMjQ0LjYwNC0uOTIyLS4wMzYtMS43ODQtLjI5Ny0yLjM5Mi0uNzU3LS41NDgtLjQzMi0uNzg0LTEuMDgtLjYyLTEuNTA5LjE2NC0uNDE2LjQ2OC0uNjI4LjcyOC0uNjg4LjI0LS4wNjQuNjUtLjA2NC45MDguMDcyLjE5Mi4xMDIuNTUyLjIyOC44NTYuNDA0LjI5Mi4xOC42MDQuMzY4LjkyOC41NC42NDQuMzQ4IDEuNjY4LjUzMiAyLjg1Mi4zMTYuNzk2LS4xNDggMS42ODQtLjQ3NiAyLjM4LS42OC41NzItLjE2LjkyLjIzNi45OC42NjQuMDY4LjQzMi0uMDcyLjcyLS40NTYuODQtLjM4NC4xMjgtLjg3Mi4wNTItMS4xMDgtLjA2OC0uOTg4LS41MDQtMi4xMDQtLjgxNi0zLjIxMi0uOTA4LS4wNzYtLjAwNC0uMDc2LS4wMDQtLjA5Mi0uMDA0aC02LjMzMnYtNS45Nmg1LjE0NmMzLjQ0LS4zNTIgNi4zMjgtLjQ3NiA4LjcwNC0uNzc2IDYuMjQtLjYxNiA4Ljg0OC0xLjgyIDguODQ4LTUuNjYgMC0zLjgxNi0yLjg4LTUuMjcyLTguNTQ4LTUuNzY0LTMuMjktLjI4NC03LjQ1Mi0uMzY4LTEwLjk1Ni0uMzg0LS4wMTYtLjAwNC0uMjctLjAwNC0uNjEyLS4wMDVsLS4wNzYtLjAwNGwtLjk3Ni0uMDQ4YTEuMTkyLjQ3NiAwIDAgMS0uNzg4LS40NzYgMS4wNiAxLjA2IDAgMCAxIC4wNDgtLjgyMWMuMDk2LS4zMTIuMjk2LS40NzYuNTg4LS41NDRsLjM5Ni0uMDk2LjcwNC0uMDQ0Yy44Mi0uMDk2IDEuOTM2LS4xNTEgMy4yMTYtLjE1MXoiLz48L3N2Zz4=)](https://htmx.org/)
[![Readme was generated by Dokugen](https://img.shields.io/badge/Readme%20was%20generated%20by-Dokugen-brightgreen)](https://www.npmjs.com/package/dokugen)
