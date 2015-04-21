$series: series(name:"cpu_load")

value: $series.last()
sparkline: $series.aggregate(func:"avg", interval:"5s", count:100)
