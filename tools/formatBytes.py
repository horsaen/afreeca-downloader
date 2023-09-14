# super cool wow wonderful amazing
def format_bytes(size):
    # 2**10 = 1024
    power = 2 ** 10
    n = 0
    power_labels = {0: '', 1: 'KB', 2: 'MB', 3: 'GB', 4: 'TB'}
    while size >= power and n < len(power_labels) - 1:
        size /= power
        n += 1
    formatted_size = f'{size:.2f}' if n > 1 else f'{int(size)}'
    return formatted_size + " " + power_labels[n]