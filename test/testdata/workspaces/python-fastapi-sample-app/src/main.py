#!/usr/bin/env python
from fastapi import FastAPI
from uvicorn import run

app = FastAPI()


@app.get("/")
async def read_root():
    return {"msg": "hello world!"}


if __name__ == "__main__":
    run(
        "main:app",
        host="localhost",
        port=8080,
        reload=False,
        log_level="info",
    )
