
## API Reference

### 1. Create product

```http
  POST /products
```
 - **Request body (JSON):**
  ```json
  {
    "name": "Product Name",
    "description": "Product Description",
    "price": 19.99,
    "quantity": 10,
    "category": 1,
    "is_available": true
  }
  ```

### 2. Get all products

```http
  GET /products
```
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
        "category": 1,
        "is_available": true,
        "category_name": "Category Name",
	    "category_description": "Category Description"
      },
      {
        "id": 2,
        "name": "Another Product",
        "description": "Another Description",
        "price": 29.99,
        "quantity": 5,
        "category": 2,
        "is_available": false,
        "category_name": "Category Name",
	    "category_description": "Category Description"
      }
    ]
    ```


### 3. Get product

```http
  GET /products/{id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `int` | **Required**. Id of item to fetch |

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
      "category": 2,
      "is_available": false,
      "category_name": "Category Name",
	  "category_description": "Category Description"
    }
    ```
### 4. Update Product by ID

```http
  PUT /products/{id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `int` | **Required**. Id of item to fetch |
- **Request:**
  - **Content-Type:** `application/json`
  - **Body:** (example with all possible fields)
    ```json
    {
      "name": "Updated Product Name",
      "description": "Updated Description",
      "price": 25.99,
      "quantity": 15,
      "category": 2,
      "is_available": false
    }
    ```

### 5. Update Product Availability
```http
  PATCH /products/availability/{id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `int` | **Required**. Id of item to fetch |

- **Request:**
  - **Content-Type:** `application/json`
  - **Body:**
    ```json
    {
      "is_available": true
    }
    ```

### 6. Delete Product by ID

```http
  DELETE /products/{id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `int` | **Required**. Id of item to fetch |



## Category
### 1. Create category

```http
  POST /categories
```
 - **Request body (JSON):**
  ```json
  {
	"nameCategory": "Category name",
	"descriptionCategory": "Category description"
}
  ```

### 2. Get all categories

```http
  GET /categories
```
- **Response:**
  - **Status:** `200 OK`
  - **Content-Type:** `application/json`
  - **Body:**
    ```json
        [
        {
            "idCategory": 2,
            "nameCategory": "Clothing",
            "descriptionCategory": "Apparel and accessories"
        },
        {
            "idCategory": 3,
            "nameCategory": "Books",
            "descriptionCategory": "Printed and digital books"
        }
        ]
    ```


### 3. Get category

```http
  GET /categories/{id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `int` | **Required**. Id of item to fetch |

- **Response:**
  - **Status:** `200 OK`
  - **Content-Type:** `application/json`
  - **Body:**
    ```json
    {
        "idCategory": 3,
        "nameCategory": "Category name",
        "descriptionCategory": "Category description"
    }
    ```
### 4. Update Category by ID

```http
  PUT /category/{id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `int` | **Required**. Id of item to fetch |
- **Request:**
  - **Content-Type:** `application/json`
  - **Body:** (example with all possible fields)
    ```json
    {
        "nameCategory": "Category name",
        "descriptionCategory": "Category description"
    }
    ```

### 5. Delete Category

```http
  DELETE /categories/{id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `int` | **Required**. Id of item to fetch |
