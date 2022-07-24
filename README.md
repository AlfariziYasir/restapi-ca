## Architecture

project ini menggunakan desain clean architecture dengan 4 layer yang terdiri dari:

- Model
- Repository
- Service
- Handler

## Features

- Command line options
- JWT token create, refresh
- ORM gorm dengan database postgres dan redis untuk caching
- Dinamis pagination dengan sort, filter, search, dll
- yml config sebagai environment variabel
- Validasi pada request create, update, update password, captcha, dan login
- Middlewares cors, access control, logger, dll
