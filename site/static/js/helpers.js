/**
 * Fills template string from session storage
 * 
 * replace {key} placeholders within template with value found under "key" in session storage
 * @param {string} template - template
 * @returns {string} Filled template
 * @throws {Error} If some key is missing from session storage
 */
function fillTemplateFromSession(template) {
  return template.replace(/{(\w+)}/g, (match, key) => {
    const value = sessionStorage.getItem(key);
    if (value === null) {
      throw new Error(`Missing session value for: ${key}`);
    }
    return value;
  });
}

/**
 * Parses JWT token and returns claims
 * @param {string} token - JWT Token
 * @returns {object} Claims
 */
function parseJwt(token) {
  const payload = token.split('.')[1];
  return JSON.parse(atob(payload));
}


/**
 * Starts Authentication timer
 * Transitions via TIMEOUTLOGOUT after timer expires
 *  
 * @param {number} exp - expiration time of authentication
 * @param {string} elId - element where to write the remaining time
 * @param {function} expFunc - expiration function
 * @returns {object} Timer
 */
function startAuthTimer(exp, elId, expFunc) {
  let timerId;

  // create a function which is called periodically
  function updateTimer() {
    const now = Math.floor(Date.now() / 1000);
    let remaining = exp - now;

    // when expired
    if (remaining <= 0) {
      remaining = 0;
      // if timer var exists => timer is running
      if (timerId) {
        // clear
        clearInterval(timerId);
        // call the expiration function
        expFunc()
      }
    }

    // display the remaining time in the element
    const minutes = Math.floor(remaining / 60);
    const seconds = remaining % 60;
    const el = document.getElementById(elId);
    if (el) {
      el.textContent = `${minutes}:${seconds.toString().padStart(2, "0")}`;
    }
  }

  // 1st run to display the remaining time immediately
  updateTimer();
  // start timer with 1s update interval
  timerId = setInterval(updateTimer, 1000);

  // return timer
  return timerId;
}

function createTable(elId, data, headerFunc, rowFunc) {

  if (!data) {
    return;
  }

  if (data.length == 0) {
    return;
  }

  // create table

  const table = document.createElement("table");

  // build header

  const header = table.createTHead();
  const headerRow = header.insertRow();

  if(headerFunc) {
    headerFunc(headerRow, data[0])
  } else {
    for (const key of Object.keys(data[0])) {
      const cell = headerRow.insertCell();
      cell.textContent = key;
    }
  }

  // build data
  const body = table.createTBody();

  for (const e of data) {
    const row = body.insertRow();

    if(rowFunc) {
      rowFunc(row, e)
    } else {
      for (const value of Object.values(e)) {
        const cell = row.insertCell();
        cell.textContent = value;
      }
    }
  }

  // append table

  el = document.getElementById(elId);
  if (el) {
    el.appendChild(table);
  }

  return;
}

function solverListingHeaderFunc(headerRow, data) {
  const headers = ["Name", "Parts 1", "Parts 2"];

  for (const header of headers) {
      const cell = document.createElement("th")
      cell.textContent = header;
      headerRow.appendChild(cell)
  }
}

function solverListingRowFunc(row, data) {
  const dayCell = row.insertCell();
  dayCell.textContent = data.name;
  dayCell.classList.add("row-key");

  const part1Cell = row.insertCell();
  const part1Link = document.createElement("a");
  part1Link.href = "#";
  part1Link.textContent = "Part 1";
  part1Link.dataset.day = data.name;
  part1Link.dataset.part = "1";
  part1Cell.appendChild(part1Link);

  const part2Cell = row.insertCell();
  const part2Link = document.createElement("a");
  part2Link.href = "#";
  part2Link.textContent = "Part 2";
  part2Link.dataset.day = data.name;
  part2Link.dataset.part = "2";
  part2Cell.appendChild(part2Link);
}