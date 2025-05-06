from mcp.server.fastmcp import FastMCP

from mcp_server.server.config import get_settings
import mcp_server.server.modules as m

get_settings()

server = FastMCP()


if __name__ == "__main__":
    server.run()