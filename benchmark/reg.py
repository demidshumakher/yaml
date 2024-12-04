import re
import json
import sys

input = sys.argv[1]
output = sys.argv[2]

with open(input, "r", encoding="utf-8") as file:
    yaml_content = file.read()

day_pattern = re.compile(r"day_number:\s*(\d+)")
lesson_pattern = re.compile(r"- subject:\s*(.*)\n\s*type:\s*(.*)\n\s*time_start:\s*'(.*)'\n\s*time_end:\s*'(.*)'\n\s*teacher_name:\s*(.*)\n\s*room:\s*'(.*)'\n\s*building:\s*'(.*)'\n\s*group:\s*(.*)\n")

data = []
days = day_pattern.findall(yaml_content)
lessons = lesson_pattern.findall(yaml_content)

n = 0
for day in days:
    day_data = {"day_number": int(day), "lessons": []}
    for lesson in lessons[n:n+2]:
        day_data["lessons"].append({
            "subject": lesson[0],
            "type": lesson[1],
            "time_start": lesson[2],
            "time_end": lesson[3],
            "teacher_name": lesson[4],
            "room": lesson[5],
            "building": lesson[6],
            "group": lesson[7],
        })
    data.append(day_data)
    n += 2

json_output = json.dumps(data, ensure_ascii=False, indent=2)

with open(output, "w", encoding="utf-8") as json_file:
    json_file.write(json_output)

