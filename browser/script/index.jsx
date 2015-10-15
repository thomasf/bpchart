import React from 'react';
import Chart from './components/Chart';
import request from 'superagent';
import moment from 'moment';
import d3 from 'd3';

require("c3/c3.css");
require("../style/style.css");

var bpColors = {
    "OK": "#388",
    "HighWarning": "#000",
    "LowWarning": "#000",
    "High": "#f0f",
    "Low": "#f0f",
}

var pulseColors = {
    "OK": "#888",
    "HighWarning": "#933",
    "LowWarning": "#933",
    "High": "#e33",
    "Low": "#e33",
}


var requestData;

var prepareData = function(data, dataOptions) {
    if (dataOptions == null) {
        dataOptions = {};
    }

    var chartData = {
        json: [],
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
        var mt = moment(v.time);
        var normt = (mt.hour()*3600) + (mt.minute()*60) + (mt.second());

        chartData.dateStrings.push(v.time);
        chartData.date.push(dt);
        chartData.ntime.push(normt);
        chartData.sys.push(v.sys);
        chartData.dia.push(v.dia);
        chartData.pulse.push(v.pulse);
        chartData.json.push(v);
    });
    var datalen = data.length;

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
                    sys: datalen > 70? "": function(v){return v},
                    dia: datalen > 70? "": function(v){return v},
                    pulse: datalen > 70? "": function(v){return v},
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
                ntime: datalen > 200? 'scatter':'bar',
                sys: datalen > 200? 'scatter':'line',
                dia: datalen > 200? 'scatter':'line',
                pulse: datalen > 200? 'scatter':'line',
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
                if (d.id) {
                    if (d.id == 'pulse' && d.value) {
                        return  pulseColors[chartData.json[d.index].pulseScore];
                    }
                    else if (d.id == 'sys' && d.value) {
                        return  bpColors[chartData.json[d.index].sysScore];
                    }
                    else                if (d.id == 'dia' && d.value) {
                        return  bpColors[chartData.json[d.index].diaScore];
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
            show: false
        },
        zoom: {
            enabled: true,
        },
        axis: {
            y2: {
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
                show: datalen > 100? false:true,
                localtime: true,
                type: dataOptions.timeSeries ? 'timeseries':'category',
                tick: {
                    multiline: false,
                    rotate: 25,
                    fit:  true,
                    culling: {
                        max: 16,
                    },
                    format: function (x) {
                        if (dataOptions.timeSeries) {
                            var dt =  moment(x);
                            return dt.format('MM-DD  HH:mm');
                        } else {
                            var dt = moment(chartData.dateStrings[x]);
                            return dt.format('MM-DD  HH:mm');
                        }
                    },
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
                fetchData({timeSeries:false});
            }}>all</button>
            <button onClick={function(){
                fetchData({timeSeries:true});
            }}>timeseries</button>

            <Chart options={componentData} element="dddd"  />

        </div>,
        document.getElementById('content')
    );
}


var fetchData = function(opts){
    var url ="/json/?dt_min=1900-01-01";

    if (opts.timeSeries) {
        url += "&avg_minutes=10";
    } else {
        url += "&avg_minutes=0";
    }
    request
        .get(url)
        .set('Accept', 'application/json')
        .end(function(err, res){
            if (!res.ok) {
                console.log(err);
                return;
            }
            requestData = res.body;
            render(prepareData(requestData, opts));
        });
}

fetchData({timeSeries:true});
