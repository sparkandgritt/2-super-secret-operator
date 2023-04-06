from bs4 import BeautifulSoup

# Load the HTML file
with open("table.html") as fp:
    soup = BeautifulSoup(fp, "html.parser")

# Find the table element
table = soup.find("table")

# Find the index of the "owner" column
headers = table.find("thead").find_all("th")
owner_index = None
for i, header in enumerate(headers):
    if header.text.strip().lower() == "owner":
        owner_index = i
        break

# Count the number of occurrences of each owner
counts = {}
for tbody in table.find_all("tbody"):
    rows = tbody.find_all("tr")
    for row in rows:
        cells = row.find_all("td")
        if owner_index is not None and cells:
            owner = cells[owner_index].text.strip()
            counts[owner] = counts.get(owner, 0) + 1

# Output the results
for owner, count in counts.items():
    print(f"{owner}: {count}")
