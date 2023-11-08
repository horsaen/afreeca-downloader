import threading
import time
import random

# Define a shared variable to store the random number
random_number = None
stop_event = threading.Event()

# Function to continuously generate a random number and update the shared variable
def generate_random_number():
    global random_number
    while not stop_event.is_set():
        random_number = random.randint(1, 100)
        time.sleep(1)

# Function to stop the random number generation thread
def stop_random_number_generator():
    stop_event.set()

# Create a thread that runs the generate_random_number function
random_number_thread = threading.Thread(target=generate_random_number)
