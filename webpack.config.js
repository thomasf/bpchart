/* global __dirname, process, module, require */
var webpack = require("webpack");

var entry = {
    app: ["./browser/script/index"]
};

// Add react hot reloader if ALKASIR_HOT is set
if (process.env.ALKASIR_HOT !== undefined ) {
    entry.app = [
        "webpack-dev-server/client?http://0.0.0.0:3000", // WebpackDevServer host and port
        "webpack/hot/only-dev-server",
        "./browser/script/index"
    ];
}

var c = {
    module: {
        loaders: [
            { test: /\.js$/, exclude: /node_modules/, loader: "babel-loader"},
            { test: /\.jsx$/, loaders: ["react-hot-loader", "babel-loader"] },
            { test: /\.css$/, loader: "style!css" },
            { test: /\.woff$/, loader: "file" },
            { test: /\.woff2$/, loader: "file" },
            { test: /\.ttf$/, loader: "file" },
            { test: /\.eot$/, loader: "file" },
            { test: /\.svg$/, loader: "file" },
        ]
    },
    plugins: [
        // new webpack.HotModuleReplacementPlugin()
        // new webpack.NoErrorsPlugin()
    ],
    entry: entry,
    output: {
        path: __dirname + "/build/assets/",
        filename: "[name].bundle.js",
        publicPath: "/assets/"
    },
    resolve: {
        extensions: ["", ".js", ".jsx"]
    },
    devtool: "cheap-source-map"
};


module.exports = c;
