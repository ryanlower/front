# front

Inspired by [camo](https://github.com/atmos/camo). Proxies images, allowing user embedded content to be served from a single secure host.

[![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy?template=https://github.com/ryanlower/front) [![Circle CI](https://circleci.com/gh/ryanlower/front.png?circle-token=20a52d09d241b53c718e4b93a48e9a8ea3e5c192)](https://circleci.com/gh/ryanlower/front)

---

### Usage

Pass the remote image url to your front deployment as a url param.

For example, to proxy an insecure gopher `http://golang.org/doc/gopher/frontpage.png`:

`https://go-front.herokuapp.com/?url=http://golang.org/doc/gopher/frontpage.png`

![HTTPS Gopher](https://go-front.herokuapp.com/?url=http://golang.org/doc/gopher/frontpage.png)

---

### Deployment

The simplest method is to deploy on heroku:

[![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy?template=https://github.com/ryanlower/front)

---

### Configuration

Config is via environment variables:

variable | description | optional
--- | --- | ---
PORT | The port to listen on | no, though automatically set on heroku
ALLOWED_CONTENT_TYPE_REGEX | Regex the upstream Content-Type must match in order to be proxied | yes, defaults to `^image/`
