
const path = require("path");
const { createLogger, format, transports } = require('winston');
const { combine, timestamp, label, printf } = format;

const myFormat = printf(({ level, message, label, timestamp }) => {
    return `${timestamp} [${label}] ${level}: ${message}`;
});

exports.createLogger = (name, level) => {
    const setting = require('./setting');

    if (level === undefined) {
        currentLevel = setting['logLevel'];
    } else {
        currentLevel = level
    }
    return createLogger({
        transports: [
            new transports.Console({
                level: currentLevel,
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
}