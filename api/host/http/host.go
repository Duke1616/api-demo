package http

import (
	"github.com/Duke1616/api-demo/api/host"
	"github.com/infraboard/mcube/http/request"
	"github.com/infraboard/mcube/http/response"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (h *handler) CreateHost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// 需要读取用户传递的参数，由于POST请求，我们从body里面取出数据
	body, err := request.ReadBody(r)
	if err != nil {
		response.Failed(w, err)
		return
	}
	h.log.Debugf("receive body: %s", string(body))
	response.Success(w, "ok")
}

func (h *handler) QueryHost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	qs := r.URL.Query()

	var (
		pageSize   = 20
		pageNumber = 1
	)

	// 从query string读取分页参数
	psStr := qs.Get("page_size")
	if psStr != "" {
		pageSize, _ = strconv.Atoi(psStr)
	}
	pnStr := qs.Get("page_number")
	if pnStr != "" {
		pageNumber, _ = strconv.Atoi(pnStr)
	}

	req := &host.QueryHostRequest{
		PageSize:   pageSize,
		PageNumber: pageNumber,
	}

	set, err := h.host.QueryHost(r.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, set)
}
