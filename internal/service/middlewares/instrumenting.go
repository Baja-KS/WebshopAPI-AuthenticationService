package middlewares

import (
	"context"
	"fmt"
	"github.com/Baja-KS/WebshopAPI-AuthenticationService/internal/database"
	"github.com/Baja-KS/WebshopAPI-AuthenticationService/internal/service"
	"github.com/go-kit/kit/metrics"
	"time"
)

type InstrumentingMiddleware struct {
	RequestCount   metrics.Counter
	RequestLatency metrics.Histogram
	Next           service.Service
}

func (i *InstrumentingMiddleware) Login(ctx context.Context, username string, password string) (user database.UserOut, token string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Login", "error", fmt.Sprint(err != nil)}
		i.RequestCount.With(lvs...).Add(1)
		i.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	user, token, err = i.Next.Login(ctx, username, password)
	return
}

func (i *InstrumentingMiddleware) Register(ctx context.Context, user database.UserIn) (msg string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Register", "error", fmt.Sprint(err != nil)}
		i.RequestCount.With(lvs...).Add(1)
		i.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	msg, err = i.Next.Register(ctx, user)
	return
}

func (i *InstrumentingMiddleware) GetAll(ctx context.Context) (users []database.UserOut, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetAll", "error", fmt.Sprint(err != nil)}
		i.RequestCount.With(lvs...).Add(1)
		i.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	users, err = i.Next.GetAll(ctx)
	return
}

func (i *InstrumentingMiddleware) AuthUser(ctx context.Context) (user database.UserOut, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "AuthUser", "error", fmt.Sprint(err != nil)}
		i.RequestCount.With(lvs...).Add(1)
		i.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	user, err = i.Next.AuthUser(ctx)
	return
}
