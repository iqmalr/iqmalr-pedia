# Daftar Skenario Test Postman untuk go-auth/v2

Base URL: http://localhost:8082
Prefix API: /api/v2

Catatan umum:
- Endpoint yang memerlukan autentikasi menggunakan header: Authorization: Bearer <access_token>
- Endpoint refresh token menggunakan header Authorization: Bearer <refresh_token>
- Role-based access:
  - /admin/...: hanya admin
  - /vendor/...: vendor dan admin

Lingkup endpoint (berdasarkan kode v2):
- GET /health
- POST /api/v2/auth/register
- POST /api/v2/auth/login
- POST /api/v2/auth/logout (butuh access token)
- POST /api/v2/auth/refresh (butuh refresh token di header Authorization)
- POST /api/v2/auth/forgot-password
- POST /api/v2/auth/reset-password
- POST /api/v2/auth/verify-email
- POST /api/v2/auth/resend-verification
- GET /api/v2/auth/me (butuh access token)
- GET /api/v2/admin/dashboard (butuh access token + role admin)
- GET /api/v2/vendor/dashboard (butuh access token + role vendor/admin)

## 1) GET /health
Skenario:
- 200 OK: Mendapatkan status dan version.
  - Tests:
    - pm.test("status 200", () => pm.response.code === 200)
    - pm.test("punya field status", () => pm.response.json().status === "OK")

## 2) POST /api/v2/auth/register
Body (contoh sukses minimal):
{
  "name": "John Doe",
  "email": "john.doe+1@example.com",
  "phone": "+6281234567890",
  "password": "StrongPass123",
  "role": "customer" // opsional: customer | vendor | admin
}
Skenario:
- 201 Created (sukses daftar user baru)
- 400 Bad Request: email sudah ada
- 400: phone sudah ada (jika phone diisi dan sudah dipakai)
- 400: format email tidak valid
- 400: password < 8 karakter
- 400: name kurang dari 3 karakter
- 400: role tidak termasuk [customer, vendor, admin]
Tests umum sukses:
- status 201
- response.json().message == "User registered successfully"
- response.json().user.id ada
- response.json().user.email == input email

## 3) POST /api/v2/auth/login
Body:
{
  "email": "john.doe+1@example.com",
  "password": "StrongPass123"
}
Skenario:
- 200 OK: kredensial benar, dapat access token + refresh token tersimpan di server
- 401 Unauthorized: password salah
- 401: email tidak terdaftar (disamarkan sebagai invalid credentials)
- 400: format email tidak valid
- 400: field wajib kosong
Tests sukses:
- status 200
- response.json().message == "Login successful"
- response.json().token (access token) ada
- response.json().user.role sesuai
Catatan: Refresh token tidak dikembalikan di response, tetapi disimpan di DB (field RefreshToken). Digunakan pada endpoint /refresh melalui header.

## 4) POST /api/v2/auth/logout
Header: Authorization: Bearer <access_token>
Skenario:
- 200 OK: token akses valid, refresh token di DB dibersihkan
- 401: tanpa Authorization header
- 401: token tidak valid / expired
- 401: salah kirim refresh token sebagai access token (harusnya tetap 401 dari middleware)
Tests sukses:
- status 200
- response.json().message == "Logout successful"

## 5) POST /api/v2/auth/refresh
Header: Authorization: Bearer <refresh_token>
Skenario:
- 200 OK: refresh token valid dan cocok dengan yang tersimpan di DB -> dapat access token baru
- 401: tanpa Authorization header
- 401: token tidak valid/expired
- 401: refresh token tidak sama dengan yang tersimpan (misal setelah logout)
- 401: pakai access token biasa sebagai refresh (verifikasi kemungkinan gagal jika claim/masa berlaku berbeda)
Tests sukses:
- status 200
- response.json().token ada

## 6) POST /api/v2/auth/forgot-password
Body:
{
  "email": "john.doe+1@example.com"
}
Skenario:
- 200 OK: email terdaftar -> token reset dibuat (link email TODO, tidak benar-benar mengirim email)
- 200 OK: email tidak terdaftar (disamarkan untuk keamanan)
- 400: format email tidak valid
Tests sukses:
- status 200
- response.json().message == "Password reset link sent to your email"
Catatan: Ambil token reset password secara manual dari DB untuk pengujian reset-password.

## 7) POST /api/v2/auth/reset-password
Body:
{
  "token": "<token_reset_dari_DB>",
  "password": "NewStrongPass123",
  "password_confirmation": "NewStrongPass123"
}
Skenario:
- 200 OK: token valid, password diperbarui, token ditandai used=true
- 400: token invalid/expired
- 400: password < 8 karakter
- 400: password_confirmation tidak sama
Tests sukses:
- status 200
- response.json().message == "Password reset successful"

## 8) POST /api/v2/auth/verify-email
Body:
{
  "token": "<token_verifikasi_dari_DB>"
}
Skenario:
- 200 OK: token valid -> email_verified_at terisi & token dihapus
- 400: token tidak valid
Tests sukses:
- status 200
- response.json().message == "Email verified successfully"

## 9) POST /api/v2/auth/resend-verification
Body:
{
  "email": "john.doe+1@example.com"
}
Skenario:
- 200 OK: user belum terverifikasi -> token baru dibuat
- 400: email sudah terverifikasi ("email already verified")
- 200 OK: email tidak terdaftar (disamarkan)
- 400: format email tidak valid
Tests sukses:
- status 200
- response.json().message == "Verification email sent"

## 10) GET /api/v2/auth/me
Header: Authorization: Bearer <access_token>
Skenario:
- 200 OK: token valid -> mendapatkan profil user
- 401: tanpa header Authorization
- 401: token invalid/expired
Tests sukses:
- status 200
- response.json().email ada dan cocok

## 11) GET /api/v2/admin/dashboard
Header: Authorization: Bearer <access_token_dengan_role_admin>
Skenario:
- 200 OK: role admin -> akses berhasil
- 403 Forbidden: role customer/vendor
- 401: tanpa token / token invalid
Tests sukses:
- status 200
- response.json().message == "Welcome to admin dashboard"

## 12) GET /api/v2/vendor/dashboard
Header: Authorization: Bearer <access_token_dengan_role_vendor_atau_admin>
Skenario:
- 200 OK: role vendor
- 200 OK: role admin
- 403 Forbidden: role customer
- 401: tanpa token / token invalid
Tests sukses:
- status 200
- response.json().message == "Welcome to vendor dashboard"

---

## Urutan Alur Uji yang Direkomendasikan di Postman
1. GET /health (cek server up)
2. POST /auth/register (buat 3 user: customer, vendor, admin) dengan email unik
3. POST /auth/login untuk masing-masing user; simpan access_token (env var: access_token_customer, access_token_vendor, access_token_admin) dan juga lakukan query DB untuk melihat refresh_token tersimpan (opsional)
4. GET /auth/me dengan access_token_customer (cek profil)
5. GET /vendor/dashboard dengan access_token_vendor (200), dengan access_token_admin (200), dengan access_token_customer (403)
6. GET /admin/dashboard dengan access_token_admin (200), dengan access_token_vendor (403), dengan access_token_customer (403)
7. POST /auth/refresh dengan refresh_token milik user (set di Authorization header) -> simpan token baru sebagai access_token_baru
8. POST /auth/logout dengan access_token_customer -> kemudian coba POST /auth/refresh dengan refresh_token lama -> harus 401
9. POST /auth/forgot-password untuk user customer -> ambil token reset dari DB, lalu POST /auth/reset-password
10. POST /auth/resend-verification untuk user baru yang belum verifikasi -> lalu POST /auth/verify-email dengan token dari DB

## Contoh Tests (snippet) untuk Postman
- Verifikasi status dan struktur dasar:
  pm.test("status 200", function () { pm.response.to.have.status(200); });
  const json = pm.response.json();
  pm.expect(json).to.have.property("message");

- Simpan token ke environment (setelah login):
  const data = pm.response.json();
  pm.environment.set("access_token_customer", data.token);

- Gunakan Authorization otomatis (Pre-request Script di folder/collection):
  const tok = pm.environment.get("access_token_customer");
  if (tok) { pm.request.headers.add({ key: "Authorization", value: `Bearer ${tok}` }); }

## Catatan Integrasi/DB
- Aplikasi menggunakan PostgreSQL dengan tabel user dan token (PasswordResetToken, EmailVerificationToken). Token reset/verifikasi tidak dikirim via email; ambil manual dari DB untuk pengujian end-to-end.
- Refresh token disimpan di kolom RefreshToken user setelah login. Endpoint /auth/refresh memvalidasi bahwa refresh token pada header sama dengan yang tersimpan di DB.

## Variabel Environment Postman yang Disarankan
- baseUrl: http://localhost:8082
- apiPrefix: /api/v2
- access_token_customer
- access_token_vendor
- access_token_admin
- refresh_token_<role> (opsional jika Anda expose via query DB/tools)
