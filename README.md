# exercise-login-back-go

## Configuración de la Base de Datos

Este proyecto utiliza SQL Server como sistema de gestión de base de datos. A continuación, se encuentran los scripts utilizados para crear la estructura necesaria en la base de datos.

### Tablas

Aquí se incluye el script para crear la tabla de usuarios necesaria para el funcionamiento de la API:

````sql
CREATE TABLE users (
    id INT IDENTITY(1,1) PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(10) NOT NULL,
    password VARCHAR(255) NOT NULL,
    CONSTRAINT UC_users_email UNIQUE (email),
    CONSTRAINT UC_users_phone UNIQUE (phone)
);
````

## Procedimientos Almacenados (Stored Procedures)

Los siguientes procedimientos almacenados se utilizan para la manipulación de datos de usuarios:

### GetUserByEmailOrUsername

Obtiene un usuario por su correo electrónico o nombre de usuario:

```sql
CREATE PROCEDURE GetUserByEmailOrUsername
    @EmailOrUsername VARCHAR(255)
AS
BEGIN
    SELECT * FROM Users
    WHERE Email = @EmailOrUsername OR Username = @EmailOrUsername
END;
```


### GetUserByEmailOrPhone
Obtiene un usuario por su correo electrónico o número de teléfono:

```sql
CREATE PROCEDURE GetUserByEmailOrPhone
    @Email VARCHAR(255),
    @Phone VARCHAR(10)
AS
BEGIN
    SELECT * FROM users
    WHERE email = @Email OR phone = @Phone
END
```

### CreateUser
Crea un nuevo usuario:

```sql
CREATE PROCEDURE CreateUser
    @Username VARCHAR(255),
    @Email VARCHAR(255),
    @Phone VARCHAR(10),
    @Password VARCHAR(255)
AS
BEGIN
    INSERT INTO users (username, email, phone, password)
    VALUES (@Username, @Email, @Phone, @Password)
END
```
