import os

import httpx
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel

NEXT_URL = os.getenv("NEXT_URL", "http://localhost:8000")

app = FastAPI()


class Result(BaseModel):
    result: int


@app.get("/sum/{n}")
async def sum(n: int) -> Result:
    if n == 0:
        return Result(result=0)
    if n < 0:
        raise HTTPException(status_code=422, detail="n must be a non-negative integer")

    async with httpx.AsyncClient() as client:
        resp = await client.get(f"{NEXT_URL}/sum/{n - 1}")
    data = Result.model_validate(resp.json())
    return Result(result=data.result + n)
