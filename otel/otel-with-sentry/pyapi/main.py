import os

import httpx
import sentry_sdk
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel

NEXT_URL = os.getenv("NEXT_URL", "http://localhost:8000")

sentry_sdk.init(
    dsn=os.getenv("SENTRY_DSN"),
    traces_sample_rate=None,
)
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
    resp.raise_for_status()
    data = Result.model_validate(resp.json())
    return Result(result=data.result + n)


@app.get("/error")
def error():
    raise RuntimeError("something went wrong")
