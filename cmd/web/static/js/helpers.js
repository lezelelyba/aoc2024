function createTable(data) {
    if(!data.length) return;

    const table = document.createElement("table");

    const thead = document.createElement("thead")
    const headerRow = document.createElement("tr");

    for (const key of Object.keys(data[0])) {
        const th = document.createElement("th");
        th.textContent = key;
        headerRow.appendChild(th);
    }

    thead.appendChild(headerRow)
    thead.appendChild(thead)
   
    const tbody = document.createElement("tbody")

    for (const item of data) {
        const row = document.createElement("tr")
        for (const val of Object.values(item)) {
           const td = document.createElement("td") 
           td.textContent = val
           row.appendChild(td)
        }
        tbody.appendChild(row);
    }

    table.appendChild(tbody)

    return table
}