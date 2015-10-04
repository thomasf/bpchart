/**
 * React-C3 Chart
 * Copyright 2015 - Cody Reichert <codyreichert@gmail.com>
 */

import c3    from 'c3';
import React from 'react';


let ChartComponent = React.createClass({

    displayName: 'React-C3-Chart',

    propTypes: {
        options: React.PropTypes.object.isRequired,
        element: React.PropTypes.string.isRequired,
    },

    chart: null,

    componentDidMount: function() {
        if (this.props.options == null) {
            return
        }
        this._generateChart(
            this.props.options,
            this.props.element
        );
    },

    componentDidUpdate: function(prevProps) {
        if (this.props.options == null) {
            return
        }

        if(prevProps.options.data.columns !== this.props.options.data.columns) {
            this._generateChart(
                this.props.options,
                this.props.element
            );
        }
    },

    componentWillUnmount: function() {
        this._destroyChart();
    },

    _generateChart: function(options, element) {
        options.bindto = `#${element}`;
        this.chart = c3.generate(options);
    },

    _destroyChart: function() {
        this.chart.destroy();
    },

    render: function() {
        return (
            <div className="c3 react-c3"
                 id={this.props.element}
                 style={this.props.styles}>
            </div>
        );
    }
});


export default ChartComponent;
