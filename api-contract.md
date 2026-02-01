# API Contract

## Authorization

```
name: login
method: POST
path: /api/auth/v1/login
description:
    Melakukan autentikasi user menggunakan email dan password.
    Jika berhasil, sistem akan mengembalikan access token dan refresh token.
auth: none
request:
    body:
    {
        "username": string,
        "email": string,
        "password": string
    }
response:
    success (code: 200):
    {
        "message": string,
        "code": string,
        "result": {
            "user": {
                "id": string,
                "name": string,
                "email": string
            },
            "access_token": string,
            "refresh_token": string
        }
    }

    error (code: 400, 401, 500):
    {
        "message": string,
        "code": int,
        "error": string
    }

```

```
name: register
method: POST
path: /api/auth/v1/register
description:
    Mendaftarkan user baru ke dalam sistem EduTask.
    Endpoint ini membuat akun user dan menghasilkan access token serta refresh token.
auth: none
request:
    body:
    {
        "name": string,
        "email": string,
        "password": string
    }
response:
    success (code: 201):
    {
        "message": string,
        "code": int,
        "result": {
            "user": {
                "id": string,
                "name": string,
                "email": string,
                "created_at": string
            },
            "access_token": string,
            "refresh_token": string
        }
    }

    error (code: 400, 500):
    {
        "message": string,
        "code": int,
        "error": string
    }
```

```
name: forget-password
method: POST
path: /api/auth/v1/forget-password
description:
    Melupakan password user dengan validasi via code yang di kirim ke email user.
    Endpoint ini digunakan ketika user lupa password akun nya.
    Client hit api ini, lalu service kirim code, dan tunggu redirect dari user atau client.
auth: none
request:
    body:
    {
        "email": string
    }
response:
    success (code: 201):
    {
        "message": string,
        "code": int,
        "result": {
            "token": string,
            "expired": int
        }
    }

    error (code: 400, 500):
    {
        "message": string,
        "code": int,
        "error": string
    }
```

```
name: refresh-token
method: POST
path: /api/auth/v1/refresh-token
description:
    Menghasilkan access token baru menggunakan refresh token yang masih valid.
    Digunakan ketika access token telah kedaluwarsa.
auth: none
request:
    body:
    {
        "refresh_token": string
    }
response:
    success (code: 200):
    {
        "message": string,
        "code": int,
        "result": {
            "access_token": string
        }
    }

    error (code: 401, 400, 500):
    {
        "message": string,
        "code": int,
        "error": string
    }

```

```
name: logout
method: POST
path: /api/auth/v1/logout
description:
    Mengakhiri sesi login user dengan cara menginvalidasi refresh token.
    Access token akan otomatis kedaluwarsa sesuai TTL.
    Dan update status is_active dalam database menjadi false.
auth: bearer
response:
    success (code: 200):
    {
        "message": string,
        "code": int
    }

    error (code: 401, 500):
    {
        "message": string,
        "code": int,
        "error": string
    }

```

## User

```
name: get-me
method: GET
path: /api/users/v1/me
description:
    Mengambil detail data user berdasarkan user_id di dalam jwt token.
    Endpoint ini digunakan untuk kebutuhan internal sistem
    seperti penampilan anggota kelas, group, atau submission.
auth: bearer
response:
    success (code: 200):
    {
        "message": string,
        "code": int,
        "result": {
            "id": string,
            "name": string,
            "email": string,
            "photo_url": string,
            "created_at": string
        }
    }

    error (code: 401, 404, 500):
    {
        "message": string,
        "code": int,
        "error": string
    }

```

```
name: change-name
method: PATCH
path: /api/users/v1/change-name
description:
    Memperbarui username yang sedang login.
    User hanya dapat memperbarui data miliknya sendiri.
auth: bearer
request:
    body:
    {
        "name": string
    }
response:
    success (code: 200):
    {
        "message": string,
        "code": int
    }

    error (code: 400, 401, 500):
    {
        "message": string,
        "code": int,
        "error": string
    }

```

```
name: update-photo
method: PATCH
path: /api/users/v1/update-photo
description:
    Memperbarui photo user yang sedang login.
    User hanya dapat memperbarui data miliknya sendiri.
auth: bearer
request:
    body:
    {
        "photo_url": string
    }
response:
    success (code: 200):
    {
        "message": string,
        "code": int
    }

    error (code: 400, 401, 500):
    {
        "message": string,
        "code": int,
        "error": string
    }

```

```
name: update-password
method: PATCH
path: /api/users/v1/update-password
description:
    Memperbarui password user yang lupa.
    User hanya dapat memperbarui data miliknya sendiri.
auth: bearer (berbeda dengan access token)
request:
    body:
    {
        "new_password": string
    }
response:
    success (code: 200):
    {
        "message": string,
        "code": int
    }

    error (code: 400, 401, 500):
    {
        "message": string,
        "code": int,
        "error": string
    }

```

```
name: reset-password
method: PATCH
path: /api/users/v1/reset-password
description:
    Memperbarui password user yang sedang login.
    User hanya dapat memperbarui data miliknya sendiri.
auth: bearer (berbeda dengan access token)
request:
    body:
    {
        "new_password": string,
        "old_password": string
    }
response:
    success (code: 200):
    {
        "message": string,
        "code": int
    }

    error (code: 400, 401, 500):
    {
        "message": string,
        "code": int,
        "error": string
    }

```

## Class

```
name: get-list-class
method: GET
path: /api/classes/v1/get-list
description:
    Mengambil daftar kelas yang diikuti oleh user.
    Data mencakup role user di masing-masing kelas.
    Include paginate dan filter
auth:
    bearer
request:
    query
    - latest
    - longest
    - name
    - page
    - limit
    - pagination
response:
    success (code: 200):
    {
        "message": "success",
        "code": int,
        "result": {
            "items": [
                {
                    "class_id": string,
                    "name": string,
                    "description": string,
                    "role": string,
                    "joined_at": string
                }
            ]
        }
    }

    error (code: 401, 500):
    {
        "message": string,
        "code": int,
        "error": string
    }
```

```
name: create-class
method: POST
path: /api/classes/v1/create
description:
    Membuat kelas baru di dalam sistem EduTask.
    User yang membuat kelas akan otomatis tergabung sebagai teacher.
auth: bearer
request:
    body:
    {
        "name": string,
        "description": string
    }
response:
    success (code: 201):
    {
        "message": string,
        "code": int,
        "result": {
            "id": string,
            "name": string,
            "description": string,
            "created_by": string,
            "created_at": string
        }
    }

    error (code: 400, 401, 500):
    {
        "message": string,
        "code": int,
        "error": string
    }

```

```
name: join-class
method: POST
path: /api/classes/v1/join
description:
    Bergabung ke dalam kelas menggunakan join code.
    User yang bergabung akan otomatis mendapatkan role student.
auth:
    bearer
request:
    body:
    {
        "join_code": string
    }
response:
    success (code: 200):
    {
        "message": string,
        "code": int,
        "result": {
            "class_id": string,
            "role": "student",
            "joined_at": string
        }
    }

    error (code: 400, 401, 404, 409, 500):
    {
        "message": string,
        "code": int,
        "error": string
    }

```

```
name: update-class
method: PUT
path: /api/classes/v1/update/{class_id}
description:
    Memperbarui informasi kelas.
    Hanya teacher di dalam kelas yang dapat melakukan aksi ini.
auth: bearer (teacher in class)
request:
    body:
    {
        "name": string,
        "description": string
    }
response:
    success (code: 200):
    {
        "message": string,
        "code": int,
        "result": {
            "id": string,
            "name": string,
            "description": string,
            "updated_at": string
        }
    }

    error (code: 400, 401, 403, 404, 500):
    {
        "message": string,
        "code": int,
        "error": string
    }

```

```
name: delete-class
method: DELETE
path: /api/classes/v1/delete/{class_id}
description:
    Menghapus kelas dari sistem.
    Hanya teacher (creator / authorized teacher) yang dapat menghapus kelas.
auth:
    bearer (teacher in class)
response:
    success (code: 200):
    {
        "message": string,
        "code": int
    }

    error (code: 401, 403, 404, 500):
    {
        "message": string,
        "code": int,
        "error": string
    }

```

## Ticket

```
name: add-ticket
method: POST
path: /api/tickets/v1/add
description:
    Membuat ticket (tugas) baru.
    Ticket awalnya berstatus draft.
auth:
    bearer (teacher in class)
request:
    body:
    {
        "title": string,
        "description": string,
        "ticket_type": "individual" | "group",
        "deadline": string,
        "classes_ids": ["string"]
    }

response:
    success (code: 201):
    {
        "message": string,
        "code": int,
        "result": {
            "id": string,
            "title": string,
            "ticket_type": string,
            "status": "draft",
            "deadline": string,
            "created_at": string
        }
    }

    error (code: 400, 401, 403, 404, 500):
    {
        "message": string,
        "code": int,
        "error": string
    }
```
