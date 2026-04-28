package http

import (
	"io"
	"net/http"

	"github.com/d-velop/dvelop-app-template-go/domain"
	"github.com/d-velop/dvelop-app-template-go/domain/plugins/conf"
	"github.com/d-velop/dvelop-sdk-go/log"
)

type bcHandler struct {
	assetBasePath string
	renderhtml    func(w io.Writer, data interface{}, templatename string) error
	configStore   domain.BcConfigRepository
	bcClient      domain.DebitorClient
}

func NewBcHandler(
	assetBasePath string,
	renderhtml func(w io.Writer, data interface{}, templatename string) error,
	configStore domain.BcConfigRepository,
	bcClient domain.DebitorClient,
) *bcHandler {
	return &bcHandler{
		assetBasePath: assetBasePath,
		renderhtml:    renderhtml,
		configStore:   configStore,
		bcClient:      bcClient,
	}
}

type BcSetupDto struct {
	BaseHtmlDto
	Title   string
	Config  domain.BcConfig
	Saved   bool
}

type DebitorsDto struct {
	BaseHtmlDto
	Title     string
	Debitoren []domain.Debitor
	Error     string
	HasConfig bool
}

func (h *bcHandler) HandleSetup() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			d := &BcSetupDto{
				Title:  "Business Central Setup",
				Config: h.configStore.GetBcConfig(),
				Saved:  req.URL.Query().Get("saved") == "1",
			}
			d.AssetBasePath = h.assetBasePath
			w.Header().Set("content-type", "text/html;charset=utf-8")
			if err := h.renderhtml(w, d, "bcsetup.html"); err != nil {
				log.Error(req.Context(), err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		case http.MethodPost:
			if err := req.ParseForm(); err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			config := domain.BcConfig{
				BaseURL:  req.FormValue("baseurl"),
				Mandant:  req.FormValue("mandant"),
				Username: req.FormValue("username"),
				Password: req.FormValue("password"),
			}
			h.configStore.SaveBcConfig(config)
			http.Redirect(w, req, conf.BasePath+"/bcsetup?saved=1", http.StatusSeeOther)
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	}
	return http.HandlerFunc(fn)
}

func (h *bcHandler) HandleDebitoren() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			config := h.configStore.GetBcConfig()
			d := &DebitorsDto{
				Title:     "Debitoren",
				HasConfig: config.BaseURL != "",
			}
			d.AssetBasePath = h.assetBasePath

			if config.BaseURL != "" {
				debitoren, err := h.bcClient.GetDebitoren(config)
				if err != nil {
					d.Error = err.Error()
				} else {
					d.Debitoren = debitoren
				}
			}

			w.Header().Set("content-type", "text/html;charset=utf-8")
			if err := h.renderhtml(w, d, "debitoren.html"); err != nil {
				log.Error(req.Context(), err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	}
	return http.HandlerFunc(fn)
}
