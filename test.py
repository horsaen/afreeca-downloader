import os
import time
from tabulate import tabulate
from threading import Thread

import plugins.bigo.main2 as thread

# Start the random number generation thread
thread.random_number_thread.start()

# Continuously retrieve and print the random number every 4 seconds
while True:
    try:
        random_number = thread.random_number
        if random_number is not None:
            print("Random Number in thread.py:", random_number)
        time.sleep(4)
    except KeyboardInterrupt:
        # Stop the random number generation thread when you press Ctrl+C
        thread.stop_random_number_generator()
        break



# users = []

# with open('users', 'r') as file:
#     for line in file:
#         uname, platform = line.strip().split(',')
#         users.append([uname, platform])

#     for user in users:
#         name, platform = user
#         thread = Thread(target=main, args=(name, platform))
#         thread.start()
#         # main(name, platform)
#         # time.sleep(1)