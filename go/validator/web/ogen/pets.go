package main

import (
	"context"
	"sync"

	"github.com/buzztaiki/sandbox/go/validator/web/ogen/petstore"
)

var _ petstore.Handler = (*petsService)(nil)

type petsService struct {
	pets map[int64]petstore.Pet
	id   int64
	mux  sync.Mutex
}

func (p *petsService) AddPet(ctx context.Context, req *petstore.Pet) (*petstore.Pet, error) {
	p.mux.Lock()
	defer p.mux.Unlock()

	p.id++
	newPet := *req
	newPet.ID.SetTo(p.id)
	p.pets[p.id] = newPet

	return &newPet, nil
}

func (p *petsService) DeletePet(ctx context.Context, params petstore.DeletePetParams) error {
	p.mux.Lock()
	defer p.mux.Unlock()

	delete(p.pets, params.PetId)
	return nil
}

func (p *petsService) GetPetById(ctx context.Context, params petstore.GetPetByIdParams) (petstore.GetPetByIdRes, error) {
	p.mux.Lock()
	defer p.mux.Unlock()

	pet, ok := p.pets[params.PetId]
	if !ok {
		// Return Not Found.
		return &petstore.GetPetByIdNotFound{}, nil
	}
	return &pet, nil
}

func (p *petsService) UpdatePet(ctx context.Context, params petstore.UpdatePetParams) error {
	p.mux.Lock()
	defer p.mux.Unlock()

	pet := p.pets[params.PetId]
	pet.Status = params.Status
	if val, ok := params.Name.Get(); ok {
		pet.Name = val
	}
	p.pets[params.PetId] = pet

	return nil
}
