## 🚀 Fiber Docs

Fiber docs 👉 **https://docs.gofiber.io/**  
Using docs via command line (offline):
```sh
go doc github.com/gofiber/fiber/v2.Ctx.Query
```

---

## 🌐 URL Querying & Parameters

**Fiber equivalents:**

| Gin                          | Fiber                            |
|-----------------------------|----------------------------------|
| `c.Param("id")`             | `c.Params("id")`                 |
| `c.Query("id")`             | `c.Query("id")`                  |
| `c.DefaultQuery("id", "")`  | `c.Query("id", "defaultValue")`  |
| `c.PostForm("name")`        | `c.FormValue("name")`            |
| `c.DefaultPostForm(...)`    | `c.FormValue("name", "default")` |
| `c.GetHeader("Authorization")` | `c.Get("Authorization")`       |

---

### 📌 `c.Params("id")`
**URL parameter:** `/user/:id`

```go
app.Get("/user/:id", func(c *fiber.Ctx) error {
    id := c.Params("id") // e.g., john
    return c.SendString("User ID is: " + id)
})
```

---

### 🔎 `c.Query("name")` and default

```go
app.Get("/search", func(c *fiber.Ctx) error {
    name := c.Query("name", "unknown") // with default value
    return c.SendString("Hello, " + name)
})
```

---

### 📥 `c.FormValue(...)` (POST form fields)

```go
app.Post("/submit", func(c *fiber.Ctx) error {
    email := c.FormValue("email")
    return c.SendString("Submitted Email: " + email)
})
```

---

### 🧾 `c.Get("Authorization")` (Headers)

```go
app.Get("/auth", func(c *fiber.Ctx) error {
    token := c.Get("Authorization")
    return c.SendString("Token: " + token)
})
```

---

## 📦 Model Binding (JSON / URI / Headers)

Fiber uses `BodyParser()` for JSON and form binding.

| Gin Function           | Fiber Equivalent                    |
|------------------------|-------------------------------------|
| `c.ShouldBindJSON(&x)` | `c.BodyParser(&x)`                  |
| `c.ShouldBindUri(&x)`  | Use `c.Params(...)` manually or custom struct |
| `c.ShouldBindHeader()` | Use `c.Get(...)` manually           |

---

### 📌 JSON Binding

```go
type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

app.Post("/register", func(c *fiber.Ctx) error {
    user := new(User)

    if err := c.BodyParser(user); err != nil {
        return c.Status(400).SendString("Invalid payload")
    }

    return c.JSON(user)
})
```

---

### 🧪 Form & Header Binding

- For **form values** → use `c.FormValue("key")`
- For **headers** → use `c.Get("Header-Name")`

---

## 📎 Extra Tip: Custom Struct for URI (manual)
Unlike Gin, Fiber doesn’t support `ShouldBindUri` directly—you extract and manually bind:

```go
type Params struct {
    ID string
}

app.Get("/user/:id", func(c *fiber.Ctx) error {
    p := Params{
        ID: c.Params("id"),
    }
    return c.JSON(p)
})
```