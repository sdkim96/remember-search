import os
from pathlib import Path
from functools import lru_cache
from pydantic import BaseModel
from dotenv import load_dotenv

ROOT_DIR = Path(__file__).parent.parent

class AppConfig(BaseModel):
    PROJECT_NAME: str 
    PROJECT_AUTHOR: str
    DB_NAME: str
    DB_USER: str
    DB_PASSWORD: str
    DB_HOST: str
    DB_PORT: int
    OPENAI_API_MAX_QUOTAS: str
    OPENAI_API_KEY: str
    ELASTIC_HOST: str
    ELASTIC_API_KEY: str

load_dotenv(dotenv_path=os.path.join(ROOT_DIR, '.env'))

@lru_cache
def get_settings() -> AppConfig:
    return AppConfig(
        PROJECT_NAME=os.getenv("PROJECT_NAME", "MCP Server"),
        PROJECT_AUTHOR=os.getenv("PROJECT_AUTHOR", "Your Name"),
        DB_NAME=os.getenv("DB_NAME", "mcp_db"),
        DB_USER=os.getenv("DB_USER", "mcp_user"),
        DB_PASSWORD=os.getenv("DB_PASSWORD", "mcp_password"),
        DB_HOST=os.getenv("DB_HOST", "localhost"),
        DB_PORT=int(os.getenv("DB_PORT", 5432)),
        OPENAI_API_MAX_QUOTAS=os.getenv("OPENAI_API_MAX_QUOTAS", "1000"),
        OPENAI_API_KEY=os.getenv("OPENAI_API_KEY", ""),
        ELASTIC_HOST=os.getenv("ELASTIC_HOST", "localhost:9200"),
        ELASTIC_API_KEY=os.getenv("ELASTIC_API_KEY", ""),
    )