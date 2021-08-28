# Bookings and Reservations

This is the repository for my booking and reservation project.

- Built in Go version 1.16
- Uses the [go-chi/cho](https://github.com/go-chi/chi) router
- Uses [alex edwards SCS](https://github.com/alexedwards/scs) for session management
- Uses [nosurf](https://github.com/justinas/nosurf) for CSRF

Create & Run Docker Container

```dockerfile
$ docker build -t  smart-hotel-app .
$ docker run -p 8080:80 smart-hotel-app
$ docker ps --format="ID\t{{.ID}}\nNAME\t{{.Names}}\nIMAGE\t{{.Image}}\nPORTS\t{{.Ports}}\nCOMMAND\t{{.Command}}\nCREATED\t{{.CreatedAt}}\nSTATUS\t{{.Status}}\n"
```
