package server

import (
    "avitoTech/internal/models"
    "avitoTech/internal/repository"
    "errors"
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/jackc/pgx/v5/pgxpool"
    "net/http"
    "strconv"
)

type Server struct {
    router     *gin.Engine
    tenderRepo *repository.TenderRepo
}

func NewServer(conn *pgxpool.Pool) *Server {
    return &Server{
        router:     gin.Default(),
        tenderRepo: repository.NewTenderRepo(conn),
    }
}

func SendError(ctx *gin.Context, err error) {
    switch {
    case errors.Is(err, repository.ErrUserNotFound):
        ctx.JSON(http.StatusUnauthorized, err.Error())
    case errors.Is(err, repository.ErrNoAccess):
        ctx.JSON(http.StatusForbidden, err.Error())
    case errors.Is(err, repository.ErrTenderNotFound):
        ctx.JSON(http.StatusNotFound, err.Error())
    default:
        ctx.JSON(http.StatusInternalServerError, err.Error())
    }
}

func (s *Server) Ping(ctx *gin.Context) {
    ctx.String(http.StatusOK, "ok")
}

func (s *Server) CreateTender(ctx *gin.Context) {
    var tender models.Tender
    err := ctx.BindJSON(&tender)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, err.Error())
        return
    }
    tender, err = s.tenderRepo.CreateTender(ctx.Request.Context(), tender)
    if err != nil {
        SendError(ctx, err)
        return
    }
    ctx.JSON(http.StatusOK, tender)
}

func (s *Server) GetTenders(ctx *gin.Context) {
    var limit = 5
    var offset int
    limitRaw := ctx.Query("limit")
    var err error
    if limitRaw != "" {
        limit, err = strconv.Atoi(limitRaw)
        if err != nil {
            ctx.JSON(http.StatusBadRequest, err)
            return
        }
    }
    offsetRaw := ctx.Query("offset")
    if offsetRaw != "" {
        offset, err = strconv.Atoi(offsetRaw)
        if err != nil {
            ctx.JSON(http.StatusBadRequest, err)
            return
        }
    }
    serviceType := ctx.QueryArray("service_type")
    tenders, err := s.tenderRepo.GetTenders(ctx.Request.Context(), limit, offset, serviceType)
    if err != nil {
        SendError(ctx, err)
        return
    }
    ctx.JSON(http.StatusOK, tenders)
}

func (s *Server) GetMyTenders(ctx *gin.Context) {
    var limit = 5
    var offset int
    limitRaw := ctx.Query("limit")
    var err error
    if limitRaw != "" {
        limit, err = strconv.Atoi(limitRaw)
        if err != nil {
            ctx.JSON(http.StatusBadRequest, err.Error())
            return
        }
    }
    offsetRaw := ctx.Query("offset")
    if offsetRaw != "" {
        offset, err = strconv.Atoi(offsetRaw)
        if err != nil {
            ctx.JSON(http.StatusBadRequest, err.Error())
            return
        }
    }
    username := ctx.Query("username")
    tenders, err := s.tenderRepo.GetTendersByUsername(ctx.Request.Context(), limit, offset, username)
    if err != nil {
        SendError(ctx, err)
        return
    }
    ctx.JSON(http.StatusOK, tenders)
}

func (s *Server) GetTenderStatus(ctx *gin.Context) {
    username := ctx.Query("username")
    tenderId := ctx.Param("tenderId")
    status, err := s.tenderRepo.GetTenderStatus(ctx.Request.Context(), tenderId, username)
    if err != nil {
        SendError(ctx, err)
        return
    }
    ctx.String(http.StatusOK, status)
}

func (s *Server) UpdateTender(ctx *gin.Context) {
    username := ctx.Query("username")
    var tender models.Tender
    tender.Id = ctx.Param("tenderId")
    err := ctx.BindJSON(&tender)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, err.Error())
        return
    }
    tender, err = s.tenderRepo.UpdateTender(ctx.Request.Context(), username, tender)
    if err != nil {
        SendError(ctx, err)
        return
    }
    ctx.JSON(http.StatusOK, tender)
}

func (s *Server) UpdateStatus(ctx *gin.Context) {
    var tender models.Tender
    username := ctx.Query("username")
    tender.Id = ctx.Param("tenderId")
    tender.Status = ctx.Query("status")
    fmt.Println(tender.Status, username)
    tender, err := s.tenderRepo.UpdateTender(ctx.Request.Context(), username, tender)
    if err != nil {
        SendError(ctx, err)
        return
    }
    ctx.JSON(http.StatusOK, tender)
}

func (s *Server) InitRoutes() {
    api := s.router.Group("/api")
    {
        api.GET("/ping", s.Ping)
        tenders := api.Group("/tenders")
        {
            tenders.GET("/", s.GetTenders)
            tenders.GET("/my", s.GetMyTenders)
            g := tenders.Group("/:tenderId")
            {
                g.GET("/status", s.GetTenderStatus)
                g.PUT("/status", s.UpdateStatus)
                g.PATCH("/edit", s.UpdateTender)
            }
        }
        api.POST("/tenders/new", s.CreateTender)
    }
}

func (s *Server) Run(address string) error {
    s.InitRoutes()
    err := s.router.Run(address)
    if err != nil {
        return err
    }
    return nil
}
