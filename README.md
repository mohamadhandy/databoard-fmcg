# Postgres

```sql
CREATE TABLE public."Admin" (
	id varchar(255) PRIMARY KEY,
	email varchar(255) NOT NULL,
	phone_number varchar(20) NOT NULL,
	"role" varchar(255) NOT NULL,
	status varchar(20) NOT NULL,
	"name" varchar(100) NOT NULL,
	created_at timestamp(0) NOT NULL,
	updated_at timestamp(0) NOT NULL
);

CREATE TABLE public."Brand" (
	id varchar(255) PRIMARY KEY,
	"name" varchar(100) NOT NULL,
	created_at timestamp(0) NOT NULL,
	updated_at timestamp(0) NOT null,
	created_by varchar(100) NOT NULL,
	updated_by varchar(100) NOT null,
	status varchar(20) NOT NULL
);

CREATE TABLE public."Category" (
	id varchar(255) PRIMARY KEY,
	"name" varchar(100) NOT NULL,
	created_at timestamp(0) NOT NULL,
	updated_at timestamp(0) NOT null,
	created_by varchar(100) NOT NULL,
	updated_by varchar(100) NOT null
);

CREATE TABLE public."Product" (
	id varchar(255) PRIMARY KEY,
	"name" varchar(255) NOT NULL,
	created_at timestamp(0) NOT NULL,
	updated_at timestamp(0) NOT null,
	created_by varchar(100) NOT NULL,
	updated_by varchar(100) NOT null,
	status varchar(20) NOT null,
	sku varchar(9) not null,
	brand_id varchar(255) not null,
	category_id varchar(255) not null,
	FOREIGN KEY (brand_id) REFERENCES "Brand"(id),
  FOREIGN KEY (category_id) REFERENCES "Category"(id)
);
```

# CREATE ADMIN
```sql
INSERT INTO "Admin" (id, name, email, phoneNumber, role, status, created_at, updated_at, password)
VALUES ('8d8a03da-aba4-473e-97dd-cfc8fecd542b', 'Admin Klikdaily', 'admin@example.com', '08121111111', 'superadmin', 'active', NOW(), NOW(), 'password');
```

# CREATE BRAND
```sql
INSERT INTO "Brand" (id, name, created_at, updated_at, created_by, updated_by)
VALUES ('001', 'Torabika', NOW(), NOW(), 'Admin Klikdaily', 'Admin Klikdaily');
```

# CREATE CATEGORY
```sql
INSERT INTO "Category" (id, name, created_at, updated_at, created_by, updated_by)
VALUES ('01', 'Minuman', NOW(), NOW(), 'Admin Klikdaily', 'Admin Klikdaily');
```

# select table query
```sql
select * from "Admin" a2 

select * from "Brand" b2

select * from "Category" c 

select * from "Product" p
```

# clone this project
1. git clone this project (branch sync-es-postgre-using-rabbitmq)
2. run with `fresh`

# docker preparation
1. rabbitmq 3.12-management
2. elasticsearch 7.14.0
3. redis:latest
4. docker postgres (if needed) mine are installed locally. 