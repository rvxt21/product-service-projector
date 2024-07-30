# Product Management API

## Description

This API provides opportunities for managing products in the system. The API supports the creation, deletion, viewing, and full and partial updating of product information.

## Endpoints

### 1. Create product

- **Method:** `POST`
- **URL:** `/products`
- **Description:** Adds new product to the system.
- **Request body (JSON):**
  ```json
  {
    "name": "Product Name",
    "description": "Product Description",
    "price": 19.99,
    "quantity": 10,
    "category": "Category Name",
    "is_available": true
  }
  ```

### 2. Get all product

- **Method:** `GET`
- **URL:** `/products`
- **Description:** Returns all product from the system.
- **Response:**
  - **Status:** `200 OK`
  - **Content-Type:** `application/json`
  - **Body:**
    ```json
    [
      {
        "id": 1,
        "name": "Product Name",
        "description": "Product Description",
        "price": 19.99,
        "quantity": 10,
        "category": "Category Name",
        "is_available": true
      },
      {
        "id": 2,
        "name": "Another Product",
        "description": "Another Description",
        "price": 29.99,
        "quantity": 5,
        "category": "Another Category",
        "is_available": false
      }
    ]
    ```
- **Notes:**
  - If there are no products in the system, the response will be an empty array `[]`.

### 3. Get Product by ID

- **Method:** `GET`
- **URL:** `/products/{id}`
- **Description:** Retrieves a product by its ID.
- **Response:**
  - **Status:** `200 OK`
  - **Content-Type:** `application/json`
  - **Body:**
    ```json
    {
      "id": 1,
      "name": "Product Name",
      "description": "Product Description",
      "price": 19.99,
      "quantity": 10,
      "category": "Category Name",
      "is_available": true
    }
    ```

### 4. Update Product by ID

- **Method:** `UPDATE`
- **URL:** `/products/{id}`
- **Description:** Updates a product's details by its ID.
- **Request:**
  - **Content-Type:** `application/json`
  - **Body:** (example with all possible fields)
    ```json
    {
      "name": "Updated Product Name",
      "description": "Updated Description",
      "price": 25.99,
      "quantity": 15,
      "category": "Updated Category",
      "is_available": false
    }
    ```

### 5. Update Product Availability

- **Method:** `PATCH`
- **URL:** `/products/availability/{id}`
- **Description:** Updates the availability status of a product by its ID.
- **Request:**
  - **Content-Type:** `application/json`
  - **Body:**
    ```json
    {
      "is_available": true
    }
    ```

### 6. Delete Product by ID

- **Method:** `DELETE`
- **URL:** `/products/{id}`
- **Description:** Deletes a product by its ID.
- **Response:**
  - **Status:** `204 No Content`
