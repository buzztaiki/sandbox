import "@typespec/http";
import "@typespec/openapi";

using Http;
using OpenAPI;

@service(#{ title: "Pet store schema" })
@info(#{ version: "3.0" })
@server("/v3")
namespace PetStore;

model Pet {
  @example(10)
  id?: int64;

  @example("doggie")
  name: string;

  photoUrls?: string[];
  status?: PetStatus;
}

/** pet status in the store. */
enum PetStatus {
  available,
  pending,
  sold,
}

@route("/pet")
@tag("pet")
namespace Pets {
  /** Add a new pet to the store */
  @summary("Add a new pet to the store")
  @post
  @operationId("addPet")
  op addPet(
    /** Create a new pet in the store */
    @body pet: Pet,
  ): Pet;

  /** Returns a single pet */
  @summary("Find pet by ID")
  @get
  @operationId("getPetById")
  op getPetById(
    /** petId ID of pet to return */
    @path petId: int64,
  ): Pet | NotFoundResponse;

  @summary("Updates a pet in the store")
  @post
  @operationId("updatePet")
  op updatePet(
    /** ID of pet that needs to be updated */
    @path petId: int64,

    /** Name of pet that needs to be updated */
    @query name?: string,

    /** Status of pet that needs to be updated */
    @query status?: PetStatus,
  ): OkResponse;

  @summary("Deletes a pet")
  @delete
  @operationId("deletePet")
  op deletePet(
    /** petId Pet ID to delete */
    @path petId: int64,
  ): OkResponse;
}
