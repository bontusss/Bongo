# Bongo Web Framework
### Your go-to web framework for perfectionists on a tight schedule.

## What is Bongo?
#### Bongo is a web application framework the provides all the essential components, allowing you build without worrying about the little things.

## Quickstart
```go
package main

import "github.com/bontusss/bongo/v1"

func main() {
    app := bongo.New()
	
	app.Router.Get("/", handler)
	app.Serve()
}
```

#### Bongo has an inbuilt cli that bootstraps and sets up initial up your project structure for you.
```bash
$ bongo new mysite
```
> NB: This project is actively under development and V1 is neither finished, tested or released!

### Features and stacks
Version 1 will stand on the shoulder of giants to implement most core features. I started this project to have full grasp
on golang and fully understand how the web works. Version 2 will have custom implementations for these core features.

- [x] Router and middlewares [Chi](https://github.com/go-chi/chi) 
- [x] Configuration => ```Bongo new``` command creates a project directory for your app with a default app config in the
```.env``` file in the projects root folder.
- [x] Logging => The built-in logger is [zap](https://github.com/uber-go/zap) with [lumberjack](https://gopkg.in/natefinch/lumberjack).
- [x] Templating engine =>  Currently supports Go templates and [Jet](https://github.com/go-jet/jet) templates.
- [ ] Database => Will support Sqlite, Postgres, mysql and mariadb
- [ ] Internalization
- [ ] Forms
- [x] Sessions and cookies [scs](https://github.com/alexedwards/scs)
- [ ] Emailing
- [ ] Caching
- [ ] CLI
- [ ] Validation
- [ ] Admin interface
- [ ] ORM
- [ ] Docs

### Contributing
This repo is currently not accepting code PRs. If you want to help out with documents, It'll be a pleasure to work with you.

### License
[MIT](https://github.com/bontusss/bongo/blob/main/LICENSE)