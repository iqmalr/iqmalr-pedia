

# REST API Development Tasks for E-commerce Multi-vendor Platform

## 1. USER & AUTHENTICATION (Prioritas: Tertinggi - MVP)

### 1.1 Authentication Endpoints

**Task 1.1.1: User Registration**
- Endpoint: `POST /auth/register`
- HTTP Method: POST
- Request Body:
  ```json
  {
    "name": "string",
    "email": "string",
    "phone": "string",
    "password": "string",
    "role": "string" // optional, default: 'customer'
  }
  ```
- Response Format (201):
  ```json
  {
    "message": "User registered successfully",
    "user": {
      "id": "integer",
      "uuid": "string",
      "name": "string",
      "email": "string",
      "phone": "string",
      "role": "string",
      "is_active": "boolean",
      "created_at": "timestamp"
    }
  }
  ```
- Validasi:
    - name: required, min 3 characters
    - email: required, valid email format, unique
    - phone: optional, valid phone format if provided
    - password: required, min 8 characters
    - role: optional, must be one of 'customer', 'vendor', 'admin'
- Authorization: None
- Prioritas: High (MVP)

**Task 1.1.2: User Login**
- Endpoint: `POST /auth/login`
- HTTP Method: POST
- Request Body:
  ```json
  {
    "email": "string",
    "password": "string"
  }
  ```
- Response Format (200):
  ```json
  {
    "message": "Login successful",
    "token": "string",
    "user": {
      "id": "integer",
      "uuid": "string",
      "name": "string",
      "email": "string",
      "role": "string",
      "is_active": "boolean"
    }
  }
  ```
- Validasi:
    - email: required, must exist in database
    - password: required, must match
- Authorization: None
- Prioritas: High (MVP)

**Task 1.1.3: User Logout**
- Endpoint: `POST /auth/logout`
- HTTP Method: POST
- Request Body: None
- Response Format (200):
  ```json
  {
    "message": "Logout successful"
  }
  ```
- Validasi: None
- Authorization: Bearer token (any authenticated user)
- Prioritas: High (MVP)

**Task 1.1.4: Refresh Token**
- Endpoint: `POST /auth/refresh`
- HTTP Method: POST
- Request Body: None
- Response Format (200):
  ```json
  {
    "token": "string"
  }
  ```
- Validasi: None
- Authorization: Refresh token
- Prioritas: High (MVP)

**Task 1.1.5: Forgot Password**
- Endpoint: `POST /auth/forgot-password`
- HTTP Method: POST
- Request Body:
  ```json
  {
    "email": "string"
  }
  ```
- Response Format (200):
  ```json
  {
    "message": "Password reset link sent to your email"
  }
  ```
- Validasi:
    - email: required, must exist in database
- Authorization: None
- Prioritas: Medium

**Task 1.1.6: Reset Password**
- Endpoint: `POST /auth/reset-password`
- HTTP Method: POST
- Request Body:
  ```json
  {
    "token": "string",
    "password": "string",
    "password_confirmation": "string"
  }
  ```
- Response Format (200):
  ```json
  {
    "message": "Password reset successful"
  }
  ```
- Validasi:
    - token: required, valid
    - password: required, min 8 characters
    - password_confirmation: required, must match password
- Authorization: None
- Prioritas: Medium

**Task 1.1.7: Email Verification**
- Endpoint: `POST /auth/verify-email`
- HTTP Method: POST
- Request Body:
  ```json
  {
    "token": "string"
  }
  ```
- Response Format (200):
  ```json
  {
    "message": "Email verified successfully"
  }
  ```
- Validasi:
    - token: required, valid
- Authorization: None
- Prioritas: Medium

**Task 1.1.8: Resend Email Verification**
- Endpoint: `POST /auth/resend-verification`
- HTTP Method: POST
- Request Body:
  ```json
  {
    "email": "string"
  }
  ```
- Response Format (200):
  ```json
  {
    "message": "Verification email sent"
  }
  ```
- Validasi:
    - email: required, must exist in database
- Authorization: None
- Prioritas: Medium

### 1.2 User Management Endpoints

**Task 1.2.1: Get Current User Profile**
- Endpoint: `GET /users/me`
- HTTP Method: GET
- Request Body: None
- Response Format (200):
  ```json
  {
    "id": "integer",
    "uuid": "string",
    "name": "string",
    "email": "string",
    "phone": "string",
    "role": "string",
    "avatar_url": "string",
    "is_active": "boolean",
    "email_verified_at": "timestamp",
    "phone_verified_at": "timestamp",
    "last_login_at": "timestamp",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi: None
- Authorization: Bearer token (any authenticated user)
- Prioritas: High (MVP)

**Task 1.2.2: Update Current User Profile**
- Endpoint: `PUT /users/me`
- HTTP Method: PUT
- Request Body:
  ```json
  {
    "name": "string",
    "phone": "string",
    "avatar_url": "string"
  }
  ```
- Response Format (200):
  ```json
  {
    "id": "integer",
    "uuid": "string",
    "name": "string",
    "email": "string",
    "phone": "string",
    "role": "string",
    "avatar_url": "string",
    "is_active": "boolean",
    "email_verified_at": "timestamp",
    "phone_verified_at": "timestamp",
    "last_login_at": "timestamp",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - name: optional, min 3 characters
    - phone: optional, valid phone format if provided
    - avatar_url: optional, valid URL if provided
- Authorization: Bearer token (any authenticated user)
- Prioritas: High (MVP)

**Task 1.2.3: Change Password**
- Endpoint: `PUT /users/me/password`
- HTTP Method: PUT
- Request Body:
  ```json
  {
    "current_password": "string",
    "new_password": "string",
    "new_password_confirmation": "string"
  }
  ```
- Response Format (200):
  ```json
  {
    "message": "Password changed successfully"
  }
  ```
- Validasi:
    - current_password: required, must match current password
    - new_password: required, min 8 characters
    - new_password_confirmation: required, must match new_password
- Authorization: Bearer token (any authenticated user)
- Prioritas: High (MVP)

**Task 1.2.4: Get User by ID (Admin only)**
- Endpoint: `GET /users/{id}`
- HTTP Method: GET
- Path Params: id (integer)
- Request Body: None
- Response Format (200):
  ```json
  {
    "id": "integer",
    "uuid": "string",
    "name": "string",
    "email": "string",
    "phone": "string",
    "role": "string",
    "avatar_url": "string",
    "is_active": "boolean",
    "email_verified_at": "timestamp",
    "phone_verified_at": "timestamp",
    "last_login_at": "timestamp",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - id: required, must exist in database
- Authorization: Bearer token (admin only)
- Prioritas: Medium

**Task 1.2.5: Update User (Admin only)**
- Endpoint: `PUT /users/{id}`
- HTTP Method: PUT
- Path Params: id (integer)
- Request Body:
  ```json
  {
    "name": "string",
    "email": "string",
    "phone": "string",
    "role": "string",
    "is_active": "boolean"
  }
  ```
- Response Format (200):
  ```json
  {
    "id": "integer",
    "uuid": "string",
    "name": "string",
    "email": "string",
    "phone": "string",
    "role": "string",
    "avatar_url": "string",
    "is_active": "boolean",
    "email_verified_at": "timestamp",
    "phone_verified_at": "timestamp",
    "last_login_at": "timestamp",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - id: required, must exist in database
    - name: optional, min 3 characters
    - email: optional, valid email format, unique
    - phone: optional, valid phone format if provided
    - role: optional, must be one of 'customer', 'vendor', 'admin'
    - is_active: optional, boolean
- Authorization: Bearer token (admin only)
- Prioritas: Medium

**Task 1.2.6: Deactivate User (Admin only)**
- Endpoint: `DELETE /users/{id}`
- HTTP Method: DELETE
- Path Params: id (integer)
- Request Body: None
- Response Format (200):
  ```json
  {
    "message": "User deactivated successfully"
  }
  ```
- Validasi:
    - id: required, must exist in database
- Authorization: Bearer token (admin only)
- Prioritas: Medium

**Task 1.2.7: List Users (Admin only)**
- Endpoint: `GET /users`
- HTTP Method: GET
- Query Params:
    - page (integer, optional)
    - limit (integer, optional)
    - search (string, optional)
    - role (string, optional)
    - is_active (boolean, optional)
- Request Body: None
- Response Format (200):
  ```json
  {
    "data": [
      {
        "id": "integer",
        "uuid": "string",
        "name": "string",
        "email": "string",
        "phone": "string",
        "role": "string",
        "is_active": "boolean",
        "created_at": "timestamp"
      }
    ],
    "pagination": {
      "total": "integer",
      "page": "integer",
      "limit": "integer",
      "total_pages": "integer"
    }
  }
  ```
- Validasi:
    - page: optional, integer, min 1
    - limit: optional, integer, min 1, max 100
    - search: optional, string
    - role: optional, must be one of 'customer', 'vendor', 'admin'
    - is_active: optional, boolean
- Authorization: Bearer token (admin only)
- Prioritas: Medium

### 1.3 User Address Management Endpoints

**Task 1.3.1: Get User Addresses**
- Endpoint: `GET /users/me/addresses`
- HTTP Method: GET
- Request Body: None
- Response Format (200):
  ```json
  {
    "data": [
      {
        "id": "integer",
        "label": "string",
        "recipient_name": "string",
        "phone": "string",
        "address_line1": "string",
        "address_line2": "string",
        "city": "string",
        "state": "string",
        "postal_code": "string",
        "country": "string",
        "is_default": "boolean",
        "latitude": "decimal",
        "longitude": "decimal",
        "created_at": "timestamp",
        "updated_at": "timestamp"
      }
    ]
  }
  ```
- Validasi: None
- Authorization: Bearer token (any authenticated user)
- Prioritas: High (MVP)

**Task 1.3.2: Create User Address**
- Endpoint: `POST /users/me/addresses`
- HTTP Method: POST
- Request Body:
  ```json
  {
    "label": "string",
    "recipient_name": "string",
    "phone": "string",
    "address_line1": "string",
    "address_line2": "string",
    "city": "string",
    "state": "string",
    "postal_code": "string",
    "country": "string",
    "is_default": "boolean",
    "latitude": "decimal",
    "longitude": "decimal"
  }
  ```
- Response Format (201):
  ```json
  {
    "id": "integer",
    "label": "string",
    "recipient_name": "string",
    "phone": "string",
    "address_line1": "string",
    "address_line2": "string",
    "city": "string",
    "state": "string",
    "postal_code": "string",
    "country": "string",
    "is_default": "boolean",
    "latitude": "decimal",
    "longitude": "decimal",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - recipient_name: required, min 3 characters
    - phone: required, valid phone format
    - address_line1: required, min 5 characters
    - city: required, min 2 characters
    - state: required, min 2 characters
    - postal_code: required, min 3 characters
    - country: optional, default: 'Indonesia'
    - label: optional, max 50 characters
    - address_line2: optional
    - is_default: optional, boolean
    - latitude: optional, decimal
    - longitude: optional, decimal
- Authorization: Bearer token (any authenticated user)
- Prioritas: High (MVP)

**Task 1.3.3: Get User Address by ID**
- Endpoint: `GET /users/me/addresses/{id}`
- HTTP Method: GET
- Path Params: id (integer)
- Request Body: None
- Response Format (200):
  ```json
  {
    "id": "integer",
    "label": "string",
    "recipient_name": "string",
    "phone": "string",
    "address_line1": "string",
    "address_line2": "string",
    "city": "string",
    "state": "string",
    "postal_code": "string",
    "country": "string",
    "is_default": "boolean",
    "latitude": "decimal",
    "longitude": "decimal",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - id: required, must exist and belong to current user
- Authorization: Bearer token (any authenticated user)
- Prioritas: High (MVP)

**Task 1.3.4: Update User Address**
- Endpoint: `PUT /users/me/addresses/{id}`
- HTTP Method: PUT
- Path Params: id (integer)
- Request Body:
  ```json
  {
    "label": "string",
    "recipient_name": "string",
    "phone": "string",
    "address_line1": "string",
    "address_line2": "string",
    "city": "string",
    "state": "string",
    "postal_code": "string",
    "country": "string",
    "is_default": "boolean",
    "latitude": "decimal",
    "longitude": "decimal"
  }
  ```
- Response Format (200):
  ```json
  {
    "id": "integer",
    "label": "string",
    "recipient_name": "string",
    "phone": "string",
    "address_line1": "string",
    "address_line2": "string",
    "city": "string",
    "state": "string",
    "postal_code": "string",
    "country": "string",
    "is_default": "boolean",
    "latitude": "decimal",
    "longitude": "decimal",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - id: required, must exist and belong to current user
    - recipient_name: optional, min 3 characters
    - phone: optional, valid phone format if provided
    - address_line1: optional, min 5 characters
    - city: optional, min 2 characters
    - state: optional, min 2 characters
    - postal_code: optional, min 3 characters
    - country: optional
    - label: optional, max 50 characters
    - address_line2: optional
    - is_default: optional, boolean
    - latitude: optional, decimal
    - longitude: optional, decimal
- Authorization: Bearer token (any authenticated user)
- Prioritas: High (MVP)

**Task 1.3.5: Delete User Address**
- Endpoint: `DELETE /users/me/addresses/{id}`
- HTTP Method: DELETE
- Path Params: id (integer)
- Request Body: None
- Response Format (200):
  ```json
  {
    "message": "Address deleted successfully"
  }
  ```
- Validasi:
    - id: required, must exist and belong to current user
- Authorization: Bearer token (any authenticated user)
- Prioritas: High (MVP)

**Task 1.3.6: Set Default Address**
- Endpoint: `PUT /users/me/addresses/{id}/default`
- HTTP Method: PUT
- Path Params: id (integer)
- Request Body: None
- Response Format (200):
  ```json
  {
    "message": "Default address updated successfully"
  }
  ```
- Validasi:
    - id: required, must exist and belong to current user
- Authorization: Bearer token (any authenticated user)
- Prioritas: High (MVP)

## 2. VENDOR MANAGEMENT (Prioritas: Tinggi - MVP)

### 2.1 Vendor Endpoints

**Task 2.1.1: Create Vendor**
- Endpoint: `POST /vendors`
- HTTP Method: POST
- Request Body:
  ```json
  {
    "name": "string",
    "description": "string",
    "contact_email": "string",
    "contact_phone": "string",
    "business_type": "string",
    "tax_id": "string",
    "address_line1": "string",
    "address_line2": "string",
    "city": "string",
    "state": "string",
    "postal_code": "string",
    "country": "string"
  }
  ```
- Response Format (201):
  ```json
  {
    "id": "integer",
    "uuid": "string",
    "owner_id": "integer",
    "name": "string",
    "slug": "string",
    "description": "string",
    "contact_email": "string",
    "contact_phone": "string",
    "status": "string",
    "business_type": "string",
    "tax_id": "string",
    "address_line1": "string",
    "address_line2": "string",
    "city": "string",
    "state": "string",
    "postal_code": "string",
    "country": "string",
    "total_products": "integer",
    "total_sales": "decimal",
    "rating_avg": "decimal",
    "total_reviews": "integer",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - name: required, min 3 characters
    - contact_email: required, valid email format
    - contact_phone: optional, valid phone format if provided
    - business_type: optional, must be one of 'individual', 'company'
    - tax_id: optional, max 50 characters
    - address_line1: optional, min 5 characters
    - city: optional, min 2 characters
    - state: optional, min 2 characters
    - postal_code: optional, min 3 characters
    - country: optional, default: 'Indonesia'
- Authorization: Bearer token (users with role 'vendor' or 'admin')
- Prioritas: High (MVP)

**Task 2.1.2: Get Vendor by ID**
- Endpoint: `GET /vendors/{id}`
- HTTP Method: GET
- Path Params: id (integer)
- Request Body: None
- Response Format (200):
  ```json
  {
    "id": "integer",
    "uuid": "string",
    "owner_id": "integer",
    "name": "string",
    "slug": "string",
    "description": "string",
    "logo_url": "string",
    "banner_url": "string",
    "contact_email": "string",
    "contact_phone": "string",
    "status": "string",
    "business_type": "string",
    "tax_id": "string",
    "address_line1": "string",
    "address_line2": "string",
    "city": "string",
    "state": "string",
    "postal_code": "string",
    "country": "string",
    "total_products": "integer",
    "total_sales": "decimal",
    "rating_avg": "decimal",
    "total_reviews": "integer",
    "verified_at": "timestamp",
    "approved_at": "timestamp",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - id: required, must exist in database
- Authorization: Bearer token (any authenticated user, vendor owner, or admin)
- Prioritas: High (MVP)

**Task 2.1.3: Get Vendor by Slug**
- Endpoint: `GET /vendors/slug/{slug}`
- HTTP Method: GET
- Path Params: slug (string)
- Request Body: None
- Response Format (200):
  ```json
  {
    "id": "integer",
    "uuid": "string",
    "owner_id": "integer",
    "name": "string",
    "slug": "string",
    "description": "string",
    "logo_url": "string",
    "banner_url": "string",
    "contact_email": "string",
    "contact_phone": "string",
    "status": "string",
    "business_type": "string",
    "tax_id": "string",
    "address_line1": "string",
    "address_line2": "string",
    "city": "string",
    "state": "string",
    "postal_code": "string",
    "country": "string",
    "total_products": "integer",
    "total_sales": "decimal",
    "rating_avg": "decimal",
    "total_reviews": "integer",
    "verified_at": "timestamp",
    "approved_at": "timestamp",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - slug: required, must exist in database
- Authorization: Bearer token (any authenticated user, vendor owner, or admin)
- Prioritas: High (MVP)

**Task 2.1.4: Update Vendor**
- Endpoint: `PUT /vendors/{id}`
- HTTP Method: PUT
- Path Params: id (integer)
- Request Body:
  ```json
  {
    "name": "string",
    "description": "string",
    "contact_email": "string",
    "contact_phone": "string",
    "business_type": "string",
    "tax_id": "string",
    "address_line1": "string",
    "address_line2": "string",
    "city": "string",
    "state": "string",
    "postal_code": "string",
    "country": "string"
  }
  ```
- Response Format (200):
  ```json
  {
    "id": "integer",
    "uuid": "string",
    "owner_id": "integer",
    "name": "string",
    "slug": "string",
    "description": "string",
    "logo_url": "string",
    "banner_url": "string",
    "contact_email": "string",
    "contact_phone": "string",
    "status": "string",
    "business_type": "string",
    "tax_id": "string",
    "address_line1": "string",
    "address_line2": "string",
    "city": "string",
    "state": "string",
    "postal_code": "string",
    "country": "string",
    "total_products": "integer",
    "total_sales": "decimal",
    "rating_avg": "decimal",
    "total_reviews": "integer",
    "verified_at": "timestamp",
    "approved_at": "timestamp",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - id: required, must exist and belong to current user or admin
    - name: optional, min 3 characters
    - contact_email: optional, valid email format if provided
    - contact_phone: optional, valid phone format if provided
    - business_type: optional, must be one of 'individual', 'company'
    - tax_id: optional, max 50 characters
    - address_line1: optional, min 5 characters
    - city: optional, min 2 characters
    - state: optional, min 2 characters
    - postal_code: optional, min 3 characters
    - country: optional
- Authorization: Bearer token (vendor owner or admin)
- Prioritas: High (MVP)

**Task 2.1.5: Upload Vendor Logo**
- Endpoint: `POST /vendors/{id}/logo`
- HTTP Method: POST
- Path Params: id (integer)
- Request Body: multipart/form-data with file
- Response Format (200):
  ```json
  {
    "message": "Logo uploaded successfully",
    "logo_url": "string"
  }
  ```
- Validasi:
    - id: required, must exist and belong to current user or admin
    - file: required, image file (jpg, png, etc), max size 2MB
- Authorization: Bearer token (vendor owner or admin)
- Prioritas: Medium

**Task 2.1.6: Upload Vendor Banner**
- Endpoint: `POST /vendors/{id}/banner`
- HTTP Method: POST
- Path Params: id (integer)
- Request Body: multipart/form-data with file
- Response Format (200):
  ```json
  {
    "message": "Banner uploaded successfully",
    "banner_url": "string"
  }
  ```
- Validasi:
    - id: required, must exist and belong to current user or admin
    - file: required, image file (jpg, png, etc), max size 2MB
- Authorization: Bearer token (vendor owner or admin)
- Prioritas: Medium

**Task 2.1.7: Update Vendor Status (Admin only)**
- Endpoint: `PUT /vendors/{id}/status`
- HTTP Method: PUT
- Path Params: id (integer)
- Request Body:
  ```json
  {
    "status": "string"
  }
  ```
- Response Format (200):
  ```json
  {
    "message": "Vendor status updated successfully",
    "status": "string"
  }
  ```
- Validasi:
    - id: required, must exist in database
    - status: required, must be one of 'pending', 'active', 'suspended', 'rejected'
- Authorization: Bearer token (admin only)
- Prioritas: High (MVP)

**Task 2.1.8: List Vendors**
- Endpoint: `GET /vendors`
- HTTP Method: GET
- Query Params:
    - page (integer, optional)
    - limit (integer, optional)
    - search (string, optional)
    - status (string, optional)
    - sort (string, optional, default: 'created_at')
    - order (string, optional, default: 'desc')
- Request Body: None
- Response Format (200):
  ```json
  {
    "data": [
      {
        "id": "integer",
        "uuid": "string",
        "name": "string",
        "slug": "string",
        "logo_url": "string",
        "contact_email": "string",
        "status": "string",
        "total_products": "integer",
        "rating_avg": "decimal",
        "total_reviews": "integer",
        "created_at": "timestamp"
      }
    ],
    "pagination": {
      "total": "integer",
      "page": "integer",
      "limit": "integer",
      "total_pages": "integer"
    }
  }
  ```
- Validasi:
    - page: optional, integer, min 1
    - limit: optional, integer, min 1, max 100
    - search: optional, string
    - status: optional, must be one of 'pending', 'active', 'suspended', 'rejected'
    - sort: optional, string, must be a valid column name
    - order: optional, string, must be 'asc' or 'desc'
- Authorization: Bearer token (any authenticated user)
- Prioritas: High (MVP)

### 2.2 Vendor User Management Endpoints

**Task 2.2.1: Add User to Vendor**
- Endpoint: `POST /vendors/{id}/users`
- HTTP Method: POST
- Path Params: id (integer)
- Request Body:
  ```json
  {
    "user_id": "integer",
    "role": "string",
    "permissions": "object"
  }
  ```
- Response Format (201):
  ```json
  {
    "id": "integer",
    "vendor_id": "integer",
    "user_id": "integer",
    "role": "string",
    "permissions": "object",
    "is_active": "boolean",
    "invited_at": "timestamp",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - id: required, must exist and belong to current user or admin
    - user_id: required, must exist in database
    - role: required, must be one of 'owner', 'admin', 'staff'
    - permissions: optional, JSON object
- Authorization: Bearer token (vendor owner or admin)
- Prioritas: Medium

**Task 2.2.2: Get Vendor Users**
- Endpoint: `GET /vendors/{id}/users`
- HTTP Method: GET
- Path Params: id (integer)
- Request Body: None
- Response Format (200):
  ```json
  {
    "data": [
      {
        "id": "integer",
        "vendor_id": "integer",
        "user_id": "integer",
        "user": {
          "id": "integer",
          "name": "string",
          "email": "string"
        },
        "role": "string",
        "permissions": "object",
        "is_active": "boolean",
        "invited_at": "timestamp",
        "joined_at": "timestamp",
        "created_at": "timestamp"
      }
    ]
  }
  ```
- Validasi:
    - id: required, must exist and belong to current user or admin
- Authorization: Bearer token (vendor owner, admin, or vendor users)
- Prioritas: Medium

**Task 2.2.3: Update Vendor User Role**
- Endpoint: `PUT /vendors/{vendorId}/users/{userId}`
- HTTP Method: PUT
- Path Params:
    - vendorId (integer)
    - userId (integer)
- Request Body:
  ```json
  {
    "role": "string",
    "permissions": "object",
    "is_active": "boolean"
  }
  ```
- Response Format (200):
  ```json
  {
    "id": "integer",
    "vendor_id": "integer",
    "user_id": "integer",
    "role": "string",
    "permissions": "object",
    "is_active": "boolean",
    "invited_at": "timestamp",
    "joined_at": "timestamp",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - vendorId: required, must exist and belong to current user or admin
    - userId: required, must exist and be associated with the vendor
    - role: optional, must be one of 'owner', 'admin', 'staff'
    - permissions: optional, JSON object
    - is_active: optional, boolean
- Authorization: Bearer token (vendor owner or admin)
- Prioritas: Medium

**Task 2.2.4: Remove User from Vendor**
- Endpoint: `DELETE /vendors/{vendorId}/users/{userId}`
- HTTP Method: DELETE
- Path Params:
    - vendorId (integer)
    - userId (integer)
- Request Body: None
- Response Format (200):
  ```json
  {
    "message": "User removed from vendor successfully"
  }
  ```
- Validasi:
    - vendorId: required, must exist and belong to current user or admin
    - userId: required, must exist and be associated with the vendor
- Authorization: Bearer token (vendor owner or admin)
- Prioritas: Medium

### 2.3 Vendor Settings Endpoints

**Task 2.3.1: Get Vendor Settings**
- Endpoint: `GET /vendors/{id}/settings`
- HTTP Method: GET
- Path Params: id (integer)
- Request Body: None
- Response Format (200):
  ```json
  {
    "data": [
      {
        "id": "integer",
        "vendor_id": "integer",
        "setting_key": "string",
        "setting_value": "object",
        "created_at": "timestamp",
        "updated_at": "timestamp"
      }
    ]
  }
  ```
- Validasi:
    - id: required, must exist and belong to current user or admin
- Authorization: Bearer token (vendor owner, admin, or vendor users)
- Prioritas: Medium

**Task 2.3.2: Update Vendor Setting**
- Endpoint: `PUT /vendors/{id}/settings/{key}`
- HTTP Method: PUT
- Path Params:
    - id (integer)
    - key (string)
- Request Body:
  ```json
  {
    "setting_value": "object"
  }
  ```
- Response Format (200):
  ```json
  {
    "id": "integer",
    "vendor_id": "integer",
    "setting_key": "string",
    "setting_value": "object",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - id: required, must exist and belong to current user or admin
    - key: required, string
    - setting_value: required, JSON object
- Authorization: Bearer token (vendor owner or admin)
- Prioritas: Medium

### 2.4 Vendor Bank Accounts Endpoints

**Task 2.4.1: Add Vendor Bank Account**
- Endpoint: `POST /vendors/{id}/bank-accounts`
- HTTP Method: POST
- Path Params: id (integer)
- Request Body:
  ```json
  {
    "bank_name": "string",
    "account_number": "string",
    "account_holder": "string",
    "is_primary": "boolean"
  }
  ```
- Response Format (201):
  ```json
  {
    "id": "integer",
    "vendor_id": "integer",
    "bank_name": "string",
    "account_number": "string",
    "account_holder": "string",
    "is_verified": "boolean",
    "is_primary": "boolean",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - id: required, must exist and belong to current user or admin
    - bank_name: required, min 2 characters
    - account_number: required, min 5 characters
    - account_holder: required, min 3 characters
    - is_primary: optional, boolean
- Authorization: Bearer token (vendor owner or admin)
- Prioritas: Medium

**Task 2.4.2: Get Vendor Bank Accounts**
- Endpoint: `GET /vendors/{id}/bank-accounts`
- HTTP Method: GET
- Path Params: id (integer)
- Request Body: None
- Response Format (200):
  ```json
  {
    "data": [
      {
        "id": "integer",
        "vendor_id": "integer",
        "bank_name": "string",
        "account_number": "string",
        "account_holder": "string",
        "is_verified": "boolean",
        "is_primary": "boolean",
        "created_at": "timestamp",
        "updated_at": "timestamp"
      }
    ]
  }
  ```
- Validasi:
    - id: required, must exist and belong to current user or admin
- Authorization: Bearer token (vendor owner, admin, or vendor users)
- Prioritas: Medium

**Task 2.4.3: Update Vendor Bank Account**
- Endpoint: `PUT /vendors/{vendorId}/bank-accounts/{accountId}`
- HTTP Method: PUT
- Path Params:
    - vendorId (integer)
    - accountId (integer)
- Request Body:
  ```json
  {
    "bank_name": "string",
    "account_number": "string",
    "account_holder": "string",
    "is_primary": "boolean"
  }
  ```
- Response Format (200):
  ```json
  {
    "id": "integer",
    "vendor_id": "integer",
    "bank_name": "string",
    "account_number": "string",
    "account_holder": "string",
    "is_verified": "boolean",
    "is_primary": "boolean",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - vendorId: required, must exist and belong to current user or admin
    - accountId: required, must exist and belong to the vendor
    - bank_name: optional, min 2 characters
    - account_number: optional, min 5 characters
    - account_holder: optional, min 3 characters
    - is_primary: optional, boolean
- Authorization: Bearer token (vendor owner or admin)
- Prioritas: Medium

**Task 2.4.4: Delete Vendor Bank Account**
- Endpoint: `DELETE /vendors/{vendorId}/bank-accounts/{accountId}`
- HTTP Method: DELETE
- Path Params:
    - vendorId (integer)
    - accountId (integer)
- Request Body: None
- Response Format (200):
  ```json
  {
    "message": "Bank account deleted successfully"
  }
  ```
- Validasi:
    - vendorId: required, must exist and belong to current user or admin
    - accountId: required, must exist and belong to the vendor
- Authorization: Bearer token (vendor owner or admin)
- Prioritas: Medium

**Task 2.4.5: Set Primary Bank Account**
- Endpoint: `PUT /vendors/{vendorId}/bank-accounts/{accountId}/primary`
- HTTP Method: PUT
- Path Params:
    - vendorId (integer)
    - accountId (integer)
- Request Body: None
- Response Format (200):
  ```json
  {
    "message": "Primary bank account updated successfully"
  }
  ```
- Validasi:
    - vendorId: required, must exist and belong to current user or admin
    - accountId: required, must exist and belong to the vendor
- Authorization: Bearer token (vendor owner or admin)
- Prioritas: Medium

## 3. CATEGORIES (Prioritas: Tinggi - MVP)

### 3.1 Category Endpoints

**Task 3.1.1: Create Category**
- Endpoint: `POST /categories`
- HTTP Method: POST
- Request Body:
  ```json
  {
    "parent_id": "integer",
    "name": "string",
    "description": "string",
    "icon": "string",
    "image_url": "string",
    "level": "integer",
    "sort_order": "integer",
    "is_active": "boolean"
  }
  ```
- Response Format (201):
  ```json
  {
    "id": "integer",
    "uuid": "string",
    "parent_id": "integer",
    "name": "string",
    "slug": "string",
    "description": "string",
    "icon": "string",
    "image_url": "string",
    "level": "integer",
    "sort_order": "integer",
    "is_active": "boolean",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - name: required, min 2 characters
    - parent_id: optional, must exist in categories table if provided
    - description: optional
    - icon: optional, max 100 characters
    - image_url: optional, valid URL if provided
    - level: optional, integer, min 1
    - sort_order: optional, integer, min 0
    - is_active: optional, boolean
- Authorization: Bearer token (admin only)
- Prioritas: High (MVP)

**Task 3.1.2: Get Category by ID**
- Endpoint: `GET /categories/{id}`
- HTTP Method: GET
- Path Params: id (integer)
- Request Body: None
- Response Format (200):
  ```json
  {
    "id": "integer",
    "uuid": "string",
    "parent_id": "integer",
    "name": "string",
    "slug": "string",
    "description": "string",
    "icon": "string",
    "image_url": "string",
    "level": "integer",
    "sort_order": "integer",
    "is_active": "boolean",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - id: required, must exist in database
- Authorization: Bearer token (any authenticated user)
- Prioritas: High (MVP)

**Task 3.1.3: Get Category by Slug**
- Endpoint: `GET /categories/slug/{slug}`
- HTTP Method: GET
- Path Params: slug (string)
- Request Body: None
- Response Format (200):
  ```json
  {
    "id": "integer",
    "uuid": "string",
    "parent_id": "integer",
    "name": "string",
    "slug": "string",
    "description": "string",
    "icon": "string",
    "image_url": "string",
    "level": "integer",
    "sort_order": "integer",
    "is_active": "boolean",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - slug: required, must exist in database
- Authorization: Bearer token (any authenticated user)
- Prioritas: High (MVP)

**Task 3.1.4: Update Category**
- Endpoint: `PUT /categories/{id}`
- HTTP Method: PUT
- Path Params: id (integer)
- Request Body:
  ```json
  {
    "parent_id": "integer",
    "name": "string",
    "description": "string",
    "icon": "string",
    "image_url": "string",
    "level": "integer",
    "sort_order": "integer",
    "is_active": "boolean"
  }
  ```
- Response Format (200):
  ```json
  {
    "id": "integer",
    "uuid": "string",
    "parent_id": "integer",
    "name": "string",
    "slug": "string",
    "description": "string",
    "icon": "string",
    "image_url": "string",
    "level": "integer",
    "sort_order": "integer",
    "is_active": "boolean",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - id: required, must exist in database
    - name: optional, min 2 characters
    - parent_id: optional, must exist in categories table if provided
    - description: optional
    - icon: optional, max 100 characters
    - image_url: optional, valid URL if provided
    - level: optional, integer, min 1
    - sort_order: optional, integer, min 0
    - is_active: optional, boolean
- Authorization: Bearer token (admin only)
- Prioritas: High (MVP)

**Task 3.1.5: Delete Category**
- Endpoint: `DELETE /categories/{id}`
- HTTP Method: DELETE
- Path Params: id (integer)
- Request Body: None
- Response Format (200):
  ```json
  {
    "message": "Category deleted successfully"
  }
  ```
- Validasi:
    - id: required, must exist in database
    - Check if category has child categories or products
- Authorization: Bearer token (admin only)
- Prioritas: Medium

**Task 3.1.6: List Categories**
- Endpoint: `GET /categories`
- HTTP Method: GET
- Query Params:
    - parent_id (integer, optional)
    - level (integer, optional)
    - is_active (boolean, optional)
    - sort (string, optional, default: 'sort_order')
    - order (string, optional, default: 'asc')
- Request Body: None
- Response Format (200):
  ```json
  {
    "data": [
      {
        "id": "integer",
        "uuid": "string",
        "parent_id": "integer",
        "name": "string",
        "slug": "string",
        "description": "string",
        "icon": "string",
        "image_url": "string",
        "level": "integer",
        "sort_order": "integer",
        "is_active": "boolean",
        "created_at": "timestamp",
        "updated_at": "timestamp"
      }
    ]
  }
  ```
- Validasi:
    - parent_id: optional, integer
    - level: optional, integer
    - is_active: optional, boolean
    - sort: optional, string, must be a valid column name
    - order: optional, string, must be 'asc' or 'desc'
- Authorization: Bearer token (any authenticated user)
- Prioritas: High (MVP)

**Task 3.1.7: Get Category Tree**
- Endpoint: `GET /categories/tree`
- HTTP Method: GET
- Query Params:
    - is_active (boolean, optional)
- Request Body: None
- Response Format (200):
  ```json
  {
    "data": [
      {
        "id": "integer",
        "uuid": "string",
        "parent_id": "integer",
        "name": "string",
        "slug": "string",
        "description": "string",
        "icon": "string",
        "image_url": "string",
        "level": "integer",
        "sort_order": "integer",
        "is_active": "boolean",
        "children": [
          {
            "id": "integer",
            "uuid": "string",
            "parent_id": "integer",
            "name": "string",
            "slug": "string",
            "description": "string",
            "icon": "string",
            "image_url": "string",
            "level": "integer",
            "sort_order": "integer",
            "is_active": "boolean",
            "children": []
          }
        ],
        "created_at": "timestamp",
        "updated_at": "timestamp"
      }
    ]
  }
  ```
- Validasi:
    - is_active: optional, boolean
- Authorization: Bearer token (any authenticated user)
- Prioritas: High (MVP)

## 4. PRODUCTS (Prioritas: Tinggi - MVP)

### 4.1 Product Endpoints

**Task 4.1.1: Create Product**
- Endpoint: `POST /products`
- HTTP Method: POST
- Request Body:
  ```json
  {
    "vendor_id": "integer",
    "name": "string",
    "description": "string",
    "short_description": "string",
    "price": "decimal",
    "compare_at_price": "decimal",
    "cost_per_item": "decimal",
    "stock": "integer",
    "low_stock_threshold": "integer",
    "track_inventory": "boolean",
    "allow_backorder": "boolean",
    "weight": "decimal",
    "length": "decimal",
    "width": "decimal",
    "height": "decimal",
    "meta_title": "string",
    "meta_description": "string",
    "meta_keywords": "string",
    "status": "string",
    "is_featured": "boolean",
    "category_ids": "array"
  }
  ```
- Response Format (201):
  ```json
  {
    "id": "integer",
    "uuid": "string",
    "vendor_id": "integer",
    "name": "string",
    "slug": "string",
    "sku": "string",
    "description": "string",
    "short_description": "string",
    "price": "decimal",
    "compare_at_price": "decimal",
    "cost_per_item": "decimal",
    "stock": "integer",
    "low_stock_threshold": "integer",
    "track_inventory": "boolean",
    "allow_backorder": "boolean",
    "weight": "decimal",
    "length": "decimal",
    "width": "decimal",
    "height": "decimal",
    "meta_title": "string",
    "meta_description": "string",
    "meta_keywords": "string",
    "status": "string",
    "is_featured": "boolean",
    "is_published": "boolean",
    "view_count": "integer",
    "sold_count": "integer",
    "rating_avg": "decimal",
    "review_count": "integer",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - vendor_id: required, must exist and belong to current user or admin
    - name: required, min 3 characters
    - description: optional
    - short_description: optional, max 500 characters
    - price: required, decimal, min 0
    - compare_at_price: optional, decimal, min 0
    - cost_per_item: optional, decimal, min 0
    - stock: required, integer, min 0
    - low_stock_threshold: optional, integer, min 0
    - track_inventory: optional, boolean
    - allow_backorder: optional, boolean
    - weight: optional, decimal, min 0
    - length: optional, decimal, min 0
    - width: optional, decimal, min 0
    - height: optional, decimal, min 0
    - meta_title: optional, max 255 characters
    - meta_description: optional, max 500 characters
    - meta_keywords: optional, max 255 characters
    - status: optional, must be one of 'draft', 'pending', 'active', 'rejected', 'out_of_stock'
    - is_featured: optional, boolean
    - category_ids: optional, array of integers, must exist in categories table
- Authorization: Bearer token (vendor owner, admin, or vendor users)
- Prioritas: High (MVP)

**Task 4.1.2: Get Product by ID**
- Endpoint: `GET /products/{id}`
- HTTP Method: GET
- Path Params: id (integer)
- Request Body: None
- Response Format (200):
  ```json
  {
    "id": "integer",
    "uuid": "string",
    "vendor_id": "integer",
    "vendor": {
      "id": "integer",
      "name": "string",
      "slug": "string"
    },
    "name": "string",
    "slug": "string",
    "sku": "string",
    "description": "string",
    "short_description": "string",
    "price": "decimal",
    "compare_at_price": "decimal",
    "cost_per_item": "decimal",
    "stock": "integer",
    "low_stock_threshold": "integer",
    "track_inventory": "boolean",
    "allow_backorder": "boolean",
    "weight": "decimal",
    "length": "decimal",
    "width": "decimal",
    "height": "decimal",
    "meta_title": "string",
    "meta_description": "string",
    "meta_keywords": "string",
    "status": "string",
    "is_featured": "boolean",
    "is_published": "boolean",
    "published_at": "timestamp",
    "view_count": "integer",
    "sold_count": "integer",
    "rating_avg": "decimal",
    "review_count": "integer",
    "categories": [
      {
        "id": "integer",
        "name": "string",
        "slug": "string"
      }
    ],
    "images": [
      {
        "id": "integer",
        "image_url": "string",
        "alt_text": "string",
        "sort_order": "integer",
        "is_primary": "boolean"
      }
    ],
    "variants": [
      {
        "id": "integer",
        "name": "string",
        "sku": "string",
        "price": "decimal",
        "stock": "integer",
        "image_url": "string",
        "is_active": "boolean"
      }
    ],
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - id: required, must exist in database
- Authorization: Bearer token (any authenticated user)
- Prioritas: High (MVP)

**Task 4.1.3: Get Product by Slug**
- Endpoint: `GET /products/slug/{slug}`
- HTTP Method: GET
- Path Params: slug (string)
- Request Body: None
- Response Format (200):
  ```json
  {
    "id": "integer",
    "uuid": "string",
    "vendor_id": "integer",
    "vendor": {
      "id": "integer",
      "name": "string",
      "slug": "string"
    },
    "name": "string",
    "slug": "string",
    "sku": "string",
    "description": "string",
    "short_description": "string",
    "price": "decimal",
    "compare_at_price": "decimal",
    "cost_per_item": "decimal",
    "stock": "integer",
    "low_stock_threshold": "integer",
    "track_inventory": "boolean",
    "allow_backorder": "boolean",
    "weight": "decimal",
    "length": "decimal",
    "width": "decimal",
    "height": "decimal",
    "meta_title": "string",
    "meta_description": "string",
    "meta_keywords": "string",
    "status": "string",
    "is_featured": "boolean",
    "is_published": "boolean",
    "published_at": "timestamp",
    "view_count": "integer",
    "sold_count": "integer",
    "rating_avg": "decimal",
    "review_count": "integer",
    "categories": [
      {
        "id": "integer",
        "name": "string",
        "slug": "string"
      }
    ],
    "images": [
      {
        "id": "integer",
        "image_url": "string",
        "alt_text": "string",
        "sort_order": "integer",
        "is_primary": "boolean"
      }
    ],
    "variants": [
      {
        "id": "integer",
        "name": "string",
        "sku": "string",
        "price": "decimal",
        "stock": "integer",
        "image_url": "string",
        "is_active": "boolean"
      }
    ],
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - slug: required, must exist in database
- Authorization: Bearer token (any authenticated user)
- Prioritas: High (MVP)

**Task 4.1.4: Update Product**
- Endpoint: `PUT /products/{id}`
- HTTP Method: PUT
- Path Params: id (integer)
- Request Body:
  ```json
  {
    "name": "string",
    "description": "string",
    "short_description": "string",
    "price": "decimal",
    "compare_at_price": "decimal",
    "cost_per_item": "decimal",
    "stock": "integer",
    "low_stock_threshold": "integer",
    "track_inventory": "boolean",
    "allow_backorder": "boolean",
    "weight": "decimal",
    "length": "decimal",
    "width": "decimal",
    "height": "decimal",
    "meta_title": "string",
    "meta_description": "string",
    "meta_keywords": "string",
    "status": "string",
    "is_featured": "boolean",
    "category_ids": "array"
  }
  ```
- Response Format (200):
  ```json
  {
    "id": "integer",
    "uuid": "string",
    "vendor_id": "integer",
    "name": "string",
    "slug": "string",
    "sku": "string",
    "description": "string",
    "short_description": "string",
    "price": "decimal",
    "compare_at_price": "decimal",
    "cost_per_item": "decimal",
    "stock": "integer",
    "low_stock_threshold": "integer",
    "track_inventory": "boolean",
    "allow_backorder": "boolean",
    "weight": "decimal",
    "length": "decimal",
    "width": "decimal",
    "height": "decimal",
    "meta_title": "string",
    "meta_description": "string",
    "meta_keywords": "string",
    "status": "string",
    "is_featured": "boolean",
    "is_published": "boolean",
    "published_at": "timestamp",
    "view_count": "integer",
    "sold_count": "integer",
    "rating_avg": "decimal",
    "review_count": "integer",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - id: required, must exist and belong to vendor of current user or admin
    - name: optional, min 3 characters
    - description: optional
    - short_description: optional, max 500 characters
    - price: optional, decimal, min 0
    - compare_at_price: optional, decimal, min 0
    - cost_per_item: optional, decimal, min 0
    - stock: optional, integer, min 0
    - low_stock_threshold: optional, integer, min 0
    - track_inventory: optional, boolean
    - allow_backorder: optional, boolean
    - weight: optional, decimal, min 0
    - length: optional, decimal, min 0
    - width: optional, decimal, min 0
    - height: optional, decimal, min 0
    - meta_title: optional, max 255 characters
    - meta_description: optional, max 500 characters
    - meta_keywords: optional, max 255 characters
    - status: optional, must be one of 'draft', 'pending', 'active', 'rejected', 'out_of_stock'
    - is_featured: optional, boolean
    - category_ids: optional, array of integers, must exist in categories table
- Authorization: Bearer token (vendor owner, admin, or vendor users)
- Prioritas: High (MVP)

**Task 4.1.5: Delete Product**
- Endpoint: `DELETE /products/{id}`
- HTTP Method: DELETE
- Path Params: id (integer)
- Request Body: None
- Response Format (200):
  ```json
  {
    "message": "Product deleted successfully"
  }
  ```
- Validasi:
    - id: required, must exist and belong to vendor of current user or admin
- Authorization: Bearer token (vendor owner or admin)
- Prioritas: Medium

**Task 4.1.6: List Products**
- Endpoint: `GET /products`
- HTTP Method: GET
- Query Params:
    - page (integer, optional)
    - limit (integer, optional)
    - search (string, optional)
    - vendor_id (integer, optional)
    - category_id (integer, optional)
    - min_price (decimal, optional)
    - max_price (decimal, optional)
    - status (string, optional)
    - is_featured (boolean, optional)
    - is_published (boolean, optional)
    - sort (string, optional, default: 'created_at')
    - order (string, optional, default: 'desc')
- Request Body: None
- Response Format (200):
  ```json
  {
    "data": [
      {
        "id": "integer",
        "uuid": "string",
        "vendor_id": "integer",
        "vendor": {
          "id": "integer",
          "name": "string",
          "slug": "string"
        },
        "name": "string",
        "slug": "string",
        "sku": "string",
        "price": "decimal",
        "compare_at_price": "decimal",
        "stock": "integer",
        "is_featured": "boolean",
        "is_published": "boolean",
        "view_count": "integer",
        "sold_count": "integer",
        "rating_avg": "decimal",
        "review_count": "integer",
        "primary_image": {
          "id": "integer",
          "image_url": "string"
        },
        "created_at": "timestamp",
        "updated_at": "timestamp"
      }
    ],
    "pagination": {
      "total": "integer",
      "page": "integer",
      "limit": "integer",
      "total_pages": "integer"
    }
  }
  ```
- Validasi:
    - page: optional, integer, min 1
    - limit: optional, integer, min 1, max 100
    - search: optional, string
    - vendor_id: optional, integer
    - category_id: optional, integer
    - min_price: optional, decimal, min 0
    - max_price: optional, decimal, min 0
    - status: optional, must be one of 'draft', 'pending', 'active', 'rejected', 'out_of_stock'
    - is_featured: optional, boolean
    - is_published: optional, boolean
    - sort: optional, string, must be a valid column name
    - order: optional, string, must be 'asc' or 'desc'
- Authorization: Bearer token (any authenticated user)
- Prioritas: High (MVP)

**Task 4.1.7: Update Product Status (Admin only)**
- Endpoint: `PUT /products/{id}/status`
- HTTP Method: PUT
- Path Params: id (integer)
- Request Body:
  ```json
  {
    "status": "string"
  }
  ```
- Response Format (200):
  ```json
  {
    "message": "Product status updated successfully",
    "status": "string"
  }
  ```
- Validasi:
    - id: required, must exist in database
    - status: required, must be one of 'draft', 'pending', 'active', 'rejected', 'out_of_stock'
- Authorization: Bearer token (admin only)
- Prioritas: High (MVP)

**Task 4.1.8: Publish Product**
- Endpoint: `PUT /products/{id}/publish`
- HTTP Method: PUT
- Path Params: id (integer)
- Request Body: None
- Response Format (200):
  ```json
  {
    "message": "Product published successfully",
    "is_published": "boolean",
    "published_at": "timestamp"
  }
  ```
- Validasi:
    - id: required, must exist and belong to vendor of current user or admin
- Authorization: Bearer token (vendor owner, admin, or vendor users)
- Prioritas: High (MVP)

**Task 4.1.9: Unpublish Product**
- Endpoint: `PUT /products/{id}/unpublish`
- HTTP Method: PUT
- Path Params: id (integer)
- Request Body: None
- Response Format (200):
  ```json
  {
    "message": "Product unpublished successfully",
    "is_published": "boolean"
  }
  ```
- Validasi:
    - id: required, must exist and belong to vendor of current user or admin
- Authorization: Bearer token (vendor owner, admin, or vendor users)
- Prioritas: High (MVP)

### 4.2 Product Variant Endpoints

**Task 4.2.1: Add Product Variant**
- Endpoint: `POST /products/{id}/variants`
- HTTP Method: POST
- Path Params: id (integer)
- Request Body:
  ```json
  {
    "name": "string",
    "sku": "string",
    "price": "decimal",
    "stock": "integer",
    "image_url": "string",
    "is_active": "boolean"
  }
  ```
- Response Format (201):
  ```json
  {
    "id": "integer",
    "product_id": "integer",
    "name": "string",
    "sku": "string",
    "price": "decimal",
    "stock": "integer",
    "image_url": "string",
    "is_active": "boolean",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - id: required, must exist and belong to vendor of current user or admin
    - name: required, min 3 characters
    - sku: optional, must be unique for this product
    - price: optional, decimal, min 0
    - stock: required, integer, min 0
    - image_url: optional, valid URL if provided
    - is_active: optional, boolean
- Authorization: Bearer token (vendor owner, admin, or vendor users)
- Prioritas: Medium

**Task 4.2.2: Get Product Variants**
- Endpoint: `GET /products/{id}/variants`
- HTTP Method: GET
- Path Params: id (integer)
- Request Body: None
- Response Format (200):
  ```json
  {
    "data": [
      {
        "id": "integer",
        "product_id": "integer",
        "name": "string",
        "sku": "string",
        "price": "decimal",
        "stock": "integer",
        "image_url": "string",
        "is_active": "boolean",
        "created_at": "timestamp",
        "updated_at": "timestamp"
      }
    ]
  }
  ```
- Validasi:
    - id: required, must exist in database
- Authorization: Bearer token (any authenticated user)
- Prioritas: Medium

**Task 4.2.3: Update Product Variant**
- Endpoint: `PUT /products/{productId}/variants/{variantId}`
- HTTP Method: PUT
- Path Params:
    - productId (integer)
    - variantId (integer)
- Request Body:
  ```json
  {
    "name": "string",
    "sku": "string",
    "price": "decimal",
    "stock": "integer",
    "image_url": "string",
    "is_active": "boolean"
  }
  ```
- Response Format (200):
  ```json
  {
    "id": "integer",
    "product_id": "integer",
    "name": "string",
    "sku": "string",
    "price": "decimal",
    "stock": "integer",
    "image_url": "string",
    "is_active": "boolean",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - productId: required, must exist and belong to vendor of current user or admin
    - variantId: required, must exist and belong to the product
    - name: optional, min 3 characters
    - sku: optional, must be unique for this product
    - price: optional, decimal, min 0
    - stock: optional, integer, min 0
    - image_url: optional, valid URL if provided
    - is_active: optional, boolean
- Authorization: Bearer token (vendor owner, admin, or vendor users)
- Prioritas: Medium

**Task 4.2.4: Delete Product Variant**
- Endpoint: `DELETE /products/{productId}/variants/{variantId}`
- HTTP Method: DELETE
- Path Params:
    - productId (integer)
    - variantId (integer)
- Request Body: None
- Response Format (200):
  ```json
  {
    "message": "Product variant deleted successfully"
  }
  ```
- Validasi:
    - productId: required, must exist and belong to vendor of current user or admin
    - variantId: required, must exist and belong to the product
- Authorization: Bearer token (vendor owner or admin)
- Prioritas: Medium

### 4.3 Product Image Endpoints

**Task 4.3.1: Upload Product Image**
- Endpoint: `POST /products/{id}/images`
- HTTP Method: POST
- Path Params: id (integer)
- Request Body: multipart/form-data with file
- Query Params:
    - alt_text (string, optional)
    - sort_order (integer, optional)
    - is_primary (boolean, optional)
- Response Format (201):
  ```json
  {
    "id": "integer",
    "product_id": "integer",
    "image_url": "string",
    "alt_text": "string",
    "sort_order": "integer",
    "is_primary": "boolean",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - id: required, must exist and belong to vendor of current user or admin
    - file: required, image file (jpg, png, etc), max size 2MB
    - alt_text: optional, max 255 characters
    - sort_order: optional, integer, min 0
    - is_primary: optional, boolean
- Authorization: Bearer token (vendor owner, admin, or vendor users)
- Prioritas: High (MVP)

**Task 4.3.2: Get Product Images**
- Endpoint: `GET /products/{id}/images`
- HTTP Method: GET
- Path Params: id (integer)
- Request Body: None
- Response Format (200):
  ```json
  {
    "data": [
      {
        "id": "integer",
        "product_id": "integer",
        "image_url": "string",
        "alt_text": "string",
        "sort_order": "integer",
        "is_primary": "boolean",
        "created_at": "timestamp",
        "updated_at": "timestamp"
      }
    ]
  }
  ```
- Validasi:
    - id: required, must exist in database
- Authorization: Bearer token (any authenticated user)
- Prioritas: High (MVP)

**Task 4.3.3: Update Product Image**
- Endpoint: `PUT /products/{productId}/images/{imageId}`
- HTTP Method: PUT
- Path Params:
    - productId (integer)
    - imageId (integer)
- Request Body:
  ```json
  {
    "alt_text": "string",
    "sort_order": "integer",
    "is_primary": "boolean"
  }
  ```
- Response Format (200):
  ```json
  {
    "id": "integer",
    "product_id": "integer",
    "image_url": "string",
    "alt_text": "string",
    "sort_order": "integer",
    "is_primary": "boolean",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - productId: required, must exist and belong to vendor of current user or admin
    - imageId: required, must exist and belong to the product
    - alt_text: optional, max 255 characters
    - sort_order: optional, integer, min 0
    - is_primary: optional, boolean
- Authorization: Bearer token (vendor owner, admin, or vendor users)
- Prioritas: Medium

**Task 4.3.4: Delete Product Image**
- Endpoint: `DELETE /products/{productId}/images/{imageId}`
- HTTP Method: DELETE
- Path Params:
    - productId (integer)
    - imageId (integer)
- Request Body: None
- Response Format (200):
  ```json
  {
    "message": "Product image deleted successfully"
  }
  ```
- Validasi:
    - productId: required, must exist and belong to vendor of current user or admin
    - imageId: required, must exist and belong to the product
- Authorization: Bearer token (vendor owner or admin)
- Prioritas: Medium

**Task 4.3.5: Set Primary Product Image**
- Endpoint: `PUT /products/{productId}/images/{imageId}/primary`
- HTTP Method: PUT
- Path Params:
    - productId (integer)
    - imageId (integer)
- Request Body: None
- Response Format (200):
  ```json
  {
    "message": "Primary product image updated successfully"
  }
  ```
- Validasi:
    - productId: required, must exist and belong to vendor of current user or admin
    - imageId: required, must exist and belong to the product
- Authorization: Bearer token (vendor owner, admin, or vendor users)
- Prioritas: Medium

### 4.4 Product Review Endpoints

**Task 4.4.1: Create Product Review**
- Endpoint: `POST /products/{id}/reviews`
- HTTP Method: POST
- Path Params: id (integer)
- Request Body:
  ```json
  {
    "rating": "integer",
    "title": "string",
    "comment": "string",
    "order_item_id": "integer"
  }
  ```
- Response Format (201):
  ```json
  {
    "id": "integer",
    "product_id": "integer",
    "user_id": "integer",
    "user": {
      "id": "integer",
      "name": "string"
    },
    "order_item_id": "integer",
    "rating": "integer",
    "title": "string",
    "comment": "string",
    "is_verified_purchase": "boolean",
    "is_approved": "boolean",
    "helpful_count": "integer",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - id: required, must exist in database
    - rating: required, integer, min 1, max 5
    - title: optional, max 255 characters
    - comment: optional
    - order_item_id: optional, must exist and belong to current user
- Authorization: Bearer token (any authenticated user)
- Prioritas: Medium

**Task 4.4.2: Get Product Reviews**
- Endpoint: `GET /products/{id}/reviews`
- HTTP Method: GET
- Path Params: id (integer)
- Query Params:
    - page (integer, optional)
    - limit (integer, optional)
    - rating (integer, optional)
    - sort (string, optional, default: 'created_at')
    - order (string, optional, default: 'desc')
- Request Body: None
- Response Format (200):
  ```json
  {
    "data": [
      {
        "id": "integer",
        "product_id": "integer",
        "user_id": "integer",
        "user": {
          "id": "integer",
          "name": "string"
        },
        "order_item_id": "integer",
        "rating": "integer",
        "title": "string",
        "comment": "string",
        "is_verified_purchase": "boolean",
        "is_approved": "boolean",
        "helpful_count": "integer",
        "created_at": "timestamp",
        "updated_at": "timestamp"
      }
    ],
    "pagination": {
      "total": "integer",
      "page": "integer",
      "limit": "integer",
      "total_pages": "integer"
    }
  }
  ```
- Validasi:
    - id: required, must exist in database
    - page: optional, integer, min 1
    - limit: optional, integer, min 1, max 100
    - rating: optional, integer, min 1, max 5
    - sort: optional, string, must be a valid column name
    - order: optional, string, must be 'asc' or 'desc'
- Authorization: Bearer token (any authenticated user)
- Prioritas: Medium

**Task 4.4.3: Update Product Review (Admin only)**
- Endpoint: `PUT /products/{productId}/reviews/{reviewId}`
- HTTP Method: PUT
- Path Params:
    - productId (integer)
    - reviewId (integer)
- Request Body:
  ```json
  {
    "is_approved": "boolean"
  }
  ```
- Response Format (200):
  ```json
  {
    "id": "integer",
    "product_id": "integer",
    "user_id": "integer",
    "user": {
      "id": "integer",
      "name": "string"
    },
    "order_item_id": "integer",
    "rating": "integer",
    "title": "string",
    "comment": "string",
    "is_verified_purchase": "boolean",
    "is_approved": "boolean",
    "helpful_count": "integer",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - productId: required, must exist in database
    - reviewId: required, must exist and belong to the product
    - is_approved: required, boolean
- Authorization: Bearer token (admin only)
- Prioritas: Medium

**Task 4.4.4: Delete Product Review (Admin only)**
- Endpoint: `DELETE /products/{productId}/reviews/{reviewId}`
- HTTP Method: DELETE
- Path Params:
    - productId (integer)
    - reviewId (integer)
- Request Body: None
- Response Format (200):
  ```json
  {
    "message": "Product review deleted successfully"
  }
  ```
- Validasi:
    - productId: required, must exist in database
    - reviewId: required, must exist and belong to the product
- Authorization: Bearer token (admin only)
- Prioritas: Medium

**Task 4.4.5: Mark Review as Helpful**
- Endpoint: `POST /products/{productId}/reviews/{reviewId}/helpful`
- HTTP Method: POST
- Path Params:
    - productId (integer)
    - reviewId (integer)
- Request Body: None
- Response Format (200):
  ```json
  {
    "message": "Review marked as helpful",
    "helpful_count": "integer"
  }
  ```
- Validasi:
    - productId: required, must exist in database
    - reviewId: required, must exist and belong to the product
- Authorization: Bearer token (any authenticated user)
- Prioritas: Low

## 5. SHOPPING CART (Prioritas: Tinggi - MVP)

### 5.1 Cart Endpoints

**Task 5.1.1: Get Cart**
- Endpoint: `GET /cart`
- HTTP Method: GET
- Request Body: None
- Response Format (200):
  ```json
  {
    "id": "integer",
    "user_id": "integer",
    "session_id": "string",
    "expires_at": "timestamp",
    "items": [
      {
        "id": "integer",
        "product_id": "integer",
        "product": {
          "id": "integer",
          "name": "string",
          "slug": "string",
          "price": "decimal",
          "stock": "integer",
          "primary_image": {
            "id": "integer",
            "image_url": "string"
          }
        },
        "variant_id": "integer",
        "variant": {
          "id": "integer",
          "name": "string",
          "price": "decimal",
          "stock": "integer"
        },
        "quantity": "integer",
        "price": "decimal",
        "created_at": "timestamp",
        "updated_at": "timestamp"
      }
    ],
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi: None
- Authorization: Bearer token (for authenticated users) or session_id (for guest users)
- Prioritas: High (MVP)

**Task 5.1.2: Add Item to Cart**
- Endpoint: `POST /cart/items`
- HTTP Method: POST
- Request Body:
  ```json
  {
    "product_id": "integer",
    "variant_id": "integer",
    "quantity": "integer"
  }
  ```
- Response Format (201):
  ```json
  {
    "id": "integer",
    "cart_id": "integer",
    "product_id": "integer",
    "product": {
      "id": "integer",
      "name": "string",
      "slug": "string",
      "price": "decimal",
      "stock": "integer"
    },
    "variant_id": "integer",
    "variant": {
      "id": "integer",
      "name": "string",
      "price": "decimal",
      "stock": "integer"
    },
    "quantity": "integer",
    "price": "decimal",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - product_id: required, must exist in database and be published
    - variant_id: optional, must exist and belong to the product if provided
    - quantity: required, integer, min 1
- Authorization: Bearer token (for authenticated users) or session_id (for guest users)
- Prioritas: High (MVP)

**Task 5.1.3: Update Cart Item**
- Endpoint: `PUT /cart/items/{id}`
- HTTP Method: PUT
- Path Params: id (integer)
- Request Body:
  ```json
  {
    "quantity": "integer"
  }
  ```
- Response Format (200):
  ```json
  {
    "id": "integer",
    "cart_id": "integer",
    "product_id": "integer",
    "product": {
      "id": "integer",
      "name": "string",
      "slug": "string",
      "price": "decimal",
      "stock": "integer"
    },
    "variant_id": "integer",
    "variant": {
      "id": "integer",
      "name": "string",
      "price": "decimal",
      "stock": "integer"
    },
    "quantity": "integer",
    "price": "decimal",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - id: required, must exist and belong to current user's cart or guest session
    - quantity: required, integer, min 1
- Authorization: Bearer token (for authenticated users) or session_id (for guest users)
- Prioritas: High (MVP)

**Task 5.1.4: Remove Item from Cart**
- Endpoint: `DELETE /cart/items/{id}`
- HTTP Method: DELETE
- Path Params: id (integer)
- Request Body: None
- Response Format (200):
  ```json
  {
    "message": "Item removed from cart successfully"
  }
  ```
- Validasi:
    - id: required, must exist and belong to current user's cart or guest session
- Authorization: Bearer token (for authenticated users) or session_id (for guest users)
- Prioritas: High (MVP)

**Task 5.1.5: Clear Cart**
- Endpoint: `DELETE /cart`
- HTTP Method: DELETE
- Request Body: None
- Response Format (200):
  ```json
  {
    "message": "Cart cleared successfully"
  }
  ```
- Validasi: None
- Authorization: Bearer token (for authenticated users) or session_id (for guest users)
- Prioritas: Medium

**Task 5.1.6: Merge Guest Cart to User Cart**
- Endpoint: `POST /cart/merge`
- HTTP Method: POST
- Request Body:
  ```json
  {
    "session_id": "string"
  }
  ```
- Response Format (200):
  ```json
  {
    "message": "Cart merged successfully",
    "cart": {
      "id": "integer",
      "user_id": "integer",
      "session_id": "string",
      "expires_at": "timestamp",
      "items": [
        {
          "id": "integer",
          "product_id": "integer",
          "variant_id": "integer",
          "quantity": "integer",
          "price": "decimal",
          "created_at": "timestamp",
          "updated_at": "timestamp"
        }
      ],
      "created_at": "timestamp",
      "updated_at": "timestamp"
    }
  }
  ```
- Validasi:
    - session_id: required, must exist in carts table
- Authorization: Bearer token (authenticated users only)
- Prioritas: Medium

## 6. ORDERS (Prioritas: Tinggi - MVP)

### 6.1 Order Endpoints

**Task 6.1.1: Create Order**
- Endpoint: `POST /orders`
- HTTP Method: POST
- Request Body:
  ```json
  {
    "shipping_address_id": "integer",
    "shipping_address": {
      "recipient_name": "string",
      "phone": "string",
      "address_line1": "string",
      "address_line2": "string",
      "city": "string",
      "state": "string",
      "postal_code": "string",
      "country": "string"
    },
    "shipping_method": "string",
    "payment_method": "string",
    "coupon_code": "string",
    "notes": "string"
  }
  ```
- Response Format (201):
  ```json
  {
    "id": "integer",
    "order_number": "string",
    "user_id": "integer",
    "shipping_name": "string",
    "shipping_phone": "string",
    "shipping_address_line1": "string",
    "shipping_address_line2": "string",
    "shipping_city": "string",
    "shipping_state": "string",
    "shipping_postal_code": "string",
    "shipping_country": "string",
    "subtotal": "decimal",
    "shipping_cost": "decimal",
    "tax_amount": "decimal",
    "discount_amount": "decimal",
    "total_amount": "decimal",
    "payment_method": "string",
    "payment_status": "string",
    "status": "string",
    "notes": "string",
    "coupon_code": "string",
    "items": [
      {
        "id": "integer",
        "vendor_id": "integer",
        "vendor": {
          "id": "integer",
          "name": "string"
        },
        "product_id": "integer",
        "product_name": "string",
        "variant_id": "integer",
        "variant_name": "string",
        "sku": "string",
        "quantity": "integer",
        "unit_price": "decimal",
        "subtotal": "decimal",
        "fulfillment_status": "string",
        "created_at": "timestamp"
      }
    ],
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - shipping_address_id: optional, must exist and belong to current user if provided
    - shipping_address: required if shipping_address_id not provided
        - recipient_name: required, min 3 characters
        - phone: required, valid phone format
        - address_line1: required, min 5 characters
        - city: required, min 2 characters
        - state: required, min 2 characters
        - postal_code: required, min 3 characters
        - country: optional, default: 'Indonesia'
    - shipping_method: required, must be a valid shipping method
    - payment_method: required, must be a valid payment method
    - coupon_code: optional, must exist and be valid if provided
    - notes: optional
- Authorization: Bearer token (any authenticated user)
- Prioritas: High (MVP)

**Task 6.1.2: Get Order by ID**
- Endpoint: `GET /orders/{id}`
- HTTP Method: GET
- Path Params: id (integer)
- Request Body: None
- Response Format (200):
  ```json
  {
    "id": "integer",
    "order_number": "string",
    "user_id": "integer",
    "shipping_name": "string",
    "shipping_phone": "string",
    "shipping_address_line1": "string",
    "shipping_address_line2": "string",
    "shipping_city": "string",
    "shipping_state": "string",
    "shipping_postal_code": "string",
    "shipping_country": "string",
    "subtotal": "decimal",
    "shipping_cost": "decimal",
    "tax_amount": "decimal",
    "discount_amount": "decimal",
    "total_amount": "decimal",
    "payment_method": "string",
    "payment_status": "string",
    "paid_at": "timestamp",
    "status": "string",
    "notes": "string",
    "coupon_code": "string",
    "items": [
      {
        "id": "integer",
        "vendor_id": "integer",
        "vendor": {
          "id": "integer",
          "name": "string"
        },
        "product_id": "integer",
        "product_name": "string",
        "variant_id": "integer",
        "variant_name": "string",
        "sku": "string",
        "quantity": "integer",
        "unit_price": "decimal",
        "subtotal": "decimal",
        "fulfillment_status": "string",
        "tracking_number": "string",
        "shipped_at": "timestamp",
        "delivered_at": "timestamp",
        "created_at": "timestamp"
      }
    ],
    "payments": [
      {
        "id": "integer",
        "payment_method": "string",
        "payment_gateway": "string",
        "transaction_id": "string",
        "amount": "decimal",
        "status": "string",
        "payment_proof_url": "string",
        "paid_at": "timestamp",
        "created_at": "timestamp"
      }
    ],
    "status_history": [
      {
        "id": "integer",
        "status": "string",
        "notes": "string",
        "created_by": "integer",
        "created_at": "timestamp"
      }
    ],
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - id: required, must exist and belong to current user or admin
- Authorization: Bearer token (order owner or admin)
- Prioritas: High (MVP)

**Task 6.1.3: Get Order by Order Number**
- Endpoint: `GET /orders/number/{orderNumber}`
- HTTP Method: GET
- Path Params: orderNumber (string)
- Request Body: None
- Response Format (200):
  ```json
  {
    "id": "integer",
    "order_number": "string",
    "user_id": "integer",
    "shipping_name": "string",
    "shipping_phone": "string",
    "shipping_address_line1": "string",
    "shipping_address_line2": "string",
    "shipping_city": "string",
    "shipping_state": "string",
    "shipping_postal_code": "string",
    "shipping_country": "string",
    "subtotal": "decimal",
    "shipping_cost": "decimal",
    "tax_amount": "decimal",
    "discount_amount": "decimal",
    "total_amount": "decimal",
    "payment_method": "string",
    "payment_status": "string",
    "paid_at": "timestamp",
    "status": "string",
    "notes": "string",
    "coupon_code": "string",
    "items": [
      {
        "id": "integer",
        "vendor_id": "integer",
        "vendor": {
          "id": "integer",
          "name": "string"
        },
        "product_id": "integer",
        "product_name": "string",
        "variant_id": "integer",
        "variant_name": "string",
        "sku": "string",
        "quantity": "integer",
        "unit_price": "decimal",
        "subtotal": "decimal",
        "fulfillment_status": "string",
        "tracking_number": "string",
        "shipped_at": "timestamp",
        "delivered_at": "timestamp",
        "created_at": "timestamp"
      }
    ],
    "payments": [
      {
        "id": "integer",
        "payment_method": "string",
        "payment_gateway": "string",
        "transaction_id": "string",
        "amount": "decimal",
        "status": "string",
        "payment_proof_url": "string",
        "paid_at": "timestamp",
        "created_at": "timestamp"
      }
    ],
    "status_history": [
      {
        "id": "integer",
        "status": "string",
        "notes": "string",
        "created_by": "integer",
        "created_at": "timestamp"
      }
    ],
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - orderNumber: required, must exist in database
- Authorization: Bearer token (order owner or admin)
- Prioritas: High (MVP)

**Task 6.1.4: List Orders**
- Endpoint: `GET /orders`
- HTTP Method: GET
- Query Params:
    - page (integer, optional)
    - limit (integer, optional)
    - status (string, optional)
    - payment_status (string, optional)
    - start_date (date, optional)
    - end_date (date, optional)
    - sort (string, optional, default: 'created_at')
    - order (string, optional, default: 'desc')
- Request Body: None
- Response Format (200):
  ```json
  {
    "data": [
      {
        "id": "integer",
        "order_number": "string",
        "user_id": "integer",
        "shipping_name": "string",
        "shipping_city": "string",
        "subtotal": "decimal",
        "shipping_cost": "decimal",
        "tax_amount": "decimal",
        "discount_amount": "decimal",
        "total_amount": "decimal",
        "payment_method": "string",
        "payment_status": "string",
        "status": "string",
        "created_at": "timestamp"
      }
    ],
    "pagination": {
      "total": "integer",
      "page": "integer",
      "limit": "integer",
      "total_pages": "integer"
    }
  }
  ```
- Validasi:
    - page: optional, integer, min 1
    - limit: optional, integer, min 1, max 100
    - status: optional, must be one of 'pending', 'processing', 'shipped', 'delivered', 'cancelled'
    - payment_status: optional, must be one of 'pending', 'paid', 'failed', 'refunded'
    - start_date: optional, date
    - end_date: optional, date
    - sort: optional, string, must be a valid column name
    - order: optional, string, must be 'asc' or 'desc'
- Authorization: Bearer token (any authenticated user, but will only show user's own orders unless admin)
- Prioritas: High (MVP)

**Task 6.1.5: Cancel Order**
- Endpoint: `PUT /orders/{id}/cancel`
- HTTP Method: PUT
- Path Params: id (integer)
- Request Body:
  ```json
  {
    "reason": "string"
  }
  ```
- Response Format (200):
  ```json
  {
    "message": "Order cancelled successfully",
    "status": "string"
  }
  ```
- Validasi:
    - id: required, must exist and belong to current user or admin
    - reason: optional
    - Check if order can be cancelled (status must be 'pending' or 'processing')
- Authorization: Bearer token (order owner or admin)
- Prioritas: High (MVP)

**Task 6.1.6: Update Order Status (Admin only)**
- Endpoint: `PUT /orders/{id}/status`
- HTTP Method: PUT
- Path Params: id (integer)
- Request Body:
  ```json
  {
    "status": "string",
    "notes": "string"
  }
  ```
- Response Format (200):
  ```json
  {
    "message": "Order status updated successfully",
    "status": "string"
  }
  ```
- Validasi:
    - id: required, must exist in database
    - status: required, must be one of 'pending', 'processing', 'shipped', 'delivered', 'cancelled'
    - notes: optional
- Authorization: Bearer token (admin only)
- Prioritas: High (MVP)

### 6.2 Order Item Endpoints

**Task 6.2.1: Update Order Item Fulfillment Status**
- Endpoint: `PUT /orders/{orderId}/items/{itemId}/status`
- HTTP Method: PUT
- Path Params:
    - orderId (integer)
    - itemId (integer)
- Request Body:
  ```json
  {
    "fulfillment_status": "string",
    "tracking_number": "string"
  }
  ```
- Response Format (200):
  ```json
  {
    "message": "Order item status updated successfully",
    "fulfillment_status": "string",
    "tracking_number": "string"
  }
  ```
- Validasi:
    - orderId: required, must exist in database
    - itemId: required, must exist and belong to the order
    - fulfillment_status: required, must be one of 'pending', 'processing', 'shipped', 'delivered', 'cancelled'
    - tracking_number: optional
- Authorization: Bearer token (vendor of the item or admin)
- Prioritas: High (MVP)

**Task 6.2.2: Get Vendor Order Items**
- Endpoint: `GET /vendors/{id}/order-items`
- HTTP Method: GET
- Path Params: id (integer)
- Query Params:
    - page (integer, optional)
    - limit (integer, optional)
    - fulfillment_status (string, optional)
    - start_date (date, optional)
    - end_date (date, optional)
    - sort (string, optional, default: 'created_at')
    - order (string, optional, default: 'desc')
- Request Body: None
- Response Format (200):
  ```json
  {
    "data": [
      {
        "id": "integer",
        "order_id": "integer",
        "order": {
          "id": "integer",
          "order_number": "string",
          "user_id": "integer",
          "status": "string",
          "created_at": "timestamp"
        },
        "vendor_id": "integer",
        "product_id": "integer",
        "product_name": "string",
        "variant_id": "integer",
        "variant_name": "string",
        "sku": "string",
        "quantity": "integer",
        "unit_price": "decimal",
        "subtotal": "decimal",
        "fulfillment_status": "string",
        "tracking_number": "string",
        "shipped_at": "timestamp",
        "delivered_at": "timestamp",
        "created_at": "timestamp"
      }
    ],
    "pagination": {
      "total": "integer",
      "page": "integer",
      "limit": "integer",
      "total_pages": "integer"
    }
  }
  ```
- Validasi:
    - id: required, must exist and belong to current user or admin
    - page: optional, integer, min 1
    - limit: optional, integer, min 1, max 100
    - fulfillment_status: optional, must be one of 'pending', 'processing', 'shipped', 'delivered', 'cancelled'
    - start_date: optional, date
    - end_date: optional, date
    - sort: optional, string, must be a valid column name
    - order: optional, string, must be 'asc' or 'desc'
- Authorization: Bearer token (vendor owner, admin, or vendor users)
- Prioritas: High (MVP)

## 7. PAYMENTS & TRANSACTIONS (Prioritas: Tinggi - MVP)

### 7.1 Payment Endpoints

**Task 7.1.1: Get Payment Methods**
- Endpoint: `GET /payments/methods`
- HTTP Method: GET
- Request Body: None
- Response Format (200):
  ```json
  {
    "data": [
      {
        "code": "string",
        "name": "string",
        "description": "string",
        "icon": "string",
        "is_active": "boolean"
      }
    ]
  }
  ```
- Validasi: None
- Authorization: Bearer token (any authenticated user)
- Prioritas: High (MVP)

**Task 7.1.2: Process Payment**
- Endpoint: `POST /payments/process`
- HTTP Method: POST
- Request Body:
  ```json
  {
    "order_id": "integer",
    "payment_method": "string",
    "payment_gateway": "string"
  }
  ```
- Response Format (201):
  ```json
  {
    "id": "integer",
    "order_id": "integer",
    "payment_method": "string",
    "payment_gateway": "string",
    "transaction_id": "string",
    "amount": "decimal",
    "status": "string",
    "payment_proof_url": "string",
    "expired_at": "timestamp",
    "payment_instructions": "object",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - order_id: required, must exist and belong to current user
    - payment_method: required, must be a valid payment method
    - payment_gateway: optional, must be a valid payment gateway
- Authorization: Bearer token (any authenticated user)
- Prioritas: High (MVP)

**Task 7.1.3: Upload Payment Proof**
- Endpoint: `POST /payments/{id}/proof`
- HTTP Method: POST
- Path Params: id (integer)
- Request Body: multipart/form-data with file
- Response Format (200):
  ```json
  {
    "message": "Payment proof uploaded successfully",
    "payment_proof_url": "string"
  }
  ```
- Validasi:
    - id: required, must exist and belong to order of current user
    - file: required, image file (jpg, png, etc), max size 2MB
- Authorization: Bearer token (any authenticated user)
- Prioritas: High (MVP)

**Task 7.1.4: Get Payment by ID**
- Endpoint: `GET /payments/{id}`
- HTTP Method: GET
- Path Params: id (integer)
- Request Body: None
- Response Format (200):
  ```json
  {
    "id": "integer",
    "order_id": "integer",
    "order": {
      "id": "integer",
      "order_number": "string"
    },
    "payment_method": "string",
    "payment_gateway": "string",
    "transaction_id": "string",
    "amount": "decimal",
    "status": "string",
    "payment_proof_url": "string",
    "paid_at": "timestamp",
    "expired_at": "timestamp",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - id: required, must exist in database
- Authorization: Bearer token (payment owner or admin)
- Prioritas: Medium

**Task 7.1.5: Update Payment Status (Admin only)**
- Endpoint: `PUT /payments/{id}/status`
- HTTP Method: PUT
- Path Params: id (integer)
- Request Body:
  ```json
  {
    "status": "string"
  }
  ```
- Response Format (200):
  ```json
  {
    "message": "Payment status updated successfully",
    "status": "string"
  }
  ```
- Validasi:
    - id: required, must exist in database
    - status: required, must be one of 'pending', 'paid', 'failed', 'refunded'
- Authorization: Bearer token (admin only)
- Prioritas: High (MVP)

### 7.2 Vendor Transaction Endpoints

**Task 7.2.1: Get Vendor Transactions**
- Endpoint: `GET /vendors/{id}/transactions`
- HTTP Method: GET
- Path Params: id (integer)
- Query Params:
    - page (integer, optional)
    - limit (integer, optional)
    - type (string, optional)
    - status (string, optional)
    - start_date (date, optional)
    - end_date (date, optional)
    - sort (string, optional, default: 'created_at')
    - order (string, optional, default: 'desc')
- Request Body: None
- Response Format (200):
  ```json
  {
    "data": [
      {
        "id": "integer",
        "vendor_id": "integer",
        "order_item_id": "integer",
        "order_item": {
          "id": "integer",
          "order_id": "integer",
          "product_name": "string"
        },
        "type": "string",
        "amount": "decimal",
        "commission_rate": "decimal",
        "commission_amount": "decimal",
        "net_amount": "decimal",
        "status": "string",
        "description": "string",
        "created_at": "timestamp",
        "updated_at": "timestamp"
      }
    ],
    "pagination": {
      "total": "integer",
      "page": "integer",
      "limit": "integer",
      "total_pages": "integer"
    }
  }
  ```
- Validasi:
    - id: required, must exist and belong to current user or admin
    - page: optional, integer, min 1
    - limit: optional, integer, min 1, max 100
    - type: optional, must be one of 'sale', 'refund', 'commission', 'withdrawal'
    - status: optional, must be one of 'pending', 'completed', 'failed'
    - start_date: optional, date
    - end_date: optional, date
    - sort: optional, string, must be a valid column name
    - order: optional, string, must be 'asc' or 'desc'
- Authorization: Bearer token (vendor owner, admin, or vendor users)
- Prioritas: Medium

**Task 7.2.2: Get Vendor Transaction Summary**
- Endpoint: `GET /vendors/{id}/transactions/summary`
- HTTP Method: GET
- Path Params: id (integer)
- Query Params:
    - start_date (date, optional)
    - end_date (date, optional)
- Request Body: None
- Response Format (200):
  ```json
  {
    "total_sales": "decimal",
    "total_refunds": "decimal",
    "total_commissions": "decimal",
    "total_withdrawals": "decimal",
    "net_earnings": "decimal",
    "pending_amount": "decimal",
    "completed_amount": "decimal",
    "failed_amount": "decimal"
  }
  ```
- Validasi:
    - id: required, must exist and belong to current user or admin
    - start_date: optional, date
    - end_date: optional, date
- Authorization: Bearer token (vendor owner, admin, or vendor users)
- Prioritas: Medium

### 7.3 Vendor Withdrawal Endpoints

**Task 7.3.1: Create Vendor Withdrawal**
- Endpoint: `POST /vendors/{id}/withdrawals`
- HTTP Method: POST
- Path Params: id (integer)
- Request Body:
  ```json
  {
    "bank_account_id": "integer",
    "amount": "decimal",
    "notes": "string"
  }
  ```
- Response Format (201):
  ```json
  {
    "id": "integer",
    "vendor_id": "integer",
    "bank_account_id": "integer",
    "bank_account": {
      "id": "integer",
      "bank_name": "string",
      "account_number": "string",
      "account_holder": "string"
    },
    "amount": "decimal",
    "admin_fee": "decimal",
    "net_amount": "decimal",
    "status": "string",
    "notes": "string",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - id: required, must exist and belong to current user or admin
    - bank_account_id: required, must exist and belong to the vendor
    - amount: required, decimal, min minimum withdrawal amount
    - notes: optional
- Authorization: Bearer token (vendor owner or admin)
- Prioritas: Medium

**Task 7.3.2: Get Vendor Withdrawals**
- Endpoint: `GET /vendors/{id}/withdrawals`
- HTTP Method: GET
- Path Params: id (integer)
- Query Params:
    - page (integer, optional)
    - limit (integer, optional)
    - status (string, optional)
    - start_date (date, optional)
    - end_date (date, optional)
    - sort (string, optional, default: 'created_at')
    - order (string, optional, default: 'desc')
- Request Body: None
- Response Format (200):
  ```json
  {
    "data": [
      {
        "id": "integer",
        "vendor_id": "integer",
        "bank_account_id": "integer",
        "bank_account": {
          "id": "integer",
          "bank_name": "string",
          "account_number": "string",
          "account_holder": "string"
        },
        "amount": "decimal",
        "admin_fee": "decimal",
        "net_amount": "decimal",
        "status": "string",
        "notes": "string",
        "processed_by": "integer",
        "processed_at": "timestamp",
        "created_at": "timestamp",
        "updated_at": "timestamp"
      }
    ],
    "pagination": {
      "total": "integer",
      "page": "integer",
      "limit": "integer",
      "total_pages": "integer"
    }
  }
  ```
- Validasi:
    - id: required, must exist and belong to current user or admin
    - page: optional, integer, min 1
    - limit: optional, integer, min 1, max 100
    - status: optional, must be one of 'pending', 'processing', 'completed', 'rejected'
    - start_date: optional, date
    - end_date: optional, date
    - sort: optional, string, must be a valid column name
    - order: optional, string, must be 'asc' or 'desc'
- Authorization: Bearer token (vendor owner, admin, or vendor users)
- Prioritas: Medium

**Task 7.3.3: Get Withdrawal by ID**
- Endpoint: `GET /withdrawals/{id}`
- HTTP Method: GET
- Path Params: id (integer)
- Request Body: None
- Response Format (200):
  ```json
  {
    "id": "integer",
    "vendor_id": "integer",
    "vendor": {
      "id": "integer",
      "name": "string"
    },
    "bank_account_id": "integer",
    "bank_account": {
      "id": "integer",
      "bank_name": "string",
      "account_number": "string",
      "account_holder": "string"
    },
    "amount": "decimal",
    "admin_fee": "decimal",
    "net_amount": "decimal",
    "status": "string",
    "notes": "string",
    "processed_by": "integer",
    "processed_at": "timestamp",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - id: required, must exist in database
- Authorization: Bearer token (vendor owner or admin)
- Prioritas: Medium

**Task 7.3.4: Process Withdrawal (Admin only)**
- Endpoint: `PUT /withdrawals/{id}/process`
- HTTP Method: PUT
- Path Params: id (integer)
- Request Body:
  ```json
  {
    "status": "string",
    "notes": "string"
  }
  ```
- Response Format (200):
  ```json
  {
    "message": "Withdrawal processed successfully",
    "status": "string",
    "processed_at": "timestamp"
  }
  ```
- Validasi:
    - id: required, must exist in database
    - status: required, must be one of 'completed', 'rejected'
    - notes: optional
- Authorization: Bearer token (admin only)
- Prioritas: Medium

## 8. SHIPPING (Prioritas: Medium)

### 8.1 Shipping Method Endpoints

**Task 8.1.1: Get Shipping Methods**
- Endpoint: `GET /shipping/methods`
- HTTP Method: GET
- Query Params:
    - vendor_id (integer, optional)
    - is_active (boolean, optional)
- Request Body: None
- Response Format (200):
  ```json
  {
    "data": [
      {
        "id": "integer",
        "vendor_id": "integer",
        "name": "string",
        "code": "string",
        "description": "string",
        "carrier": "string",
        "base_cost": "decimal",
        "cost_per_kg": "decimal",
        "estimated_days": "string",
        "is_active": "boolean",
        "created_at": "timestamp",
        "updated_at": "timestamp"
      }
    ]
  }
  ```
- Validasi:
    - vendor_id: optional, integer
    - is_active: optional, boolean
- Authorization: Bearer token (any authenticated user)
- Prioritas: Medium

**Task 8.1.2: Create Shipping Method (Vendor or Admin)**
- Endpoint: `POST /shipping/methods`
- HTTP Method: POST
- Request Body:
  ```json
  {
    "vendor_id": "integer",
    "name": "string",
    "code": "string",
    "description": "string",
    "carrier": "string",
    "base_cost": "decimal",
    "cost_per_kg": "decimal",
    "estimated_days": "string",
    "is_active": "boolean"
  }
  ```
- Response Format (201):
  ```json
  {
    "id": "integer",
    "vendor_id": "integer",
    "name": "string",
    "code": "string",
    "description": "string",
    "carrier": "string",
    "base_cost": "decimal",
    "cost_per_kg": "decimal",
    "estimated_days": "string",
    "is_active": "boolean",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - vendor_id: optional, must exist and belong to current user if provided
    - name: required, min 2 characters
    - code: required, min 2 characters, must be unique for vendor
    - description: optional
    - carrier: optional, min 2 characters
    - base_cost: required, decimal, min 0
    - cost_per_kg: required, decimal, min 0
    - estimated_days: optional, max 50 characters
    - is_active: optional, boolean
- Authorization: Bearer token (vendor owner or admin)
- Prioritas: Medium

**Task 8.1.3: Update Shipping Method (Vendor or Admin)**
- Endpoint: `PUT /shipping/methods/{id}`
- HTTP Method: PUT
- Path Params: id (integer)
- Request Body:
  ```json
  {
    "name": "string",
    "code": "string",
    "description": "string",
    "carrier": "string",
    "base_cost": "decimal",
    "cost_per_kg": "decimal",
    "estimated_days": "string",
    "is_active": "boolean"
  }
  ```
- Response Format (200):
  ```json
  {
    "id": "integer",
    "vendor_id": "integer",
    "name": "string",
    "code": "string",
    "description": "string",
    "carrier": "string",
    "base_cost": "decimal",
    "cost_per_kg": "decimal",
    "estimated_days": "string",
    "is_active": "boolean",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - id: required, must exist and belong to vendor of current user or admin
    - name: optional, min 2 characters
    - code: optional, min 2 characters, must be unique for vendor
    - description: optional
    - carrier: optional, min 2 characters
    - base_cost: optional, decimal, min 0
    - cost_per_kg: optional, decimal, min 0
    - estimated_days: optional, max 50 characters
    - is_active: optional, boolean
- Authorization: Bearer token (vendor owner or admin)
- Prioritas: Medium

**Task 8.1.4: Delete Shipping Method (Vendor or Admin)**
- Endpoint: `DELETE /shipping/methods/{id}`
- HTTP Method: DELETE
- Path Params: id (integer)
- Request Body: None
- Response Format (200):
  ```json
  {
    "message": "Shipping method deleted successfully"
  }
  ```
- Validasi:
    - id: required, must exist and belong to vendor of current user or admin
- Authorization: Bearer token (vendor owner or admin)
- Prioritas: Medium

**Task 8.1.5: Calculate Shipping Cost**
- Endpoint: `POST /shipping/calculate`
- HTTP Method: POST
- Request Body:
  ```json
  {
    "vendor_id": "integer",
    "destination": {
      "city": "string",
      "state": "string",
      "postal_code": "string",
      "country": "string"
    },
    "weight": "decimal",
    "items": [
      {
        "product_id": "integer",
        "quantity": "integer"
      }
    ]
  }
  ```
- Response Format (200):
  ```json
  {
    "data": [
      {
        "id": "integer",
        "vendor_id": "integer",
        "name": "string",
        "code": "string",
        "description": "string",
        "carrier": "string",
        "cost": "decimal",
        "estimated_days": "string"
      }
    ]
  }
  ```
- Validasi:
    - vendor_id: required, must exist in database
    - destination: required
        - city: required, min 2 characters
        - state: required, min 2 characters
        - postal_code: required, min 3 characters
        - country: optional, default: 'Indonesia'
    - weight: required, decimal, min 0
    - items: required, array
        - product_id: required, must exist in database
        - quantity: required, integer, min 1
- Authorization: Bearer token (any authenticated user)
- Prioritas: Medium

## 9. PROMOTIONS & DISCOUNTS (Prioritas: Medium)

### 9.1 Coupon Endpoints

**Task 9.1.1: Create Coupon (Vendor or Admin)**
- Endpoint: `POST /coupons`
- HTTP Method: POST
- Request Body:
  ```json
  {
    "vendor_id": "integer",
    "code": "string",
    "description": "string",
    "discount_type": "string",
    "discount_value": "decimal",
    "min_purchase": "decimal",
    "max_discount": "decimal",
    "usage_limit": "integer",
    "usage_limit_per_user": "integer",
    "valid_from": "timestamp",
    "valid_until": "timestamp",
    "is_active": "boolean"
  }
  ```
- Response Format (201):
  ```json
  {
    "id": "integer",
    "vendor_id": "integer",
    "code": "string",
    "description": "string",
    "discount_type": "string",
    "discount_value": "decimal",
    "min_purchase": "decimal",
    "max_discount": "decimal",
    "usage_limit": "integer",
    "usage_limit_per_user": "integer",
    "used_count": "integer",
    "valid_from": "timestamp",
    "valid_until": "timestamp",
    "is_active": "boolean",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - vendor_id: optional, must exist and belong to current user if provided
    - code: required, min 2 characters, must be unique
    - description: optional
    - discount_type: required, must be one of 'percentage', 'fixed'
    - discount_value: required, decimal, min 0
    - min_purchase: optional, decimal, min 0
    - max_discount: optional, decimal, min 0
    - usage_limit: optional, integer, min 1
    - usage_limit_per_user: optional, integer, min 1
    - valid_from: required, timestamp
    - valid_until: required, timestamp, must be after valid_from
    - is_active: optional, boolean
- Authorization: Bearer token (vendor owner or admin)
- Prioritas: Medium

**Task 9.1.2: Get Coupon by Code**
- Endpoint: `GET /coupons/{code}`
- HTTP Method: GET
- Path Params: code (string)
- Request Body: None
- Response Format (200):
  ```json
  {
    "id": "integer",
    "vendor_id": "integer",
    "vendor": {
      "id": "integer",
      "name": "string"
    },
    "code": "string",
    "description": "string",
    "discount_type": "string",
    "discount_value": "decimal",
    "min_purchase": "decimal",
    "max_discount": "decimal",
    "usage_limit": "integer",
    "usage_limit_per_user": "integer",
    "used_count": "integer",
    "valid_from": "timestamp",
    "valid_until": "timestamp",
    "is_active": "boolean",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - code: required, must exist in database
- Authorization: Bearer token (any authenticated user)
- Prioritas: Medium

**Task 9.1.3: Update Coupon (Vendor or Admin)**
- Endpoint: `PUT /coupons/{id}`
- HTTP Method: PUT
- Path Params: id (integer)
- Request Body:
  ```json
  {
    "description": "string",
    "discount_type": "string",
    "discount_value": "decimal",
    "min_purchase": "decimal",
    "max_discount": "decimal",
    "usage_limit": "integer",
    "usage_limit_per_user": "integer",
    "valid_from": "timestamp",
    "valid_until": "timestamp",
    "is_active": "boolean"
  }
  ```
- Response Format (200):
  ```json
  {
    "id": "integer",
    "vendor_id": "integer",
    "code": "string",
    "description": "string",
    "discount_type": "string",
    "discount_value": "decimal",
    "min_purchase": "decimal",
    "max_discount": "decimal",
    "usage_limit": "integer",
    "usage_limit_per_user": "integer",
    "used_count": "integer",
    "valid_from": "timestamp",
    "valid_until": "timestamp",
    "is_active": "boolean",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - id: required, must exist and belong to vendor of current user or admin
    - description: optional
    - discount_type: optional, must be one of 'percentage', 'fixed'
    - discount_value: optional, decimal, min 0
    - min_purchase: optional, decimal, min 0
    - max_discount: optional, decimal, min 0
    - usage_limit: optional, integer, min 1
    - usage_limit_per_user: optional, integer, min 1
    - valid_from: optional, timestamp
    - valid_until: optional, timestamp, must be after valid_from
    - is_active: optional, boolean
- Authorization: Bearer token (vendor owner or admin)
- Prioritas: Medium

**Task 9.1.4: Delete Coupon (Vendor or Admin)**
- Endpoint: `DELETE /coupons/{id}`
- HTTP Method: DELETE
- Path Params: id (integer)
- Request Body: None
- Response Format (200):
  ```json
  {
    "message": "Coupon deleted successfully"
  }
  ```
- Validasi:
    - id: required, must exist and belong to vendor of current user or admin
- Authorization: Bearer token (vendor owner or admin)
- Prioritas: Medium

**Task 9.1.5: List Coupons**
- Endpoint: `GET /coupons`
- HTTP Method: GET
- Query Params:
    - page (integer, optional)
    - limit (integer, optional)
    - vendor_id (integer, optional)
    - is_active (boolean, optional)
    - valid_only (boolean, optional)
    - sort (string, optional, default: 'created_at')
    - order (string, optional, default: 'desc')
- Request Body: None
- Response Format (200):
  ```json
  {
    "data": [
      {
        "id": "integer",
        "vendor_id": "integer",
        "vendor": {
          "id": "integer",
          "name": "string"
        },
        "code": "string",
        "description": "string",
        "discount_type": "string",
        "discount_value": "decimal",
        "min_purchase": "decimal",
        "max_discount": "decimal",
        "usage_limit": "integer",
        "usage_limit_per_user": "integer",
        "used_count": "integer",
        "valid_from": "timestamp",
        "valid_until": "timestamp",
        "is_active": "boolean",
        "created_at": "timestamp",
        "updated_at": "timestamp"
      }
    ],
    "pagination": {
      "total": "integer",
      "page": "integer",
      "limit": "integer",
      "total_pages": "integer"
    }
  }
  ```
- Validasi:
    - page: optional, integer, min 1
    - limit: optional, integer, min 1, max 100
    - vendor_id: optional, integer
    - is_active: optional, boolean
    - valid_only: optional, boolean
    - sort: optional, string, must be a valid column name
    - order: optional, string, must be 'asc' or 'desc'
- Authorization: Bearer token (any authenticated user)
- Prioritas: Medium

**Task 9.1.6: Validate Coupon**
- Endpoint: `POST /coupons/validate`
- HTTP Method: POST
- Request Body:
  ```json
  {
    "code": "string",
    "user_id": "integer",
    "subtotal": "decimal",
    "vendor_id": "integer"
  }
  ```
- Response Format (200):
  ```json
  {
    "valid": "boolean",
    "message": "string",
    "discount_amount": "decimal",
    "coupon": {
      "id": "integer",
      "vendor_id": "integer",
      "code": "string",
      "description": "string",
      "discount_type": "string",
      "discount_value": "decimal",
      "min_purchase": "decimal",
      "max_discount": "decimal",
      "usage_limit": "integer",
      "usage_limit_per_user": "integer",
      "used_count": "integer",
      "valid_from": "timestamp",
      "valid_until": "timestamp",
      "is_active": "boolean"
    }
  }
  ```
- Validasi:
    - code: required, min 2 characters
    - user_id: required, integer
    - subtotal: required, decimal, min 0
    - vendor_id: optional, integer
- Authorization: Bearer token (any authenticated user)
- Prioritas: Medium

## 10. NOTIFICATIONS (Prioritas: Rendah)

### 10.1 Notification Endpoints

**Task 10.1.1: Get User Notifications**
- Endpoint: `GET /notifications`
- HTTP Method: GET
- Query Params:
    - page (integer, optional)
    - limit (integer, optional)
    - is_read (boolean, optional)
    - type (string, optional)
    - sort (string, optional, default: 'created_at')
    - order (string, optional, default: 'desc')
- Request Body: None
- Response Format (200):
  ```json
  {
    "data": [
      {
        "id": "integer",
        "type": "string",
        "title": "string",
        "message": "string",
        "data": "object",
        "is_read": "boolean",
        "read_at": "timestamp",
        "created_at": "timestamp"
      }
    ],
    "pagination": {
      "total": "integer",
      "page": "integer",
      "limit": "integer",
      "total_pages": "integer"
    },
    "unread_count": "integer"
  }
  ```
- Validasi:
    - page: optional, integer, min 1
    - limit: optional, integer, min 1, max 100
    - is_read: optional, boolean
    - type: optional, string
    - sort: optional, string, must be a valid column name
    - order: optional, string, must be 'asc' or 'desc'
- Authorization: Bearer token (any authenticated user)
- Prioritas: Low

**Task 10.1.2: Mark Notification as Read**
- Endpoint: `PUT /notifications/{id}/read`
- HTTP Method: PUT
- Path Params: id (integer)
- Request Body: None
- Response Format (200):
  ```json
  {
    "message": "Notification marked as read",
    "is_read": "boolean",
    "read_at": "timestamp"
  }
  ```
- Validasi:
    - id: required, must exist and belong to current user
- Authorization: Bearer token (any authenticated user)
- Prioritas: Low

**Task 10.1.3: Mark All Notifications as Read**
- Endpoint: `PUT /notifications/read-all`
- HTTP Method: PUT
- Request Body: None
- Response Format (200):
  ```json
  {
    "message": "All notifications marked as read"
  }
  ```
- Validasi: None
- Authorization: Bearer token (any authenticated user)
- Prioritas: Low

**Task 10.1.4: Delete Notification**
- Endpoint: `DELETE /notifications/{id}`
- HTTP Method: DELETE
- Path Params: id (integer)
- Request Body: None
- Response Format (200):
  ```json
  {
    "message": "Notification deleted successfully"
  }
  ```
- Validasi:
    - id: required, must exist and belong to current user
- Authorization: Bearer token (any authenticated user)
- Prioritas: Low

## 11. MEDIA (Prioritas: Rendah)

### 11.1 Media Endpoints

**Task 11.1.1: Upload Media**
- Endpoint: `POST /media/upload`
- HTTP Method: POST
- Request Body: multipart/form-data with file
- Query Params:
    - entity_type (string, optional)
    - entity_id (integer, optional)
- Response Format (201):
  ```json
  {
    "id": "integer",
    "uuid": "string",
    "vendor_id": "integer",
    "user_id": "integer",
    "entity_type": "string",
    "entity_id": "integer",
    "file_path": "string",
    "file_name": "string",
    "file_size": "integer",
    "mime_type": "string",
    "width": "integer",
    "height": "integer",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - file: required, file, max size 10MB
    - entity_type: optional, string
    - entity_id: optional, integer
- Authorization: Bearer token (any authenticated user)
- Prioritas: Low

**Task 11.1.2: Get Media by ID**
- Endpoint: `GET /media/{id}`
- HTTP Method: GET
- Path Params: id (integer)
- Request Body: None
- Response Format (200):
  ```json
  {
    "id": "integer",
    "uuid": "string",
    "vendor_id": "integer",
    "user_id": "integer",
    "entity_type": "string",
    "entity_id": "integer",
    "file_path": "string",
    "file_name": "string",
    "file_size": "integer",
    "mime_type": "string",
    "width": "integer",
    "height": "integer",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
  ```
- Validasi:
    - id: required, must exist in database
- Authorization: Bearer token (media owner or admin)
- Prioritas: Low

**Task 11.1.3: Delete Media**
- Endpoint: `DELETE /media/{id}`
- HTTP Method: DELETE
- Path Params: id (integer)
- Request Body: None
- Response Format (200):
  ```json
  {
    "message": "Media deleted successfully"
  }
  ```
- Validasi:
    - id: required, must exist and belong to current user or admin
- Authorization: Bearer token (media owner or admin)
- Prioritas: Low

## 12. ADMIN & LOGS (Prioritas: Rendah)

### 12.1 Admin Activity Log Endpoints

**Task 12.1.1: Get Admin Activity Logs (Admin only)**
- Endpoint: `GET /admin/logs`
- HTTP Method: GET
- Query Params:
    - page (integer, optional)
    - limit (integer, optional)
    - user_id (integer, optional)
    - action (string, optional)
    - entity_type (string, optional)
    - entity_id (integer, optional)
    - start_date (date, optional)
    - end_date (date, optional)
    - sort (string, optional, default: 'created_at')
    - order (string, optional, default: 'desc')
- Request Body: None
- Response Format (200):
  ```json
  {
    "data": [
      {
        "id": "integer",
        "user_id": "integer",
        "user": {
          "id": "integer",
          "name": "string",
          "email": "string"
        },
        "action": "string",
        "entity_type": "string",
        "entity_id": "integer",
        "old_values": "object",
        "new_values": "object",
        "ip_address": "string",
        "user_agent": "string",
        "created_at": "timestamp"
      }
    ],
    "pagination": {
      "total": "integer",
      "page": "integer",
      "limit": "integer",
      "total_pages": "integer"
    }
  }
  ```
- Validasi:
    - page: optional, integer, min 1
    - limit: optional, integer, min 1, max 100
    - user_id: optional, integer
    - action: optional, string
    - entity_type: optional, string
    - entity_id: optional, integer
    - start_date: optional, date
    - end_date: optional, date
    - sort: optional, string, must be a valid column name
    - order: optional, string, must be 'asc' or 'desc'
- Authorization: Bearer token (admin only)
- Prioritas: Low

### 12.2 Dashboard Endpoints

**Task 12.2.1: Get Dashboard Stats (Admin only)**
- Endpoint: `GET /admin/dashboard`
- HTTP Method: GET
- Query Params:
    - start_date (date, optional)
    - end_date (date, optional)
- Request Body: None
- Response Format (200):
  ```json
  {
    "users": {
      "total": "integer",
      "new": "integer",
      "active": "integer"
    },
    "vendors": {
      "total": "integer",
      "pending": "integer",
      "active": "integer"
    },
    "products": {
      "total": "integer",
      "pending": "integer",
      "active": "integer"
    },
    "orders": {
      "total": "integer",
      "pending": "integer",
      "processing": "integer",
      "shipped": "integer",
      "delivered": "integer",
      "cancelled": "integer",
      "total_revenue": "decimal"
    },
    "payments": {
      "total": "integer",
      "pending": "integer",
      "paid": "integer",
      "failed": "integer",
      "refunded": "integer"
    }
  }
  ```
- Validasi:
    - start_date: optional, date
    - end_date: optional, date
- Authorization: Bearer token (admin only)
- Prioritas: Medium

**Task 12.2.2: Get Vendor Dashboard Stats**
- Endpoint: `GET /vendors/{id}/dashboard`
- HTTP Method: GET
- Path Params: id (integer)
- Query Params:
    - start_date (date, optional)
    - end_date (date, optional)
- Request Body: None
- Response Format (200):
  ```json
  {
    "products": {
      "total": "integer",
      "active": "integer",
      "out_of_stock": "integer"
    },
    "orders": {
      "total": "integer",
      "pending": "integer",
      "processing": "integer",
      "shipped": "integer",
      "delivered": "integer",
      "cancelled": "integer",
      "total_revenue": "decimal"
    },
    "transactions": {
      "total_sales": "decimal",
      "total_commissions": "decimal",
      "net_earnings": "decimal",
      "pending_withdrawal": "decimal"
    },
    "reviews": {
      "total": "integer",
      "pending": "integer",
      "average_rating": "decimal"
    }
  }
  ```
- Validasi:
    - id: required, must exist and belong to current user or admin
    - start_date: optional, date
    - end_date: optional, date
- Authorization: Bearer token (vendor owner, admin, or vendor users)
- Prioritas: Medium