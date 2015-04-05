# front

Proxies and manipulates images, allowing user embedded content to be served from a single secure host at various sizes.

[![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy?template=https://github.com/ryanlower/front) [![Circle CI](https://circleci.com/gh/ryanlower/front.png?circle-token=20a52d09d241b53c718e4b93a48e9a8ea3e5c192)](https://circleci.com/gh/ryanlower/front)

---

### Usage

**As proxy**

Pass the remote image url to your front deployment as a url param. For example, to proxy an insecure gopher `http://golang.org/doc/gopher/frontpage.png`:

[`https://go-front.herokuapp.com/?url=http://golang.org/doc/gopher/frontpage.png`](https://go-front.herokuapp.com/?url=http://golang.org/doc/gopher/frontpage.png)

![HTTPS Gopher](https://go-front.herokuapp.com/?url=http://golang.org/doc/gopher/frontpage.png)

**To resize images**

Add width and height params to the standard proxy url. For example, for a smaller gopher:

[`https://go-front.herokuapp.com/?url=http://golang.org/doc/gopher/frontpage.png&width=125&height=175`](https://go-front.herokuapp.com/?url=http://golang.org/doc/gopher/frontpage.png&width=125&height=175)

![HTTPS Gopher](https://go-front.herokuapp.com/?url=http://golang.org/doc/gopher/frontpage.png&width=125&height=175)

---

### Deployment

The simplest method is to deploy on heroku:

[![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy?template=https://github.com/ryanlower/front)

You'll probably want to add a CDN (e.g. CloudFront) to your front deployment for caching

---

### Configuration

Config is via environment variables:

variable | description | optional
--- | --- | ---
PORT | The port to listen on | no, though automatically set on heroku
ALLOWED_CONTENT_TYPE_REGEX | Regex the upstream Content-Type must match in order to be proxied | yes, defaults to `^image/`
