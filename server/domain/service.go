package domain

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/leandrohsilveira/simple-bank/server/store"
)

type DomainService struct{}

func (service DomainService) Store(ctx fiber.Ctx) (*store.Store, error) {
	instance, ok := ctx.Locals(store.StoreCtxKey).(*store.Store)

	if !ok {
		return &store.Store{}, fmt.Errorf("store not found in context")
	}

	return instance, nil
}
