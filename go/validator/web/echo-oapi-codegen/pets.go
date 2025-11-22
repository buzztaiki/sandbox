package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/buzztaiki/sandbox/go/validator/web/echo-oapi-codegen/petstore"
	"github.com/labstack/echo/v4"
)

type PetHandler struct {
	lastID int64
	mu     sync.Mutex
	pets   map[int64]*petstore.Pet
}

func NewPetHandler() *PetHandler {
	return &PetHandler{
		lastID: 0,
		pets:   map[int64]*petstore.Pet{},
	}
}

// Add a new pet to the store
// (POST /pet)
func (h *PetHandler) AddPet(ctx echo.Context) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	var pet petstore.Pet
	ctx.Bind(&pet)

	h.lastID += 1
	index := h.lastID
	pet.Id = &index
	h.pets[index] = &pet

	return ctx.JSON(http.StatusOK, &pet)
}

// Deletes a pet
// (DELETE /pet/{petId})
func (h *PetHandler) DeletePet(ctx echo.Context, petId int64) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	_, ok := h.pets[petId]
	if !ok {
		return echo.ErrNotFound
	}
	delete(h.pets, petId)
	return ctx.NoContent(http.StatusOK)
}

// Find pet by ID
// (GET /pet/{petId})
func (h *PetHandler) GetPetById(ctx echo.Context, petId int64) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	log.Println("aaa")

	pet, ok := h.pets[petId]
	if !ok {
		return echo.ErrNotFound
	}

	return ctx.JSON(http.StatusOK, &pet)
}

// Updates a pet in the store
// (POST /pet/{petId})
func (h *PetHandler) UpdatePet(ctx echo.Context, petId int64, params petstore.UpdatePetParams) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	pet, ok := h.pets[petId]
	if !ok {
		return echo.ErrNotFound
	}

	pet.Name = *params.Name
	pet.Status = params.Status

	return ctx.NoContent(http.StatusOK)
}
