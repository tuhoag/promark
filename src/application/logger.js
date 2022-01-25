
const path = require("path");
const { createLogger, format, transports } = require('winston');
const { combine, timestamp, label, printf } = format;

const myFormat = printf(({ level, message, label, timestamp }) => {
    return `${timestamp} [${label}] ${level}: ${message}`;
});

module.exports = (name, level) => createLogger({
    // format: combine(
    //     label({ label: 'promark' }),
    //     timestamp(),
    //     myFormat
    // ),
    transports: [
        new transports.Console({
            level: level,
            json: false,
            timestamp: true,
            prettyPrint: true,
            depth: true,
            colorize: true,
            format: combine(
                label({ label: path.basename(name) }),
                timestamp(),
                myFormat
            ),
        })
        // new transports.File({ filename: 'error.log', level: 'error'}),
    ]
});