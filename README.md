[![Coverage](https://codecov.io/gh/hikjik/gophkeeper/branch/dev/graph/badge.svg?token=7XARTW7JX9)](https://codecov.io/gh/hikjik/gophkeeper)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://tldrlegal.com/license/gnu-lesser-general-public-license-v3-(lgpl-3))
[![Lint tests](https://github.com/hikjik/gophkeeper/actions/workflows/lint-tests.yml/badge.svg)](https://github.com/hikjik/gophkeeper/actions/workflows/lint-tests.yml)
[![Unit tests](https://github.com/hikjik/gophkeeper/actions/workflows/unit-tests.yml/badge.svg)](https://github.com/hikjik/gophkeeper/actions/workflows/unit-tests.yml)
------
# Менеджер паролей GophKeeper

GophKeeper представляет собой клиент-серверную систему, позволяющую пользователю безопасно хранить
логины, пароли, данные банковских карт, произвольные текстовые и бинарные данные.

Сервер поддерживает следующий функционал:
 * регистрация, аутентификация и авторизация пользователей;
 * хранение приватных данных пользователей;
 * синхронизация данных между несколькими авторизованными клиентами одного владельца;
 * передача приватных данных владельцу по запросу.

Клиент реализует следующую бизнес-логику:
 * регистрация, аутентификация и авторизация пользователей на удалённом сервере;
 * доступ к приватным данным по запросу.

## Настройка и запуск сервера

Перед запуском сервера необходимо создать конфигурационный файл с настройками
или задать переменные окружения. Пример конфигурационного файла:

```
# server-config.yaml

grpc:
  address: 127.0.0.1:9090
db:
  url: postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable
auth:
  key: xiuw1bi4r98vd1(&*6
  expiration_time: 24h
hasher:
  key: jc7YSHpH287)(*2bSq
```

Пример настройки сервера через переменные окружения:

```
# .env

GRPC_ADDRESS=127.0.0.1:9090
DB_URL=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable
AUTH_KEY=xiuw1bi4r98vd1(&*6
AUTH_EXPIRATION_TIME=24h
HASHER_KEY=jc7YSHpH287)(*2bSq
```

После настройки запуск сервера осуществляется командной:

```
./gophkeeper server
```

## Настройка и запуск клиента

Перед запуском клиента необходимо создать конфигурационный файл с настройками
или задать переменные окружения. Пример конфигурационного файла:

```
# client-config.yaml

grpc:
  address: 127.0.0.1:9090
encryption:
  key: WYJcWgkItShq513L21E1CFuz6uQWDy3p
```

Пример настройки клиента через переменные окружения:

```
# .env

GRPC_ADDRESS=127.0.0.1:9090
ENCRYPTION_KEY=WYJcWgkItShq513L21E1CFuz6uQWDy3p
```

## Процедуры регистрации, аутентификации, авторизации

При регистрации пользователя необходимо указать адрес электронной почты и пароль.
Пример команды регистрации:

```
./gophkeeper-cli auth register -e user@mail.ru -p 123456
```

В случае успешного выполнения запроса регистрации нового пользователя,
сервер вернет в ответ токен доступа.
Токен будет сохранен в текстовый файл.

В случае необходимости, токен доступа можно запросить повторно с помощью команды:

```
./gophkeeper-cli auth login -e user@mail.ru -p 123456
```

## Хранение приватных данных пользователя

После записи полученного при регистрации токена доступа в переменную окружения TOKEN
становятся доступны команды для управления приватными данными пользователя.

### Добавление приватных данных

1. Пример команды, сохраняющей данные банковской карты:

```
./gophkeeper-cli secret create card \
  --name visa \
  --number 1111222233334444 \
  --date 12/22 \
  --holder Alexandr \
  --code 512
```

2. Пример команды создания пары логин пароль:

```
./gophkeeper-cli secret create credentials \
  --name yandex-mail \
  --login user@yandex.ru \
  --password 12345678
```

3. Пример команды для сохранения текстовой информации

```
./gophkeeper-cli secret create text \
  --name pushkin \
  --data "medny vsadnik"
```

3. Пример команды для сохранения бинарных данных

```
./gophkeeper-cli secret create bin \
  --name code \
  -f main.go
```

### Получение данных

Для получения приватных данных необходимо указать название секрета, пример:

```
./gophkeeper-cli secret get --name visa
```

Также можно вывести список всех приватных данных пользователя:

```
./gophkeeper-cli secret list
```

### Редактирование и удаление данных

Пример редактирования данных о банковской карте:

```
./gophkeeper-cli secret update card \
  --name visa \
  --number 1111222233335555 \
  --date 12/30 \
  --holder "Alexandr Pushkin" \
  --code 256
```

Пример команды удаления данных:

```
./gophkeeper-cli secret delete --name visa
```
