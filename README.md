# api-starter

โครงตั้งต้น (starter) สำหรับทำ REST API ด้วย Go — ใช้ [Fiber v2](https://docs.gofiber.io/) เป็น web framework, [GORM](https://gorm.io/) + [GORM Gen](https://gorm.io/gen/) เป็น ORM และเชื่อมต่อ Microsoft SQL Server

มีโมดูลตัวอย่าง `user` ที่ทำครบทุกชั้น (router → handler → service → repository) ไว้ให้ copy ไปทำโมดูลใหม่

## สิ่งที่ต้องมี

- Go (เวอร์ชันล่าสุดที่ยัง support)
- Microsoft SQL Server
- [air](https://github.com/air-verse/air) สำหรับ live reload (ไม่บังคับ)

## เริ่มต้นใช้งาน

```bash
# 1. ตั้งค่า config
cp .env.template .env        # แล้วแก้ค่า DB_* ให้ตรงกับเครื่อง

# 2. ติดตั้ง dependencies
go mod tidy

# 3. สร้างตาราง + generate query code (ต้องต่อ DB ได้)
go run ./gen

# 4. รัน server
go run .                     # หรือ `air` ถ้าต้องการ live reload
```

เปิด <http://localhost:3005/api/v1> จะได้ข้อมูลของ service กลับมา

## Config

อ่านค่าจาก environment variable โดยโหลดไฟล์ `/configs/.env` ก่อนถ้ามี (ไว้ mount ตอน deploy) ถ้าไม่มีจะใช้ `.env` ในโฟลเดอร์ที่รัน ค่าที่ตั้งไว้ใน environment อยู่แล้วจะไม่ถูกไฟล์ `.env` ทับ

| ตัวแปร | ความหมาย | ตัวอย่าง |
| --- | --- | --- |
| `PRODUCT_CODE` | รหัส product | `starter-api` |
| `MODULE_NAME` | ชื่อ service (ใช้เป็นชื่อ Fiber app) | `starter-api` |
| `BUILD` | ชื่อ environment — ถ้าเป็น `DEV` จะ log SQL ทุกคำสั่ง | `DEV` |
| `RELEASE` | เวอร์ชัน | `1.0.0` |
| `PORT` | port ของ HTTP server (ค่าเริ่มต้น `3005`) | `3005` |
| `DB_SERVER` | host ของ SQL Server | `localhost` |
| `DB_PORT` | port ของ SQL Server | `1433` |
| `DB_NAME` | ชื่อ database | |
| `DB_USER` | user | |
| `DB_PASSWORD` | password | |
| `DB_FED_AUTH` | โหมด Azure AD auth (`fedauth`) เว้นว่างได้ถ้าใช้ SQL Server authentication | |

เพิ่ม config ใหม่ได้ที่ struct `Config` ใน `pkg/env/env.go` โดยใส่ tag `env:"ชื่อตัวแปร"`

## โครงสร้างโปรเจกต์

```
.
├── main.go                 # จุดเริ่มโปรแกรม: โหลด config, ต่อ DB, สร้าง Fiber app, graceful shutdown
├── external/
│   └── routers.go          # รวม route ทั้งหมด + error handler กลาง (ตอบเป็น JSON)
├── internal/               # โมดูลตามฟีเจอร์ 1 โฟลเดอร์ต่อ 1 โมดูล
│   └── user/
│       ├── router.go       # ประกอบ repository → service → handler และผูก route
│       ├── handler.go      # ชั้น HTTP: อ่าน request, แปลง error เป็น status code
│       ├── service.go      # business logic + validation
│       └── repository.go   # คุยกับ database ผ่าน GORM
├── domain/
│   ├── entity/             # struct ของตาราง (ใช้ migrate และ generate query)
│   └── model/              # รูปร่างของ request / response
├── gen/
│   └── db_migrate.go       # คำสั่ง migrate schema + generate query code
├── pkg/                    # โค้ดกลางที่ไม่ผูกกับฟีเจอร์
│   ├── db/                 # เชื่อมต่อ SQL Server
│   ├── env/                # โหลด config
│   └── util/               # ค่าคงที่ที่ใช้ร่วมกัน
├── .env.template           # ตัวอย่างไฟล์ config
└── .air.toml               # config ของ air
```

ทิศทางการเรียก: `handler → service → repository → database` แต่ละชั้นรู้จักชั้นถัดไปผ่าน interface เท่านั้น (`Service`, `Repository`) จึง mock เพื่อเขียน test ได้

## API

Base path: `/api/v1`

| Method | Path | ทำอะไร |
| --- | --- | --- |
| GET | `/api/v1` | ข้อมูลของ service (product code, build, release) |
| POST | `/api/v1/users` | สร้าง user |
| GET | `/api/v1/users` | รายการ user ทั้งหมด |
| GET | `/api/v1/users/:id` | ดู user ตาม id |

ตัวอย่างสร้าง user:

```bash
curl -X POST http://localhost:3005/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"firstName":"Somchai","lastName":"Jaidee","email":"somchai@example.com","password":"12345678"}'
```

- password ถูกเก็บเป็น bcrypt hash และไม่ถูกส่งกลับใน response
- `roleId` ไม่บังคับ ถ้าส่งมาต้องเป็น id ที่มีอยู่ในตาราง `user_roles`

เมื่อเกิด error จะตอบเป็น JSON รูปแบบเดียวกันทั้งระบบ:

```json
{ "message": "email is already in use" }
```

| Status | กรณี |
| --- | --- |
| 400 | body ผิดรูปแบบ หรือ validate ไม่ผ่าน |
| 404 | ไม่พบข้อมูล |
| 409 | email ซ้ำ |
| 500 | error อื่น ๆ (รายละเอียดอยู่ใน log ไม่ส่งให้ client) |

## Database

Schema มาจาก struct ใน `domain/entity` คำสั่งนี้จะ AutoMigrate ตาราง แล้ว generate query code แบบ type-safe ของ GORM Gen ไว้ที่ `./query`:

```bash
go run ./gen      # ต้องรันจาก root ของโปรเจกต์
```

รันใหม่ทุกครั้งที่แก้ entity ถ้าเพิ่ม entity ใหม่ให้เพิ่มทั้งใน `AutoMigrate` และ `ApplyBasic` ที่ `gen/db_migrate.go`

> column ชนิด string ที่มี index ต้องกำหนด `size` เสมอ (เช่น `gorm:"size:255;uniqueIndex"`) เพราะ SQL Server สร้าง index บน `nvarchar(MAX)` ไม่ได้

repository ของโมดูล `user` ตอนนี้เรียก GORM ตรง ๆ จึง build ได้โดยไม่ต้อง generate ก่อน ถ้าต้องการใช้ query ที่ generate แล้วให้เรียก `query.SetDefault(db)` ตอนเริ่มโปรแกรม

## เพิ่มโมดูลใหม่

1. เพิ่ม entity ใน `domain/entity` และ request/response ใน `domain/model`
2. เพิ่ม entity ใน `gen/db_migrate.go` แล้วรัน `go run ./gen`
3. สร้าง `internal/<module>/` โดย copy 4 ไฟล์จาก `internal/user` มาแก้
4. ลงทะเบียน router ใน `external/routers.go`:

   ```go
   user.Router(v1, db)
   product.Router(v1, db) // โมดูลใหม่
   ```

## Swagger

route `/api-doc/*` ผูกกับ Swagger UI ไว้แล้ว แต่ยังไม่ได้ generate เอกสาร ถ้าจะใช้งาน:

1. ใส่ annotation ของ [swag](https://github.com/swaggo/swag) ที่ `main.go` และ handler
2. รัน `swag init` เพื่อสร้างโฟลเดอร์ `docs/`
3. import `_ "api-starter/docs"` ใน `main.go`
