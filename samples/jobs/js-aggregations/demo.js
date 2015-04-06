#!/usr/bin/env node

'use strict';

// var op = {
// 	$trim: {
// 		series: 'test',
// 		keep: 10
// 	}
// };

// console.log(JSON.stringify(op));

var v = Math.random() * 1000;

if (Math.random() > 0.8) {
	v *= 10000;
}

var op = {
	$push: {
		series: 'test',
		value: v,
	}
};

console.log(JSON.stringify(op));

op = {
	value: {
		$pick: {
			prop: 'value',
			from: {
				$last: {
					series: 'test',
					default: {
						ts: 0,
						value: 0
					}
				}
			}
		}
	},

	color: {
		$if: {
			condition: {
				$anomaly: {
					series: 'test',
					period: 100,
					value: {
						$pick: {
							prop: 'value',
							from: {
								$last: {
									series: 'test',
									default: {
										ts: 0,
										value: 0
									}
								}
							}
						}
					}
				}
			},
			then: 'red',
			else: 'white'
		}
	},

	sparkline: {
		$pick: {
			prop: 'value',
			from: {
				$aggregate: {
					series: 'test',
					op: 'sum',
					interval: 5,
					count: 20
				}
			}
		}
	}
};

console.log(JSON.stringify(op));

// op = {
// 	value: {
// 		$pick: {
// 			prop: 'value',
// 			from: {
// 				$last: {
// 					series: 'test'
// 				}
// 			}
// 		}
// 	}
// };

// console.log(JSON.stringify(op));