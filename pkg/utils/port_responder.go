package utils

type Response struct {
	statusCode int
	body       string
}

func (r *Response) GetStatus() int {
	return r.statusCode
}

func (r *Response) GetBody() string {
	return r.body
}

type Responder struct {
	pathResponders map[string]*PathResponder
}

func (r *Responder) GetResponse(path string) *Response {
	if r != nil {
		pr := r.pathResponders[path]
		if pr != nil {
			resp := pr.getNext()
			if resp != nil {
				return resp
			}
		}
	}
	return &Response{statusCode: 404, body: path + " Not Found!"}
}

func (r *Responder) AddPathResponse(path, body string) *Responder {
	_ = r.AddPath(path).Add(body)
	return r
}

func (r *Responder) AddPath(path string) *PathResponder {
	pr := r.pathResponders[path]
	if pr == nil {
		pr = &PathResponder{owner: r, path: path}
		if r.pathResponders == nil {
			r.pathResponders = map[string]*PathResponder{path: pr}
		} else {
			r.pathResponders[path] = pr
		}
	}
	return pr
}

func (r *PathResponder) getNext() (resp *Response) {
	rLen := len(r.responses)
	if rLen != 0 {
		i := r.next
		if rLen <= i {
			i = 0
		}
		resp = r.responses[i]
		r.next = i + 1
	}
	return
}

type PathResponder struct {
	owner     *Responder
	path      string
	responses []*Response
	next      int
}

func (r *PathResponder) Add(body string) *PathResponder {
	return r.AddWithStatus(200, body)
}

func (r *PathResponder) AddWithStatus(statusCode int, body string) *PathResponder {
	r.plus(&Response{statusCode: statusCode, body: body})
	return r
}

func (r *PathResponder) Repeat(count int) *PathResponder {
	if count <= 0 {
		panic("Repeat < 1")
	}
	rLen := len(r.responses)
	if rLen == 0 {
		panic("Can't Repeat - None added yet!")
	}
	resp := r.responses[rLen-1]
	for i := 0; i < count; i++ {
		r.plus(resp)
	}
	return r
}

func (r *PathResponder) NoMore() *Responder {
	return r.owner
}

func (r *PathResponder) plus(resp *Response) {
	r.responses = append(r.responses, resp)
}
