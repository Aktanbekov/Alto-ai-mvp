import json

with open("result.json", "r") as f:
    data = json.load(f)

messages = data["messages"]
print(messages[3]["text"])  # check how many total messages
