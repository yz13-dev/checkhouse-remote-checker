# Checkhouse Remote Checker

Нужен для разворачивания на VPS и проверки мониторов из удаленных локаций.

---

Один из важных моментов, это использование переменной окружения `REGION` для указания региона, также region должен быть указан в теле токена.

Токен принимается в заголовке `Authorization` в формате `Bearer {token}`.

`REGION="{CONUTRY_CODE}:{CITY}"`

```.env
REGION="ru:moscow"
JWT_SECRET="..."
```

```bash
docker run --rm \
  -v $(pwd)/.env:/code/.env \
  -e GIN_MODE=release \
  -p 8080:8080 \
  api:latest
```
