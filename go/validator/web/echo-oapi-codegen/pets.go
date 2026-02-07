package main

import (
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
func (h *PetHandler) AddPet(c echo.Context) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	var pet petstore.Pet
	c.Bind(&pet)

	h.lastID += 1
	index := h.lastID
	pet.Id = &index
	h.pets[index] = &pet

	return c.JSON(http.StatusOK, &pet)
}

// Deletes a pet
// (DELETE /pet/{petId})
func (h *PetHandler) DeletePet(c echo.Context, petId int64) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	_, ok := h.pets[petId]
	if !ok {
		return echo.ErrNotFound
	}
	delete(h.pets, petId)
	return c.NoContent(http.StatusOK)
}

// Find pet by ID
// (GET /pet/{petId})
func (h *PetHandler) GetPetById(c echo.Context, petId int64) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	pet, ok := h.pets[petId]
	if !ok {
		return echo.ErrNotFound
	}

	return c.JSON(http.StatusOK, &pet)
}

// Updates a pet in the store
// (POST /pet/{petId})
func (h *PetHandler) UpdatePet(c echo.Context, petId int64, params petstore.UpdatePetParams) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	pet, ok := h.pets[petId]
	if !ok {
		return echo.ErrNotFound
	}

	if params.Name != nil {
		pet.Name = *params.Name
	}
	if params.Status != nil {
		pet.Status = params.Status
	}
	h.pets[petId] = pet

	return c.NoContent(http.StatusOK)
}
