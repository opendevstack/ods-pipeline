import pytest
from httpx import AsyncClient

from main import app


@pytest.mark.anyio
async def test_main_base_endpoint_should_return_hello_world():
    async with AsyncClient(app=app, base_url="http://test") as ac:
        response = await ac.get("/")
    assert response.status_code == 200
    assert response.json() == {"msg": "hello world!"}
