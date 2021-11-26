import pytest

# see docs https://anyio.readthedocs.io/en/stable/testing.html
@pytest.fixture
def anyio_backend():
    return 'asyncio'
