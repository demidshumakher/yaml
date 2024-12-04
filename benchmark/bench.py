import subprocess
import datetime


n = 100
print(n)


# now = datetime.datetime.now()
# for i in range(n):
#     subprocess.run(["python", "reg.py", "input.yaml", "output.json"], stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
# time = datetime.datetime.now() - now
# print("regex python:", time)


now = datetime.datetime.now()
for i in range(n):
    subprocess.run(["./regexGo.exe", "input.yaml", "output.json"], stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)

time = datetime.datetime.now() - now
print("regex go:", time)


now = datetime.datetime.now()
for i in range(n):
    subprocess.run(["./main.exe", "input.yaml", "output.json", "lib"], stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)

time = datetime.datetime.now() - now
print("converter:", time )


now = datetime.datetime.now()
for i in range(n):
    subprocess.run(["./main.exe", "input.yaml", "output.json"], stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
time = datetime.datetime.now() - now
print("my convert:", time)