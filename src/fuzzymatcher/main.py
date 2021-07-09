import asyncio
import logging
from robomotion import plugin
from init import *

async def main():
    await plugin.start()

if __name__ == "__main__":
    logging.basicConfig()
    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)
    loop.run_until_complete(main())