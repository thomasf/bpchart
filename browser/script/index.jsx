import React from 'react';
import Chart from './components/Chart';
import request from 'superagent';
import moment from 'moment';
import d3 from 'd3';

require("c3/c3.css");
require("../style/style.css");

// Scoring rules: (TODO, start at 0, no negative values, use array)
//
//  0 = normal
//  1 = warning high
//  2 = high
//  -1 = warning low

var scoreSys = function(v) {
    if (v > 139) {
        return 2;
    }
    if (v > 119) {
        return 1;
    }
    if (v < 90) {
        return -1;
    }
    return 0;
}


var scoreDia = function(v) {
    if (v > 89) {
        return 2;
    }
    if (v > 79) {
        return 1;
    }
    if (v < 60) {
        return -1;
    }
    return 0;
}

var scorePulse = function(v) {
    if (v > 90) {
        return 2;
    }
    if (v > 80) {
        return 1;
    }
    return 0;
}

var bpColors = {
    0: "#888",
    1: "#000",
    2: "#f0f",
}

var colorBp = function(score){
    // todo
}



var requestData;

var prepareData = function(data, dataOptions) {
    if (dataOptions == null) {
        dataOptions = {};
    }

    if (dataOptions.timeSeries) {
        var avgData = [];
        var prevTime;
        data.forEach(function(v) {
            if (prevTime == null) {
                avgData.push(v);
                prevTime = new Date(v.time);
                return null;
            }
            var curTime = new Date(v.time);
            if ((curTime-prevTime) < (10 * 60 * 1000)) {
                var prev = avgData.pop();
                var newAvg = {
                    time: v.time,
                    sys: (prev.sys + v.sys) / 2,
                    dia: (prev.dia + v.dia) / 2,
                    pulse: (prev.pulse + v.pulse) / 2,
                }
                avgData.push(newAvg);
                prevTime = curTime;
                return null;
            }
            avgData.push(v)
            prevTime = curTime;
            return null;
        })
        data = avgData;
    }

    var chartData = {
        dateStrings: [],
        date: ['date'],
        time: ['time'],
        ntime: ['ntime'],
        sys: ['sys'],
        dia: ['dia'],
        pulse: ['pulse'],
    };
    data.forEach(function(v){
        var d = moment(v.time).format('YYYY-MM-DD');
        var dt = moment(v.time).format('YYYY-MM-DD HH:mm:ss');
        var t = moment(v.time).format('HH:mm:ss');
        var mt = moment(v.time);
        var normt = (mt.hour()*3600) + (mt.minute()*60) + (mt.second());

        chartData.dateStrings.push(v.time);
        chartData.date.push(dt);
        chartData.time.push(t);
        chartData.ntime.push(normt);
        chartData.sys.push(v.sys);
        chartData.dia.push(v.dia);
        chartData.pulse.push(v.pulse);
    });

    var componentData = {
        size: {
            height:900,
            width:1000,
        },
        data: {
            x: 'date',
            xFormat: '%Y-%m-%d %H:%M:%S',

            axes: {
                'ntime': 'y2',
            },
            labels: {
                format: {
                    sys: function(v){return Math.round(v)},
                    dia: function(v){return Math.round(v)},
                    pulse: function(v){return Math.round(v)},
                },
            },
            columns: [
                chartData.date,
                chartData.pulse,
                chartData.ntime,
                chartData.sys,
                chartData.dia,
            ],
            groups: [
                // ['sys', 'dia']
            ],
            types: {
                ntime: 'bar',
                sys: 'line',
                dia: 'line',
                pulse: 'line',
            },
            names: {
                ntime: "Time",
                sys: "Systolic (max)",
                dia: "Diastolic (min)",
                pulse: "Pulse",
            },
            colors: {
                sys: '#aee',
                dia: '#aee',
                pulse: '#ccc',
                ntime: '#f8f8f8',
            },
            color: function (color, d) {
                // d will be 'id' when called for legends
                //return d.id && d.id === 'data3' ? d3.rgb(color).darker(d.value / 150) : color;
                if (d.id) {
                    if (d.id == 'pulse') {
                        if (d.value) {
                            if (d.value > 90) {
                                return "#e33";
                            }

                            if (d.value > 80) {
                                return "#933";
                            }
                            return "#888";
                        }
                    }
                    if (d.id == 'sys') {
                        if (d.value) {
                            if (d.value > 139) {
                                return "#f0f";
                            }
                            if (d.value > 119) {
                                return "#000";
                            }
                            if (d.value < 90) {
                                return "#f0f";
                            }
                            return "#388";
                        }
                    }
                    if (d.id == 'dia') {
                        if (d.value) {
                            if (d.value > 89) {
                                return "#f0f";
                            }
                            if (d.value > 79) {
                                return "#000";
                            }
                            if (d.value < 60) {
                                return "#f0f";
                            }
                            return "#388";
                        }
                    }
                }
                return color;
            }
        },
        grid: {
            y: {
                lines: [

                ]
            },
            y: {
                lines: [
                    // {value:90, class: 'gridsoft'},
                    // {value:120, class: 'gridsoft'},
                    // {value:140, class: 'gridsoft'},
                    // {value:140, text:"Sys max", position: "start"},
                    // {value:90, text:"Dia max", position: "end"},
                ]
            }
        },
        subchart: {
            show: true
        },
        axis: {
            y2: {
                // type: 'category',
                tick: {
                    format: function(x) {
                        var ts = moment(0).seconds(x);
                        if (ts.dayOfYear() == 1) {
                            return ts.format("HH:mm");
                        }
                        return "";
                    },
                },
                show: true
            },
            x: {
                localtime: true,
                type: dataOptions.timeSeries ? 'timeseries':'category',
                tick: {
                    // fit:  false,
                    culling: {
                        max: 14,
                    },
                    // centered: true,
                    format: function (x) {
                        if (dataOptions.timeSeries) {
                            var dt =  moment(x);
                            return dt.format('MM-DD HH:mm');
                        } else {
                            var dt = moment(chartData.dateStrings[x]);
                            return dt.format('MM-DD HH:mm');
                        }
                    },
                    // values: chartData.date
                    // count: 4,
                }
            }
        }

    };
    return componentData;

}


var render = function(componentData){
    React.render(
            <div>
            <button onClick={function(){
                render(prepareData(requestData));
            }}>all</button>
            <button onClick={function(){
                render(prepareData(requestData, {timeSeries:true}));
            }}>timeseries</button>

            <Chart options={componentData} element="dddd"  />

        </div>,
        document.getElementById('content')
    );
}



request
.get('/json/')
.set('Accept', 'application/json')
.end(function(err, res){
    if (!res.ok) {
        console.log(err);
        return;
    }

    requestData = res.body;
    render(prepareData(requestData, {timeSeries: true}))
});
