import json
import os
import io

if __name__ == "__main__":
    cache = "./data/research/lists/"
    # extract industries and sectors
    combines = []
    for exchange in ["nasdaq", "nyse"]:
        location = os.path.join(cache, f"exchange_is_{exchange}.json")
        if not os.path.exists(location):
            print(f"[error] could not find json file for {exchange}")
            continue
        with open(location, "r") as f:
            data = json.load(f)
            combines.extend(data)
    sectors, industries = set(), set()
    # get the sectors
    for c in combines:
        sectors.add(c["sector"])
        industries.add(c["industry"])
    sectors.remove("")
    industries.remove("")

    # save them
    with open(os.path.join(cache, "sectors.json"), "w") as f:
        json.dump(list(sectors), f)
    with open(os.path.join(cache, "industries.json"), "w") as f:
        json.dump(list(industries), f)
