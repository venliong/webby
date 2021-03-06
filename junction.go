package webby

// Demuxer for Request Method. Implement RouteHandler interface!
type Junction struct {
	ALL, GET, POST, HEAD, DELETE, PUT, PATCH, OPTIONS, AJAX, WS RouteHandler
}

func (jn Junction) View(w *Web) {
	if w.IsWebSocketRequest() {
		if jn.WS != nil {
			jn.WS.View(w)
			return
		}
	}

	switch w.Req.Method {
	case "GET":
		if w.IsAjaxRequest() {
			if jn.AJAX != nil {
				jn.AJAX.View(w)
				return
			}
		}
		if jn.GET != nil {
			jn.GET.View(w)
			return
		}
	case "POST":
		if w.IsAjaxRequest() {
			if jn.AJAX != nil {
				jn.AJAX.View(w)
				return
			}
		}
		if jn.POST != nil {
			jn.POST.View(w)
			return
		}
	case "HEAD":
		if jn.HEAD != nil {
			jn.HEAD.View(w)
			return
		}
	case "DELETE":
		if jn.DELETE != nil {
			jn.DELETE.View(w)
			return
		}
	case "PUT":
		if jn.PUT != nil {
			jn.PUT.View(w)
			return
		}
	case "PATCH":
		if jn.PATCH != nil {
			jn.PATCH.View(w)
			return
		}
	case "OPTIONS":
		if jn.OPTIONS != nil {
			jn.OPTIONS.View(w)
			return
		}
	}

	if jn.ALL != nil {
		jn.ALL.View(w)
		return
	}

	w.Error404()
	return
}

// Chainable version of Junction
type PipeJunction struct {
	jn Junction
}

func NewJunction() PipeJunction {
	return PipeJunction{Junction{}}
}

func (pi PipeJunction) GetJunction() Junction {
	return pi.jn
}

func (pi PipeJunction) Get(get RouteHandler) PipeJunction {
	pi.jn.GET = get
	return pi
}

func (pi PipeJunction) Post(post RouteHandler) PipeJunction {
	pi.jn.POST = post
	return pi
}

func (pi PipeJunction) Head(head RouteHandler) PipeJunction {
	pi.jn.HEAD = head
	return pi
}

func (pi PipeJunction) Delete(del RouteHandler) PipeJunction {
	pi.jn.DELETE = del
	return pi
}

func (pi PipeJunction) Put(put RouteHandler) PipeJunction {
	pi.jn.PUT = put
	return pi
}

func (pi PipeJunction) Patch(patch RouteHandler) PipeJunction {
	pi.jn.PATCH = patch
	return pi
}

func (pi PipeJunction) Options(options RouteHandler) PipeJunction {
	pi.jn.OPTIONS = options
	return pi
}

func (pi PipeJunction) Ajax(ajax RouteHandler) PipeJunction {
	pi.jn.AJAX = ajax
	return pi
}

func (pi PipeJunction) Websocket(ws RouteHandler) PipeJunction {
	pi.jn.WS = ws
	return pi
}

func (pi PipeJunction) All(all RouteHandler) PipeJunction {
	pi.jn.ALL = all
	return pi
}

func (pi PipeJunction) Any(any RouteHandler) PipeJunction {
	return pi.All(any)
}

func (pi PipeJunction) Fallback(fallback RouteHandler) PipeJunction {
	return pi.All(fallback)
}
