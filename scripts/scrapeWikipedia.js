function arrayFromColumn(table, columnNumber) {
  const tr = table.querySelectorAll("tbody tr");

  var tableRows = Array.from(tr).map((x) => {
    var theTD = x.querySelectorAll("td");
    var theText = Array.from(theTD).map((item) => item.innerText);
    return theText;
  });

  var isolatedColumn = tableRows.map((x) => {
    return x[columnNumber];
  });
  return isolatedColumn;
}

// function getTables returns an array of tables on the Wikipedia page. The arg
// startIndex is inclusive but endIndex is not.
function getTables(startIndex, endIndex) {
  const tables = document.querySelectorAll(".wikitable");
  const tablesArray = Array.from(tables);
  return tablesArray.slice(startIndex, endIndex);
}

function extractColumnFromArrayOfTables(tablesArray, columnNumber) {
  for (i = 0; i < tablesArray.length; i++) {}
}
