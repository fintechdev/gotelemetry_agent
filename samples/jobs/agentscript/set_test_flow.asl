/* Load a series */

$series = series(name:"cpu_load")
$last = $series.avg("5s")
$sparkline = $series.aggregate(func:"avg", interval:"5s", count:100)

$series.trim(since:"10s")

/* Populate output properties using data from the series */

value = $last
sparkline = $sparkline

/* Check for anomalies */

if anomaly(data: $sparkline.values, value: $last) {
	color = "red"
} else {
	color = "white"
}
