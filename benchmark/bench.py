import subprocess
import datetime


n = 100
print(n)


now = datetime.datetime.now()
for i in range(n):
    subprocess.run(["yaml-to-json-go", "convert", "input.yaml", "output.json"], stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)

time = datetime.datetime.now() - now
print("util:", time )


now = datetime.datetime.now()
for i in range(n):
    subprocess.run(["./myConverter.exe", "input.yaml", "output.json"], stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
time = datetime.datetime.now() - now
print("my:", time)


now = datetime.datetime.now()
for i in range(n):
    subprocess.run(["python", "reg.py", "input.yaml", "output.json"], stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
time = datetime.datetime.now() - now
print("regex:", time)


now = datetime.datetime.now()
for i in range(n):
    subprocess.run(["./regexGo.exe", "input.yaml", "output.json"], stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)

time = datetime.datetime.now() - now
print("regexGo:", time)
