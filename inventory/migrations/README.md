# 🗺️ Руководство по написанию миграций MongoDB

## 📋 Формат миграций

**Все миграции** — это **массив JSON-объектов**, где каждый объект описывает **одну операцию**.

```json
[
  { "create": "collection_name" },
  { "createIndexes": "collection_name", "indexes": [...] },
  { "insert": "collection_name", "documents": [...] }
]
```

### Структура файлов:
```
migrations/
├── 0001_create_parts.up.json     ← Создание
├── 0001_create_parts.down.json   ← Откат
├── 0002_add_users.up.json
└── 0002_add_users.down.json
```

**Важно:** Для каждой миграции **обязательно** нужны **`.up.json`** и **`.down.json`**.

---

## 🛠️ Основные команды

### 1. **`create`** — Создание коллекции

```json
{ "create": "parts" }
```

**`.up.json`:**
```json
[ { "create": "parts" } ]
```

**`.down.json`:**
```json
[ { "drop": "parts" } ]
```

---

### 2. **`createIndexes`** — Создание индексов

```json
{
  "createIndexes": "parts",
  "indexes": [
    {
      "key": { "part_uuid": 1 },
      "name": "idx_part_uuid_unique",
      "unique": true
    },
    {
      "key": { "status": 1, "created_at": -1 },
      "name": "idx_status_created_desc"
    }
  ]
}
```

**`.up.json`:**
```json
[
  { "create": "parts" },
  {
    "createIndexes": "parts",
    "indexes": [
      {
        "key": { "part_uuid": 1 },
        "name": "idx_part_uuid_unique",
        "unique": true
      }
    ]
  }
]
```

**`.down.json`:**
```json
[
  {
    "dropIndexes": "parts",
    "indexNames": ["idx_part_uuid_unique"]
  },
  { "drop": "parts" }
]
```

---

### 3. **`insert`** — Вставка документов

```json
{
  "insert": "parts",
  "documents": [
    {
      "part_uuid": "part-001",
      "name": "Fuel Tank",
      "status": "active",
      "created_at": { "$date": "2025-01-01T00:00:00Z" }
    }
  ]
}
```

**`.up.json`:**
```json
[
  { "create": "parts" },
  {
    "insert": "parts",
    "documents": [
      {
        "part_uuid": "part-001",
        "name": "Fuel Tank",
        "status": "active",
        "created_at": { "$date": "2025-01-01T00:00:00Z" }
      }
    ]
  }
]
```

**`.down.json`:**
```json
[ { "drop": "parts" } ]
```

---

### 4. **`drop`** — Удаление коллекции

```json
{ "drop": "parts" }
```

---

### 5. **`dropIndexes`** — Удаление индексов

```json
{
  "dropIndexes": "parts",
  "indexNames": ["idx_part_uuid_unique", "idx_status_created_desc"]
}
```

**Удалить ВСЕ индексы (кроме `_id_`):**
```json
{ "dropIndexes": "parts", "indexNames": ["*"] }
```

---

### 6. **`renameCollection`** — Переименование коллекции

```json
{ "renameCollection": "old_name", "to": "new_name" }
```

**`.up.json`:**
```json
[ { "renameCollection": "parts_v1", "to": "parts" } ]
```

**`.down.json`:**
```json
[ { "renameCollection": "parts", "to": "parts_v1" } ]
```

---

## 📊 Типы индексов

| Тип | Пример |
|-----|--------|
| **Один поле** | `{ "key": { "status": 1 } }` |
| **Несколько полей** | `{ "key": { "status": 1, "created_at": -1 } }` |
| **Уникальный** | `{ "key": { "part_uuid": 1 }, "unique": true }` |
| **TTL (автоудаление)** | `{ "key": { "expires_at": 1 }, "expireAfterSeconds": 3600 }` |
| **Text** | `{ "key": { "$**": "text" }, "name": "full_text" }` |

---

## 🧪 Полные примеры миграций

### Миграция 1: Создание коллекции `parts`

**`0001_create_parts.up.json`:**
```json
[
  { "create": "parts" },
  {
    "createIndexes": "parts",
    "indexes": [
      {
        "key": { "part_uuid": 1 },
        "name": "idx_part_uuid_unique",
        "unique": true
      },
      {
        "key": { "status": 1, "created_at": -1 },
        "name": "idx_status_created_desc"
      }
    ]
  },
  {
    "insert": "parts",
    "documents": [
      {
        "part_uuid": "part-001",
        "name": "Fuel Tank",
        "status": "active",
        "created_at": { "$date": "2025-01-01T00:00:00Z" }
      }
    ]
  }
]
```

**`0001_create_parts.down.json`:**
```json
[ { "drop": "parts" } ]
```

---

### Миграция 2: Добавление коллекции `users`

**`0002_create_users.up.json`:**
```json
[
  { "create": "users" },
  {
    "createIndexes": "users",
    "indexes": [
      {
        "key": { "email": 1 },
        "name": "idx_email_unique",
        "unique": true
      },
      {
        "key": { "user_uuid": 1 },
        "name": "idx_user_uuid_unique",
        "unique": true
      }
    ]
  }
]
```

**`0002_create_users.down.json`:**
```json
[ { "drop": "users" } ]
```

---

## ⚠️ Corner Cases и подводные камни

### 1. **Даты в MongoDB**
```json
"created_at": { "$date": "2025-01-01T00:00:00Z" }  // ✅ Правильно
"created_at": "2025-01-01T00:00:00Z"               // ❌ Строкой НЕ работает!
```

### 2. **ObjectId**
```json
"_id": { "$oid": "507f1f77bcf86cd799439011" }     // ✅ Правильно
```

### 3. **Массивы в документах**
```json
{
  "tags": ["rocket", "engine", "v1"],
  "metadata": { "version": "1.0", "active": true }
}
```

### 4. **"Dirty" миграции**
```
error: Dirty database version 1. Fix and force version.
```

**Решение:**
```bash
migrate -path ./migrations -database "mongodb://..." force 1
```

### 5. **Порядок операций критичен!**
```json
[
  { "create": "parts" },                    // 1️⃣ Сначала создать
  { "createIndexes": "parts", ... },        // 2️⃣ Потом индексы
  { "insert": "parts", ... }               // 3️⃣ Потом данные
]
```

**❌ Неправильно:**
```json
[
  { "insert": "parts", ... },              // Ошибка! Коллекция не существует
  { "create": "parts" }
]
```

### 6. **Имена файлов — строго по маске**
```
0001_xxx.up.json
0002_yyy.up.json
0001_xxx.down.json  ← Тот же номер!
```

**❌ НЕ работает:**
```
1_create.up.json      // Ноль слева обязателен
001_create.up.json    // Только цифры + подчёркивание
```

---

## 🚀 Запуск миграций

```bash
# Применить все
migrate -path ./migrations -database "mongodb://..." up

# Откатить на 1 шаг
migrate -path ./migrations -database "mongodb://..." down 1

# Проверить статус
migrate -path ./migrations -database "mongodb://..." version

# Принудительно исправить dirty
migrate -path ./migrations -database "mongodb://..." force 1
```

---

## 📝 Шаблон для новой миграции

```json
// NNNN_description.up.json
[
  { "create": "collection_name" },
  {
    "createIndexes": "collection_name",
    "indexes": [
      {
        "key": { "field_name": 1 },
        "name": "idx_field_name",
        "unique": true
      }
    ]
  }
]

// NNNN_description.down.json
[
  { "drop": "collection_name" }
]
```

---

## ✅ Checklist перед коммитом

- [ ] Есть `.up.json` и `.down.json`
- [ ] Номера миграций последовательные
- [ ] `.down.json` полностью откатывает `.up.json`
- [ ] Проверены даты (`$date`)
- [ ] Имена индексов уникальные (`idx_` + описание)
- [ ] Тестировал локально: `up` → `down` → `up`
