from sqlalchemy import create_engine, text
from sqlalchemy.orm import sessionmaker

from mcp_server.server.config import get_settings

APP_CONFIG = get_settings()

DB_URL = "postgresql+psycopg://{user}:{password}@{host}:{port}/{db}".format(
    user=APP_CONFIG.DB_USER,
    password=APP_CONFIG.DB_PASSWORD,
    host=APP_CONFIG.DB_HOST,
    port=APP_CONFIG.DB_PORT,
    db=APP_CONFIG.DB_NAME,
)
try:
    _engine = create_engine(url=DB_URL)
    Session = sessionmaker(bind=_engine)

    with Session() as s:
        res = s.execute(text("SELECT 1")).scalar()
        print("Health check: ", res)
except Exception as e:
    print(f"Error connecting to the database: {e}")
    raise


