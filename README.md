# sagasql

## Environment variable
Aplikasi membutuhkan variable yang wajib diisi untuk koneksi ke database dan keamanan JWT
.  
copy atau rename file `.env-example` menjadi `.env`. Sesuaikan isinya seperti pada contoh
file ini akan diload ketika program dijalankan.

## Database
Aplikasi memerlukan database `PostgreSQL` dengan nama database `testsaga`.  
adapun Table yang dibutuhkan sebagai berikut : (tidak sempat dibuat automigration)

        
        CREATE TABLE users (
        username VARCHAR ( 50 ) PRIMARY KEY,
        name VARCHAR (100) NOT NULL,
        email VARCHAR ( 255 ) UNIQUE NOT NULL,
        password VARCHAR (100) NOT NULL,
        role VARCHAR(50) NOT NULL,
        created_at BIGINT NOT NULL,
        updated_at BIGINT NOT NULL
        )
        
        
        CREATE TABLE products (
        product_id SERIAL PRIMARY KEY,
        name VARCHAR (100) UNIQUE NOT NULL,
        price BIGINT,
        image VARCHAR(50),
        created_by VARCHAR REFERENCES users(username),
        created_at BIGINT NOT NULL
        )

design database dibuat se flat mungkin hanya untuk mempermudah test ini.

## Endpoint

postman config disertakan pada folder postman_file. sedikit sample akan ditulis dibawah ini :

### USER
1. `POST` `{{url}}/api/v1/register-force` meregister user tanpa auth admin    
   Body :
```json
{
  "username":"muchlis",
  "email": "whois.who@gmail.com",
  "name": "muchlis",
  "password": "Password",
  "role": "ADMIN"
}
```
role hanya dapat di isi dengan ADMIN atau NORMAL

2. `POST` `{{url}}/api/v1/login`  login dengan mengembalikan access token dan refresh token. token ini nantinya yang akan terus
dilampirkan pada setiap request pada `header` key : `Authorization` dan value `Bearer {token_tanpa_curly_brace}`
Body : 
```json
{
  "username":"muchlis",
  "password":"Password"
}
```

### Product
1. `GET` `{{url}}/api/v1/products` menampilkan list produk, memerlukan token
2. `POST` `{{url}}/api/v1/products` menambahkan products  
   Body :
```json
{
  "name": "Mangga",
  "price": 50000
}
```
3. `POST` `{{url}}/api/v1/products-image/:id` mengupload gambar product.  
gunakan form-data dengan key "image" dan value {gambarnya}.


### Daftar lengkap map url
```
	/*
	// url mapping
	api := app.Group("/api/v1")

	//USER
	api.Get("/users/:username", userHandler.Get)
	api.Get("/users", userHandler.Find)
	api.Post("/login", userHandler.Login)
	api.Post("/refresh", userHandler.RefreshToken)
	api.Get("/profile", middle.NormalAuth(), userHandler.GetProfile)
	api.Post("/register-force", userHandler.Register)                                // <- seharusnya gunakan middleware agar hanya admin yang bisa meregistrasi
	api.Post("/register", middle.NormalAuth(config.RoleAdmin), userHandler.Register) // <- hanya admin yang bisa meregistrasi
	api.Put("/users/:username", middle.NormalAuth(config.RoleAdmin), userHandler.Edit)
	api.Delete("/users/:username", middle.NormalAuth(config.RoleAdmin), userHandler.Delete)

	//PRODUCT
	api.Get("/products/:id", middle.NormalAuth(), productHandler.Get)
	api.Get("/products", middle.NormalAuth(), productHandler.Find)
	api.Post("/products", middle.NormalAuth(), productHandler.Insert)
	api.Put("/products/:id", middle.NormalAuth(), productHandler.Edit)
	api.Delete("/products/:id", middle.NormalAuth(), productHandler.Delete)
	api.Post("/products-image/:id", middle.NormalAuth(), productHandler.UploadImage) // <- upload image multipath*/
```


