# wallet_service

## использование

- устанавливаем модули
- запускаем рядом БД в докере
- запускаем сам сервис

```
make install-modules
make dbrun
make run
```

## дропнуть докер с БД

```
make stopdb
```

## API

### Создание кошелька

```
POST /wallet/new
name: str (required)
```

---

### Пополнение кошелька

```
POST /wallet/topup
wallet_id: uint32 (required)
amount: uint32 (required)
client_operation_hash: str (required)  // защита от задублирования операции
```

---

### Перевод денежных средств

```
POST /wallet/transfer
wallet_id_from: uint32 (required)
wallet_id_to: uint32 (required)
amount: uint32 (required)
client_operation_hash: uint32 (required)  // защита от задублирования операции
```

---

### Отчет о транзакциях

```
GET /wallet/report
wallet_id: uint32 (required)
date_from: int64 (required)
date_to: int64 (required)
type: int (0: получить все, 1: получить пополнения, 2: получить выводы)
```
