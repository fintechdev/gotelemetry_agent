/* Load a series */

$series = series(name:"cpu_load")

/* Populate output properties using data from the series */

value = $series.last()
sparkline = $series.aggregate(func:"avg", interval:"5s", count:100)
